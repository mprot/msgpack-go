package msgpack

import (
	"encoding"
	"encoding/binary"
	"io"
	"math"
	"time"
)

// Writer defines a writer for MessagePack encoded data.
type Writer struct {
	w io.Writer
}

// NewWriter creates a writer for MessagePack encoded data which writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// WriteNil writes a nil value to the MessagePack stream.
func (w *Writer) WriteNil() error {
	buf := [1]byte{tagNil}
	_, err := w.w.Write(buf[:])
	return err
}

// WriteBool writes a boolean value to the MessagePack stream.
func (w *Writer) WriteBool(b bool) error {
	buf := [1]byte{tagFalse}
	if b {
		buf[0] = tagTrue
	}

	_, err := w.w.Write(buf[:])
	return err
}

// WriteInt8 writes an 8-bit integer value to the MessagePack stream.
func (w *Writer) WriteInt8(i int8) error {
	var buf [2]byte
	switch {
	case i >= 0:
		buf[0] = posFixintTag(uint8(i))
		_, err := w.w.Write(buf[:1])
		return err

	case i > -32:
		buf[0] = negFixintTag(i)
		_, err := w.w.Write(buf[:1])
		return err

	default:
		buf[0] = tagInt8
		buf[1] = byte(i)
		_, err := w.w.Write(buf[:2])
		return err
	}
}

// WriteInt16 writes a 16-bit integer value to the MessagePack stream.
func (w *Writer) WriteInt16(i int16) error {
	switch {
	case math.MinInt8 <= i && i <= math.MaxInt8:
		return w.WriteInt8(int8(i))
	default:
		buf := [3]byte{tagInt16}
		binary.BigEndian.PutUint16(buf[1:], uint16(i))
		_, err := w.w.Write(buf[:])
		return err
	}
}

// WriteInt32 writes a 32-bit integer value to the MessagePack stream.
func (w *Writer) WriteInt32(i int32) error {
	switch {
	case math.MinInt8 <= i && i <= math.MaxInt8:
		return w.WriteInt8(int8(i))
	case math.MinInt16 <= i && i <= math.MaxInt16:
		return w.WriteInt16(int16(i))
	default:
		buf := [5]byte{tagInt32}
		binary.BigEndian.PutUint32(buf[1:], uint32(i))
		_, err := w.w.Write(buf[:])
		return err
	}
}

// WriteInt64 writes a 64-bit integer value to the MessagePack stream.
func (w *Writer) WriteInt64(i int64) error {
	switch {
	case math.MinInt8 <= i && i <= math.MaxInt8:
		return w.WriteInt8(int8(i))
	case math.MinInt16 <= i && i <= math.MaxInt16:
		return w.WriteInt16(int16(i))
	case math.MinInt32 <= i && i <= math.MaxInt32:
		return w.WriteInt32(int32(i))
	default:
		buf := [9]byte{tagInt64}
		binary.BigEndian.PutUint64(buf[1:], uint64(i))
		_, err := w.w.Write(buf[:])
		return err
	}
}

// WriteInt writes an integer value to the MessagePack stream.
func (w *Writer) WriteInt(i int) error {
	return w.WriteInt64(int64(i))
}

// WriteUint8 writes an 8-bit unsigned integer value to the MessagePack stream.
func (w *Writer) WriteUint8(i uint8) error {
	var buf [2]byte
	switch {
	case i < 128:
		buf[0] = posFixintTag(i)
		_, err := w.w.Write(buf[:1])
		return err

	default:
		buf[0] = tagUint8
		buf[1] = i
		_, err := w.w.Write(buf[:])
		return err
	}
}

// WriteUint16 writes a 16-bit unsigned integer value to the MessagePack stream.
func (w *Writer) WriteUint16(i uint16) error {
	switch {
	case i <= math.MaxUint8:
		return w.WriteUint8(uint8(i))
	default:
		buf := [3]byte{tagUint16}
		binary.BigEndian.PutUint16(buf[1:], i)
		_, err := w.w.Write(buf[:])
		return err
	}
}

// WriteUint32 writes a 32-bit unsigned integer value to the MessagePack stream.
func (w *Writer) WriteUint32(i uint32) error {
	switch {
	case i <= math.MaxUint8:
		return w.WriteUint8(uint8(i))
	case i <= math.MaxUint16:
		return w.WriteUint16(uint16(i))
	default:
		buf := [5]byte{tagUint32}
		binary.BigEndian.PutUint32(buf[1:], i)
		_, err := w.w.Write(buf[:])
		return err
	}
}

// WriteUint64 writes a 64-bit unsigned integer value to the MessagePack stream.
func (w *Writer) WriteUint64(i uint64) error {
	switch {
	case i <= math.MaxUint8:
		return w.WriteUint8(uint8(i))
	case i <= math.MaxUint16:
		return w.WriteUint16(uint16(i))
	case i <= math.MaxUint32:
		return w.WriteUint32(uint32(i))
	default:
		buf := [9]byte{tagUint64}
		binary.BigEndian.PutUint64(buf[1:], i)
		_, err := w.w.Write(buf[:])
		return err
	}
}

// WriteUint writes an unsigned integer value to the MessagePack stream.
func (w *Writer) WriteUint(i uint) error {
	return w.WriteUint64(uint64(i))
}

// WriteFloat32 writes a 32-bit floating-point value to the MessagePack stream.
func (w *Writer) WriteFloat32(f float32) error {
	buf := [5]byte{tagFloat32}
	binary.BigEndian.PutUint32(buf[1:], math.Float32bits(f))
	_, err := w.w.Write(buf[:])
	return err
}

// WriteFloat64 writes a 64-bit floating-point value to the MessagePack stream.
func (w *Writer) WriteFloat64(f float64) error {
	buf := [9]byte{tagFloat64}
	binary.BigEndian.PutUint64(buf[1:], math.Float64bits(f))
	_, err := w.w.Write(buf[:])
	return err
}

// WriteBytes writes a binary value to the MessagePack stream.
func (w *Writer) WriteBytes(b []byte) error {
	return w.writeBlob(tagBin8, b)
}

// WriteString writes a string value to the MessagePack stream.
func (w *Writer) WriteString(s string) error {
	if len(s) <= 31 {
		buf := [1]byte{fixstrTag(len(s))}
		_, err := w.w.Write(buf[:])
		if err == nil {
			_, err = io.WriteString(w.w, s)
		}
		return err
	}
	return w.writeBlob(tagStr8, []byte(s))
}

// WriteArrayHeader writes the header of an array value to the MessagePack stream.
func (w *Writer) WriteArrayHeader(length int) error {
	if length <= 15 {
		buf := [1]byte{fixarrayTag(length)}
		_, err := w.w.Write(buf[:])
		return err
	}
	return w.writeCollectionHeader(tagArray16, length)
}

// WriteMapHeader writes the header of an map value to the MessagePack stream.
func (w *Writer) WriteMapHeader(length int) error {
	if length <= 15 {
		buf := [1]byte{fixmapTag(length)}
		_, err := w.w.Write(buf[:])
		return err
	}
	return w.writeCollectionHeader(tagMap16, length)
}

// WriteRaw writes raw bytes to the MessagePack stream, which represent an
// already encoded section.
func (w *Writer) WriteRaw(r Raw) error {
	_, err := w.w.Write(r)
	return err
}

// WriteExt writes an extension value to the MessagePack stream.
func (w *Writer) WriteExt(typ int8, v encoding.BinaryMarshaler) error {
	if typ < 0 {
		return newInvalidExtensionError(typ)
	}

	data, err := v.MarshalBinary()
	if err != nil {
		return err
	}
	return w.writeExtension(typ, data)
}

// WriteTime writes a time value to the MessagePack stream.
func (w *Writer) WriteTime(tm time.Time) error {
	// To avoid undefined behaviour for timestamps before the
	// year 1678 or after 2262, the 64-bit UNIX timestamp is
	// used here.
	var buf [12]byte
	binary.BigEndian.PutUint64(buf[4:], uint64(tm.Unix()))
	return w.writeExtension(extTime, buf[:])
}

func (w *Writer) writeBlob(baseTag byte, blob []byte) error {
	var (
		buf [5]byte
		p   []byte
	)

	n := len(blob)
	switch {
	case n <= math.MaxUint8:
		buf[0] = baseTag
		buf[1] = byte(n)
		p = buf[:2]

	case n <= math.MaxUint16:
		buf[0] = baseTag + 1
		binary.BigEndian.PutUint16(buf[1:], uint16(n))
		p = buf[:3]

	case n <= math.MaxUint32:
		buf[0] = baseTag + 2
		binary.BigEndian.PutUint32(buf[1:], uint32(n))
		p = buf[:5]

	default:
		return errLengthLimitExceeded
	}

	_, err := w.w.Write(p)
	if err == nil {
		_, err = w.w.Write(blob)
	}
	return err
}

func (w *Writer) writeCollectionHeader(baseTag byte, length int) error {
	var buf [5]byte
	switch {
	case length <= math.MaxUint16:
		buf[0] = baseTag
		binary.BigEndian.PutUint16(buf[1:], uint16(length))
		_, err := w.w.Write(buf[:3])
		return err

	case length <= math.MaxUint32:
		buf[0] = baseTag + 1
		binary.BigEndian.PutUint32(buf[1:], uint32(length))
		_, err := w.w.Write(buf[:5])
		return err

	default:
		return errLengthLimitExceeded
	}
}

func (w *Writer) writeExtension(typ int8, data []byte) error {
	n := len(data)
	switch n {
	case 1:
		return w.writeFixExt(tagFixExt1, typ, data)
	case 2:
		return w.writeFixExt(tagFixExt2, typ, data)
	case 4:
		return w.writeFixExt(tagFixExt4, typ, data)
	case 8:
		return w.writeFixExt(tagFixExt8, typ, data)
	case 16:
		return w.writeFixExt(tagFixExt16, typ, data)
	}

	var (
		buf [6]byte
		p   []byte
	)
	switch {
	case n <= math.MaxUint8:
		buf[0] = tagExt8
		buf[1] = byte(n)
		buf[2] = byte(typ)
		p = buf[:3]

	case n <= math.MaxUint16:
		buf[0] = tagExt16
		binary.BigEndian.PutUint16(buf[1:], uint16(n))
		buf[3] = byte(typ)
		p = buf[:4]

	case n <= math.MaxUint32:
		buf[0] = tagExt32
		binary.BigEndian.PutUint32(buf[1:], uint32(n))
		buf[5] = byte(typ)
		p = buf[:6]

	default:
		return errLengthLimitExceeded
	}

	_, err := w.w.Write(p)
	if err == nil {
		_, err = w.w.Write(data)
	}
	return err
}

func (w *Writer) writeFixExt(tag byte, typ int8, data []byte) error {
	buf := [2]byte{tag, byte(typ)}
	_, err := w.w.Write(buf[:])
	if err == nil {
		_, err = w.w.Write(data)
	}
	return err
}
