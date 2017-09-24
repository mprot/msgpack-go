package msgpack

import (
	"bytes"
	"math"
	"strings"
	"testing"
	"time"
)

func TestWriterWrite(t *testing.T) {
	tests := []struct {
		write func(*Writer) error
		data  []byte
	}{
		// nil
		{
			write: func(w *Writer) error { return w.WriteNil() },
			data:  []byte{tagNil},
		},
		// bool
		{
			write: func(w *Writer) error { return w.WriteBool(false) },
			data:  []byte{tagFalse},
		},
		{
			write: func(w *Writer) error { return w.WriteBool(true) },
			data:  []byte{tagTrue},
		},
		// int
		{
			write: func(w *Writer) error { return w.WriteInt64(7) },
			data:  []byte{posFixintTag(7)},
		},
		{
			write: func(w *Writer) error { return w.WriteInt64(-7) },
			data:  []byte{negFixintTag(-7)},
		},
		{
			write: func(w *Writer) error { return w.WriteInt64(math.MinInt8) },
			data:  []byte{tagInt8, 0x80},
		},
		{
			write: func(w *Writer) error { return w.WriteInt64(math.MinInt16) },
			data:  []byte{tagInt16, 0x80, 0x0},
		},
		{
			write: func(w *Writer) error { return w.WriteInt64(math.MinInt32) },
			data:  []byte{tagInt32, 0x80, 0x0, 0x0, 0x0},
		},
		{
			write: func(w *Writer) error { return w.WriteInt64(math.MinInt64) },
			data:  []byte{tagInt64, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
		// uint
		{
			write: func(w *Writer) error { return w.WriteUint64(7) },
			data:  []byte{posFixintTag(7)},
		},
		{
			write: func(w *Writer) error { return w.WriteUint64(math.MaxUint8) },
			data:  []byte{tagUint8, 0xff},
		},
		{
			write: func(w *Writer) error { return w.WriteUint64(math.MaxUint16) },
			data:  []byte{tagUint16, 0xff, 0xff},
		},
		{
			write: func(w *Writer) error { return w.WriteUint64(math.MaxUint32) },
			data:  []byte{tagUint32, 0xff, 0xff, 0xff, 0xff},
		},
		{
			write: func(w *Writer) error { return w.WriteUint64(math.MaxUint64) },
			data:  []byte{tagUint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		// float
		{
			write: func(w *Writer) error { return w.WriteFloat32(3.141592) },
			data:  []byte{tagFloat32, 0x40, 0x49, 0x0f, 0xd8},
		},
		{
			write: func(w *Writer) error { return w.WriteFloat64(3.141592) },
			data:  []byte{tagFloat64, 0x40, 0x09, 0x21, 0xfa, 0xfc, 0x8b, 0x00, 0x7a},
		},
		// bytes
		{
			write: func(w *Writer) error { return w.WriteBytes(bytes.Repeat([]byte{'0'}, 1)) },
			data:  append([]byte{tagBin8, 0x1}, bytes.Repeat([]byte{'0'}, 1)...),
		},
		{
			write: func(w *Writer) error { return w.WriteBytes(bytes.Repeat([]byte{'0'}, math.MaxUint8+1)) },
			data:  append([]byte{tagBin16, 0x1, 0x0}, bytes.Repeat([]byte{'0'}, math.MaxUint8+1)...),
		},
		{
			write: func(w *Writer) error { return w.WriteBytes(bytes.Repeat([]byte{'0'}, math.MaxUint16+1)) },
			data:  append([]byte{tagBin32, 0x0, 0x1, 0x0, 0x0}, bytes.Repeat([]byte{'0'}, math.MaxUint16+1)...),
		},
		// string
		{
			write: func(w *Writer) error { return w.WriteString(strings.Repeat("0", 7)) },
			data:  append([]byte{fixstrTag(7)}, bytes.Repeat([]byte{'0'}, 7)...),
		},
		{
			write: func(w *Writer) error { return w.WriteString(strings.Repeat("0", 32)) },
			data:  append([]byte{tagStr8, 0x20}, bytes.Repeat([]byte{'0'}, 32)...),
		},
		{
			write: func(w *Writer) error { return w.WriteString(strings.Repeat("0", math.MaxUint8+1)) },
			data:  append([]byte{tagStr16, 0x1, 0x0}, bytes.Repeat([]byte{'0'}, math.MaxUint8+1)...),
		},
		{
			write: func(w *Writer) error { return w.WriteString(strings.Repeat("0", math.MaxUint16+1)) },
			data:  append([]byte{tagStr32, 0x0, 0x01, 0x0, 0x0}, bytes.Repeat([]byte{'0'}, math.MaxUint16+1)...),
		},
		// array
		{
			write: func(w *Writer) error { return w.WriteArrayHeader(7) },
			data:  []byte{fixarrayTag(7)},
		},
		{
			write: func(w *Writer) error { return w.WriteArrayHeader(math.MaxUint16) },
			data:  []byte{tagArray16, 0xff, 0xff},
		},
		{
			write: func(w *Writer) error { return w.WriteArrayHeader(math.MaxUint32) },
			data:  []byte{tagArray32, 0xff, 0xff, 0xff, 0xff},
		},
		// map
		{
			write: func(w *Writer) error { return w.WriteMapHeader(7) },
			data:  []byte{fixmapTag(7)},
		},
		{
			write: func(w *Writer) error { return w.WriteMapHeader(math.MaxUint16) },
			data:  []byte{tagMap16, 0xff, 0xff},
		},
		{
			write: func(w *Writer) error { return w.WriteMapHeader(math.MaxUint32) },
			data:  []byte{tagMap32, 0xff, 0xff, 0xff, 0xff},
		},
		// ext
		{
			write: func(w *Writer) error { return w.WriteExt(0xd, bytesMarshaler(bytes.Repeat([]byte{'0'}, 1))) },
			data:  append([]byte{tagFixExt1, 0xd}, bytes.Repeat([]byte{'0'}, 1)...),
		},
		{
			write: func(w *Writer) error { return w.WriteExt(0xd, bytesMarshaler(bytes.Repeat([]byte{'0'}, 2))) },
			data:  append([]byte{tagFixExt2, 0xd}, bytes.Repeat([]byte{'0'}, 2)...),
		},
		{
			write: func(w *Writer) error { return w.WriteExt(0xd, bytesMarshaler(bytes.Repeat([]byte{'0'}, 4))) },
			data:  append([]byte{tagFixExt4, 0xd}, bytes.Repeat([]byte{'0'}, 4)...),
		},
		{
			write: func(w *Writer) error { return w.WriteExt(0xd, bytesMarshaler(bytes.Repeat([]byte{'0'}, 8))) },
			data:  append([]byte{tagFixExt8, 0xd}, bytes.Repeat([]byte{'0'}, 8)...),
		},
		{
			write: func(w *Writer) error { return w.WriteExt(0xd, bytesMarshaler(bytes.Repeat([]byte{'0'}, 16))) },
			data:  append([]byte{tagFixExt16, 0xd}, bytes.Repeat([]byte{'0'}, 16)...),
		},
		{
			write: func(w *Writer) error {
				return w.WriteExt(0xd, bytesMarshaler(bytes.Repeat([]byte{'0'}, 7)))
			},
			data: append([]byte{tagExt8, 0x7, 0xd}, bytes.Repeat([]byte{'0'}, 7)...),
		},
		{
			write: func(w *Writer) error {
				return w.WriteExt(0xd, bytesMarshaler(bytes.Repeat([]byte{'0'}, math.MaxUint8+1)))
			},
			data: append([]byte{tagExt16, 0x1, 0x0, 0xd}, bytes.Repeat([]byte{'0'}, math.MaxUint8+1)...),
		},
		{
			write: func(w *Writer) error {
				return w.WriteExt(0xd, bytesMarshaler(bytes.Repeat([]byte{'0'}, math.MaxUint16+1)))
			},
			data: append([]byte{tagExt32, 0x0, 0x1, 0x0, 0x0, 0xd}, bytes.Repeat([]byte{'0'}, math.MaxUint16+1)...),
		},
		// time
		{
			write: func(w *Writer) error {
				return w.WriteTime(time.Date(2017, time.September, 26, 13, 14, 15, 0, time.UTC))
			},
			data: []byte{tagExt8, 12, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0xca, 0x52, 0xa7},
		},
	}

	var buf bytes.Buffer
	for _, test := range tests {
		buf.Reset()
		err := test.write(NewWriter(&buf))
		if err != nil {
			t.Errorf("unexpected write error: %v", err)
		} else if !bytes.Equal(buf.Bytes(), test.data) {
			t.Errorf("unexpected data: %x", buf.Bytes())
		}
	}
}

type bytesMarshaler []byte

func (m bytesMarshaler) MarshalBinary() ([]byte, error) {
	return []byte(m), nil
}
