package msgpack

import (
	"encoding"
	"encoding/binary"
	"io"
	"math"
	"sync"
	"time"
)

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

var readerPool sync.Pool

// Reader defines a reader for MessagePack encoded data.
type Reader struct {
	r     io.Reader
	buf   []byte
	first int
	last  int
	err   error
}

// NewReader creates a reader for MessagePack encoded data read from r.
func NewReader(r io.Reader) *Reader {
	if v := readerPool.Get(); v != nil {
		reader := v.(*Reader)
		reader.r = r
		reader.first = 0
		reader.last = 0
		reader.err = nil
		return reader
	}

	return &Reader{
		r:   r,
		buf: make([]byte, 1024),
	}
}

// NewReaderBytes creates a reader for MessagePack encoded data. The reader
// reads all data from p.
func NewReaderBytes(p []byte) *Reader {
	return &Reader{
		r:     eofReader{},
		buf:   p,
		first: 0,
		last:  len(p),
	}
}

func releaseReader(r *Reader) {
	r.r = nil
	readerPool.Put(r)
}

// Peek returns the type for the next element without moving the
// read pointer.
func (r *Reader) Peek() (Type, error) {
	tag, err := r.peek()
	if err != nil {
		return "", err
	}
	return r.peekType(tag)
}

// ReadNil reads a nil value from the MessagePack stream.
func (r *Reader) ReadNil() error {
	tag, err := r.peek()
	switch {
	case err != nil:
		return err
	case tag != tagNil:
		return r.typeErr(tag, Nil)
	default:
		r.advance(1)
		return nil
	}
}

// ReadBool reads a boolean value from the MessagePack stream.
func (r *Reader) ReadBool() (bool, error) {
	tag, err := r.peek()
	if err != nil {
		return false, err
	}

	switch tag {
	case tagNil, tagFalse:
		r.advance(1)
		return false, nil
	case tagTrue:
		r.advance(1)
		return true, nil
	default:
		return false, r.typeErr(tag, Bool)
	}
}

// ReadInt8 reads an 8-bit integer value from the MessagePack stream.
func (r *Reader) ReadInt8() (int8, error) {
	i, err := r.ReadInt64()
	switch {
	case err != nil:
		return 0, err
	case i < math.MinInt8 || i > math.MaxInt8:
		return 0, intOverflowError{}
	default:
		return int8(i), nil
	}
}

// ReadInt16 reads a 16-bit integer value from the MessagePack stream.
func (r *Reader) ReadInt16() (int16, error) {
	i, err := r.ReadInt64()
	switch {
	case err != nil:
		return 0, err
	case i < math.MinInt16 || i > math.MaxInt16:
		return 0, intOverflowError{}
	default:
		return int16(i), nil
	}
}

// ReadInt32 reads a 32-bit integer value from the MessagePack stream.
func (r *Reader) ReadInt32() (int32, error) {
	i, err := r.ReadInt64()
	switch {
	case err != nil:
		return 0, err
	case i < math.MinInt32 || i > math.MaxInt32:
		return 0, intOverflowError{}
	default:
		return int32(i), nil
	}
}

// ReadInt64 reads a 64-bit integer value from the MessagePack stream.
func (r *Reader) ReadInt64() (int64, error) {
	tag, err := r.peek()
	if err != nil {
		return 0, err
	}

	switch {
	case isPosFixintTag(tag):
		i := readPosFixint(tag)
		r.advance(1)
		return int64(i), nil

	case isNegFixintTag(tag):
		i := readNegFixint(tag)
		r.advance(1)
		return int64(i), nil
	}

	switch tag {
	case tagNil:
		r.advance(1)
		return 0, nil

	// signed int types
	case tagInt8:
		buf, err := r.read(2)
		if err != nil {
			return 0, err
		}
		return int64(int8(buf[1])), nil

	case tagInt16:
		buf, err := r.read(3)
		if err != nil {
			return 0, err
		}
		return int64(int16(binary.BigEndian.Uint16(buf[1:]))), nil

	case tagInt32:
		buf, err := r.read(5)
		if err != nil {
			return 0, err
		}
		return int64(int32(binary.BigEndian.Uint32(buf[1:]))), nil

	case tagInt64:
		buf, err := r.read(9)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint64(buf[1:])), nil

	// unsigned int types
	case tagUint8:
		buf, err := r.read(2)
		if err != nil {
			return 0, err
		}
		return int64(buf[1]), nil

	case tagUint16:
		buf, err := r.read(3)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint16(buf[1:])), nil

	case tagUint32:
		buf, err := r.read(5)
		if err != nil {
			return 0, err
		}
		return int64(binary.BigEndian.Uint32(buf[1:])), nil

	case tagUint64:
		buf, err := r.read(9)
		if err != nil {
			return 0, err
		}
		ui := binary.BigEndian.Uint64(buf[1:])
		if ui > math.MaxInt64 {
			return 0, intOverflowError{}
		}
		return int64(ui), nil

	default:
		return 0, r.typeErr(tag, Int)
	}
}

// ReadInt reads an integer value from the MessagePack stream.
func (r *Reader) ReadInt() (int, error) {
	i, err := r.ReadInt64()
	switch {
	case err != nil:
		return 0, err
	case i < int64(minInt) || i > int64(maxInt):
		return 0, intOverflowError{}
	default:
		return int(i), nil
	}
}

// ReadUint8 reads an 8-bit unsigned integer value from the MessagePack stream.
func (r *Reader) ReadUint8() (uint8, error) {
	ui, err := r.ReadUint64()
	switch {
	case err != nil:
		return 0, err
	case ui > math.MaxUint8:
		return 0, intOverflowError{}
	default:
		return uint8(ui), nil
	}
}

// ReadUint16 reads a 16-bit unsigned integer value from the MessagePack stream.
func (r *Reader) ReadUint16() (uint16, error) {
	ui, err := r.ReadUint64()
	switch {
	case err != nil:
		return 0, err
	case ui > math.MaxUint16:
		return 0, intOverflowError{}
	default:
		return uint16(ui), nil
	}
}

// ReadUint32 reads a 32-bit unsigned integer value from the MessagePack stream.
func (r *Reader) ReadUint32() (uint32, error) {
	ui, err := r.ReadUint64()
	switch {
	case err != nil:
		return 0, err
	case ui > math.MaxUint32:
		return 0, intOverflowError{}
	default:
		return uint32(ui), nil
	}
}

// ReadUint64 reads a 64-bit unsigned integer value from the MessagePack stream.
func (r *Reader) ReadUint64() (uint64, error) {
	tag, err := r.peek()
	if err != nil {
		return 0, err
	}

	if isPosFixintTag(tag) {
		i := readPosFixint(tag)
		r.advance(1)
		return uint64(i), nil
	}

	switch tag {
	case tagNil:
		r.advance(1)
		return 0, nil

	// signed int types
	case tagInt8:
		buf, err := r.read(2)
		if err != nil {
			return 0, err
		}
		i8 := int8(buf[1])
		if i8 < 0 {
			return 0, intOverflowError{}
		}
		return uint64(i8), nil

	case tagInt16:
		buf, err := r.read(3)
		if err != nil {
			return 0, err
		}
		i16 := int16(binary.BigEndian.Uint16(buf[1:]))
		if i16 < 0 {
			return 0, intOverflowError{}
		}
		return uint64(i16), nil

	case tagInt32:
		buf, err := r.read(5)
		if err != nil {
			return 0, err
		}
		i32 := int32(binary.BigEndian.Uint32(buf[1:]))
		if i32 < 0 {
			return 0, intOverflowError{}
		}
		return uint64(i32), nil

	case tagInt64:
		buf, err := r.read(9)
		if err != nil {
			return 0, err
		}
		i64 := int64(binary.BigEndian.Uint64(buf[1:]))
		if i64 < 0 {
			return 0, intOverflowError{}
		}
		return uint64(i64), nil

	// unsigned int types
	case tagUint8:
		buf, err := r.read(2)
		if err != nil {
			return 0, err
		}
		return uint64(buf[1]), nil

	case tagUint16:
		buf, err := r.read(3)
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint16(buf[1:])), nil

	case tagUint32:
		buf, err := r.read(5)
		if err != nil {
			return 0, err
		}
		return uint64(binary.BigEndian.Uint32(buf[1:])), nil

	case tagUint64:
		buf, err := r.read(9)
		if err != nil {
			return 0, err
		}
		return binary.BigEndian.Uint64(buf[1:]), nil

	default:
		return 0, r.typeErr(tag, Uint)
	}
}

// ReadUint reads an unsigned integer value from the MessagePack stream.
func (r *Reader) ReadUint() (uint, error) {
	ui, err := r.ReadUint64()
	switch {
	case err != nil:
		return 0, err
	case ui > uint64(maxUint):
		return 0, intOverflowError{}
	default:
		return uint(ui), nil
	}
}

// ReadFloat32 reads a 32-bit floating-point value from the MessagePack stream.
func (r *Reader) ReadFloat32() (float32, error) {
	f, err := r.ReadFloat64()
	switch {
	case err != nil:
		return 0, err
	case f < -math.MaxFloat32 || f > math.MaxFloat32:
		return 0, floatOverflowError{}
	default:
		return float32(f), nil
	}
}

// ReadFloat64 reads a 64-bit floating-point value from the MessagePack stream.
func (r *Reader) ReadFloat64() (float64, error) {
	tag, err := r.peek()
	if err != nil {
		return 0, err
	}

	switch tag {
	case tagNil:
		r.advance(1)
		return 0, nil

	case tagFloat32:
		buf, err := r.read(5)
		if err != nil {
			return 0, err
		}
		return float64(math.Float32frombits(binary.BigEndian.Uint32(buf[1:]))), nil

	case tagFloat64:
		buf, err := r.read(9)
		if err != nil {
			return 0, err
		}
		return math.Float64frombits(binary.BigEndian.Uint64(buf[1:])), nil

	default:
		return 0, r.typeErr(tag, Float)
	}
}

// ReadBytes reads a binary value from the MessagePack stream. The data
// will be copied to p if it fits. Otherwise a new slice will be allocated.
func (r *Reader) ReadBytes(p []byte) ([]byte, error) {
	n, err := r.readBlobHeader(Bytes)
	if err != nil || n == 0 {
		return nil, err
	}

	if cap(p) < n {
		p = make([]byte, n)
	}
	p = p[:n]
	err = r.readFull(p)
	return p, err
}

// ReadBytesNoCopy reads a binary value from the MessagePack stream. The data
// is directly passed from the underlying buffer. This could result in a buffer
// reallocation.
func (r *Reader) ReadBytesNoCopy() ([]byte, error) {
	n, err := r.readBlobHeader(Bytes)
	if err != nil {
		return nil, err
	}
	return r.read(n)
}

// ReadString reads a string value from the MessagePack stream.
func (r *Reader) ReadString() (string, error) {
	n, err := r.readBlobHeader(String)
	if err != nil || n == 0 {
		return "", err
	}

	p := make([]byte, n)
	err = r.readFull(p)
	return string(p), err
}

// ReadArrayHeader reads the header of an array value from the MessagePack stream
// and returns the number of elements.
func (r *Reader) ReadArrayHeader() (int, error) {
	return r.readCollectionHeader(tagArray16, Array, func(tag byte) (uint8, bool) {
		if isFixarrayTag(tag) {
			return readFixarray(tag), true
		}
		return 0, false
	})
}

// ReadArrayHeaderWithSize reads the header of an array value from the MessagePack
// stream and compares it to the given size. If the size of the array does not
// equal the given size, an error will be returned.
func (r *Reader) ReadArrayHeaderWithSize(size int) error {
	n, err := r.ReadArrayHeader()
	switch {
	case err != nil:
		return err
	case n != size:
		return errorf("invalid array header size %d", n)
	default:
		return nil
	}
}

// ReadMapHeader reads the header of an map value from the MessagePack stream
// and returns the number of elements.
func (r *Reader) ReadMapHeader() (int, error) {
	return r.readCollectionHeader(tagMap16, Map, func(tag byte) (uint8, bool) {
		if isFixmapTag(tag) {
			return readFixmap(tag), true
		}
		return 0, false
	})
}

// ReadRaw reads the next value from the MessagePack stream into raw.
func (r *Reader) ReadRaw(raw Raw) (Raw, error) {
	raw = raw[:0]
	err := r.readRaw(func(p []byte) { raw = append(raw, p...) })
	return raw, err
}

// ReadExt reads an extension value from the MessagePack stream. The given extension type
// must match the extension type found in the stream.
func (r *Reader) ReadExt(typ int8, v encoding.BinaryUnmarshaler) error {
	data, err := r.readExtension(typ)
	if err != nil {
		return err
	}
	return v.UnmarshalBinary(data)
}

// ReadTime reads a time value from the MessagePack stream.
func (r *Reader) ReadTime() (time.Time, error) {
	data, err := r.readExtension(extTime)
	if err != nil {
		return time.Time{}, err
	}

	switch len(data) {
	case 4: // 32-bit seconds
		seconds := binary.BigEndian.Uint32(data)
		return time.Unix(int64(seconds), 0).UTC(), nil
	case 8: // 34-bit seconds + 30-bit nanoseconds
		tm := binary.BigEndian.Uint64(data)
		return time.Unix(int64(tm&0x3ffffffff), int64(tm>>34)).UTC(), nil
	case 12: // 64-bit seconds + 32-bit nanoseconds
		secs := binary.BigEndian.Uint64(data[4:])
		nsecs := binary.BigEndian.Uint32(data)
		return time.Unix(int64(secs), int64(nsecs)).UTC(), nil
	default:
		return time.Time{}, errorf("invalid timestamp length %d", len(data))
	}
}

// Skip skips the next value in the MessagePack stream.
func (r *Reader) Skip() error {
	return r.readRaw(func([]byte) {})
}

func (r *Reader) readBlobHeader(expectedType Type) (int, error) {
	tag, err := r.peek()
	if err != nil {
		return 0, err
	}

	switch tag {
	case tagNil:
		r.advance(1)
		return 0, nil

	case tagBin8, tagStr8:
		buf, err := r.read(2)
		if err != nil {
			return 0, err
		}
		return int(buf[1]), nil

	case tagBin16, tagStr16:
		buf, err := r.read(3)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(buf[1:])), nil

	case tagBin32, tagStr32:
		buf, err := r.read(5)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(buf[1:])), nil

	default:
		if !isFixstrTag(tag) {
			return 0, r.typeErr(tag, expectedType)
		}
		r.advance(1)
		return int(readFixstr(tag)), nil
	}
}

func (r *Reader) readCollectionHeader(tagBase byte, expectedTyp Type, readFix func(byte) (uint8, bool)) (int, error) {
	tag, err := r.peek()
	if err != nil {
		return 0, err
	}

	if n, ok := readFix(tag); ok {
		r.advance(1)
		return int(n), nil
	}

	switch tag {
	case tagNil:
		r.advance(1)
		return 0, nil

	case tagBase: // 16 bit
		buf, err := r.read(3)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16(buf[1:])), nil

	case tagBase + 1: // 32 bit
		buf, err := r.read(5)
		if err != nil {
			return 0, err
		}
		n := binary.BigEndian.Uint32(buf[1:])
		if uint(n) > uint(maxInt) {
			return 0, intOverflowError{}
		}
		return int(n), nil

	default:
		return 0, r.typeErr(tag, expectedTyp)
	}
}

func (r *Reader) readExtension(typ int8) ([]byte, error) {
	tag, err := r.peek()
	if err != nil {
		return nil, err
	}

	var (
		header []byte
		n      int
	)
	switch tag {
	case tagFixExt1:
		header, err = r.peekn(2)
		n = 1
	case tagFixExt2:
		header, err = r.peekn(2)
		n = 2
	case tagFixExt4:
		header, err = r.peekn(2)
		n = 4
	case tagFixExt8:
		header, err = r.peekn(2)
		n = 8
	case tagFixExt16:
		header, err = r.peekn(2)
		n = 16
	case tagExt8:
		if header, err = r.peekn(3); err == nil {
			n = int(header[1])
		}
	case tagExt16:
		if header, err = r.peekn(4); err == nil {
			n = int(binary.BigEndian.Uint16(header[1:]))
		}
	case tagExt32:
		if header, err = r.peekn(6); err == nil {
			n = int(binary.BigEndian.Uint32(header[1:]))
		}
	default:
		return nil, r.typeErr(tag, Ext)
	}

	if err != nil {
		return nil, err
	} else if t := int8(header[len(header)-1]); typ != t {
		return nil, newInvalidExtensionError(t)
	}

	data, err := r.read(len(header) + n)
	if err != nil {
		return nil, err
	}
	return data[len(header):], nil
}

func (r *Reader) typeErr(tag byte, expected Type) error {
	actual, err := r.peekType(tag)
	if err != nil {
		actual = Ext
	}
	return TypeError{
		Actual:   actual,
		Expected: expected,
	}
}

// precondition: tag is not consumed
func (r *Reader) peekType(tag byte) (Type, error) {
	switch tag {
	case tagFixExt1, tagFixExt2, tagFixExt4, tagFixExt8, tagFixExt16:
		p, err := r.peekn(2)
		if err != nil {
			return "", err
		}
		return extType(int8(p[1])), nil
	case tagExt8:
		p, err := r.peekn(3)
		if err != nil {
			return "", err
		}
		return extType(int8(p[2])), nil
	case tagExt16:
		p, err := r.peekn(4)
		if err != nil {
			return "", err
		}
		return extType(int8(p[3])), nil
	case tagExt32:
		p, err := r.peekn(6)
		if err != nil {
			return "", err
		}
		return extType(int8(p[5])), nil
	}

	return tagType(tag), nil
}

func (r *Reader) peek() (byte, error) {
	if r.first == r.last {
		if err := r.fillBuf(1); err != nil {
			return 0, err
		}
	}
	return r.buf[r.first], nil
}

func (r *Reader) peekn(n int) ([]byte, error) {
	if r.first+n > r.last {
		if err := r.fillBuf(n); err != nil {
			return nil, err
		}
	}
	return r.buf[r.first : r.first+n], nil
}

func (r *Reader) read(n int) ([]byte, error) {
	p, err := r.peekn(n)
	if err == nil {
		r.advance(n)
	}
	return p, err
}

func (r *Reader) readFull(p []byte) error {
	var n int
	if r.first != r.last {
		n = copy(p, r.buf[r.first:r.last])
		r.advance(n)
	}
	if n == len(p) {
		return nil
	}

	if r.err == nil {
		if _, r.err = io.ReadFull(r.r, p[n:]); r.err == io.EOF {
			r.err = io.ErrUnexpectedEOF
		}
	}
	return r.err
}

func (r *Reader) readRaw(f func([]byte)) error {
	tag, err := r.peek()
	if err != nil {
		return err
	}

	switch tag {
	case tagNil, tagFalse, tagTrue:
		p, err := r.read(1)
		f(p)
		return err
	case tagInt8, tagUint8:
		p, err := r.read(2)
		f(p)
		return err
	case tagInt16, tagUint16:
		p, err := r.read(3)
		f(p)
		return err
	case tagInt32, tagUint32, tagFloat32:
		p, err := r.read(5)
		f(p)
		return err
	case tagInt64, tagUint64, tagFloat64:
		p, err := r.read(9)
		f(p)
		return err
	case tagStr8, tagBin8:
		p, err := r.peekn(2)
		if err == nil {
			p, err = r.read(2 + int(p[1]))
			f(p)
		}
		return err
	case tagStr16, tagBin16:
		p, err := r.peekn(3)
		if err == nil {
			p, err = r.read(3 + int(binary.BigEndian.Uint16(p[1:])))
			f(p)
		}
		return err
	case tagStr32, tagBin32:
		p, err := r.peekn(5)
		if err == nil {
			p, err = r.read(5 + int(binary.BigEndian.Uint32(p[1:])))
			f(p)
		}
		return err
	case tagArray16:
		p, err := r.read(3)
		if err != nil {
			return err
		}
		f(p)
		for i, n := uint16(0), binary.BigEndian.Uint16(p[1:]); i < n; i++ {
			if err := r.readRaw(f); err != nil {
				return err
			}
		}
		return nil
	case tagArray32:
		p, err := r.read(5)
		if err != nil {
			return err
		}
		f(p)
		for i, n := uint32(0), binary.BigEndian.Uint32(p[1:]); i < n; i++ {
			if err := r.readRaw(f); err != nil {
				return err
			}
		}
		return nil
	case tagMap16:
		p, err := r.read(3)
		if err != nil {
			return err
		}
		f(p)
		for i, n := uint16(0), 2*binary.BigEndian.Uint16(p[1:]); i < n; i++ {
			if err := r.readRaw(f); err != nil {
				return err
			}
		}
		return nil
	case tagMap32:
		p, err := r.read(5)
		if err != nil {
			return err
		}
		f(p)
		for i, n := uint32(0), 2*binary.BigEndian.Uint32(p[1:]); i < n; i++ {
			if err := r.readRaw(f); err != nil {
				return err
			}
		}
		return nil
	case tagFixExt1:
		p, err := r.read(3)
		f(p)
		return err
	case tagFixExt2:
		p, err := r.read(4)
		f(p)
		return err
	case tagFixExt4:
		p, err := r.read(6)
		f(p)
		return err
	case tagFixExt8:
		p, err := r.read(10)
		f(p)
		return err
	case tagFixExt16:
		p, err := r.read(18)
		f(p)
		return err
	case tagExt8:
		p, err := r.peekn(3)
		if err == nil {
			p, err = r.read(3 + int(p[1]))
			f(p)
		}
		return err
	case tagExt16:
		p, err := r.peekn(4)
		if err == nil {
			p, err = r.read(4 + int(binary.BigEndian.Uint16(p[1:])))
			f(p)
		}
		return err
	case tagExt32:
		p, err := r.peekn(6)
		if err == nil {
			p, err = r.read(6 + int(binary.BigEndian.Uint32(p[1:])))
			f(p)
		}
		return err
	}

	switch {
	case isPosFixintTag(tag) || isNegFixintTag(tag):
		p, err := r.read(1)
		f(p)
		return err
	case isFixstrTag(tag):
		p, err := r.read(1 + int(readFixstr(tag)))
		f(p)
		return err
	case isFixarrayTag(tag):
		p, err := r.read(1)
		if err != nil {
			return err
		}
		f(p)
		for i, n := 0, int(readFixarray(tag)); i < n; i++ {
			if r.readRaw(f); err != nil {
				return err
			}
		}
		return nil
	case isFixmapTag(tag):
		p, err := r.read(1)
		if err != nil {
			return err
		}
		f(p)
		for i, n := 0, 2*int(readFixmap(tag)); i < n; i++ {
			if r.readRaw(f); err != nil {
				return err
			}
		}
		return nil
	}

	return errorf("unknown tag %#02x", tag)
}

func (r *Reader) advance(n int) {
	r.first += n
}

func (r *Reader) fillBuf(minSize int) error {
	if r.err != nil {
		return r.err
	}

	if r.first != 0 {
		copy(r.buf, r.buf[r.first:r.last])
		r.last -= r.first
		r.first = 0
	}
	if len(r.buf) < minSize {
		r.buf = make([]byte, minSize)
	}

	for minSize > r.last {
		n, err := r.r.Read(r.buf[r.last:])
		if err != nil {
			r.err = err
			return err
		}
		r.last += n
	}
	return nil
}

type eofReader struct{}

func (r eofReader) Read(p []byte) (int, error) { return 0, io.EOF }
