package msgpack

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestReaderRead(t *testing.T) {
	tests := []struct {
		data  []byte
		value interface{}
		read  func(*Reader) (interface{}, error)
	}{
		// nil
		{
			data: []byte{tagNil},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadNil() },
		},
		// bool
		{
			data:  []byte{tagNil},
			value: false,
			read:  func(r *Reader) (interface{}, error) { return r.ReadBool() },
		},
		{
			data:  []byte{tagFalse},
			value: false,
			read:  func(r *Reader) (interface{}, error) { return r.ReadBool() },
		},
		{
			data:  []byte{tagTrue},
			value: true,
			read:  func(r *Reader) (interface{}, error) { return r.ReadBool() },
		},
		// int
		{
			data:  []byte{tagNil},
			value: int64(0),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{posFixintTag(7)},
			value: int64(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{negFixintTag(-7)},
			value: int64(-7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagInt8, 0xc0},
			value: int64(-64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagInt16, 0xff, 0xc0},
			value: int64(-64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagInt32, 0xff, 0xff, 0xff, 0xc0},
			value: int64(-64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagInt64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xc0},
			value: int64(-64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagUint8, 0x40},
			value: int64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagUint16, 0x0, 0x40},
			value: int64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagUint32, 0x0, 0x0, 0x0, 0x40},
			value: int64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagUint64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40},
			value: int64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data:  []byte{tagUint8, 0x07},
			value: int8(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt8() },
		},
		{
			data:  []byte{tagUint8, 0x07},
			value: int16(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt16() },
		},
		{
			data:  []byte{tagUint8, 0x07},
			value: int32(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt32() },
		},
		{
			data:  []byte{tagUint8, 0x07},
			value: int(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadInt() },
		},
		// uint
		{
			data:  []byte{tagNil},
			value: uint64(0),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{posFixintTag(7)},
			value: uint64(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagInt8, 0x40},
			value: uint64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagInt16, 0x0, 0x40},
			value: uint64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagInt32, 0x0, 0x0, 0x0, 0x40},
			value: uint64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagInt64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40},
			value: uint64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagUint8, 0x40},
			value: uint64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagUint16, 0x0, 0x40},
			value: uint64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagUint32, 0x0, 0x0, 0x0, 0x40},
			value: uint64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagUint64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x40},
			value: uint64(64),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data:  []byte{tagUint8, 0x07},
			value: uint8(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint8() },
		},
		{
			data:  []byte{tagUint8, 0x07},
			value: uint16(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint16() },
		},
		{
			data:  []byte{tagUint8, 0x07},
			value: uint32(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint32() },
		},
		{
			data:  []byte{tagUint8, 0x07},
			value: uint(7),
			read:  func(r *Reader) (interface{}, error) { return r.ReadUint() },
		},
		// float
		{
			data:  []byte{tagNil},
			value: float64(0),
			read:  func(r *Reader) (interface{}, error) { return r.ReadFloat64() },
		},
		{
			data:  []byte{tagFloat32, 0x40, 0x49, 0x0f, 0xd8},
			value: float64(float32(3.141592)),
			read:  func(r *Reader) (interface{}, error) { return r.ReadFloat64() },
		},
		{
			data:  []byte{tagFloat64, 0x40, 0x09, 0x21, 0xfa, 0xfc, 0x8b, 0x00, 0x7a},
			value: float64(3.141592),
			read:  func(r *Reader) (interface{}, error) { return r.ReadFloat64() },
		},
		{
			data:  []byte{tagFloat32, 0x40, 0x49, 0x0f, 0xd8},
			value: float32(3.141592),
			read:  func(r *Reader) (interface{}, error) { return r.ReadFloat32() },
		},
		// string/bytes
		{
			data:  []byte{tagNil},
			value: "",
			read:  func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		{
			data:  []byte{tagBin8, 0x03, 'f', 'o', 'o'},
			value: "foo",
			read:  func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		{
			data:  []byte{tagStr8, 0x03, 'f', 'o', 'o'},
			value: "foo",
			read:  func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		{
			data:  []byte{tagBin16, 0x0, 0x03, 'f', 'o', 'o'},
			value: "foo",
			read:  func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		{
			data:  []byte{tagStr16, 0x0, 0x03, 'f', 'o', 'o'},
			value: "foo",
			read:  func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		{
			data:  []byte{tagBin32, 0x0, 0x0, 0x0, 0x03, 'f', 'o', 'o'},
			value: "foo",
			read:  func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		{
			data:  []byte{tagStr32, 0x0, 0x0, 0x0, 0x03, 'f', 'o', 'o'},
			value: "foo",
			read:  func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		{
			data:  []byte{fixstrTag(3), 'f', 'o', 'o'},
			value: "foo",
			read:  func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		// array
		{
			data:  []byte{tagNil},
			value: 0,
			read:  func(r *Reader) (interface{}, error) { return r.ReadArrayHeader() },
		},
		{
			data:  []byte{fixarrayTag(3)},
			value: 3,
			read:  func(r *Reader) (interface{}, error) { return r.ReadArrayHeader() },
		},
		{
			data:  []byte{tagArray16, 0x0, 0x03},
			value: 3,
			read:  func(r *Reader) (interface{}, error) { return r.ReadArrayHeader() },
		},
		{
			data:  []byte{tagArray32, 0x0, 0x0, 0x0, 0x03},
			value: 3,
			read:  func(r *Reader) (interface{}, error) { return r.ReadArrayHeader() },
		},
		// array with size
		{
			data: []byte{tagNil},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadArrayHeaderWithSize(0) },
		},
		{
			data: []byte{fixarrayTag(3)},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadArrayHeaderWithSize(3) },
		},
		{
			data: []byte{tagArray16, 0x0, 0x03},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadArrayHeaderWithSize(3) },
		},
		{
			data: []byte{tagArray32, 0x0, 0x0, 0x0, 0x03},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadArrayHeaderWithSize(3) },
		},
		// map
		{
			data:  []byte{tagNil},
			value: 0,
			read:  func(r *Reader) (interface{}, error) { return r.ReadMapHeader() },
		},
		{
			data:  []byte{fixmapTag(3)},
			value: 3,
			read:  func(r *Reader) (interface{}, error) { return r.ReadMapHeader() },
		},
		{
			data:  []byte{tagMap16, 0x0, 0x03},
			value: 3,
			read:  func(r *Reader) (interface{}, error) { return r.ReadMapHeader() },
		},
		{
			data:  []byte{tagMap32, 0x0, 0x0, 0x0, 0x03},
			value: 3,
			read:  func(r *Reader) (interface{}, error) { return r.ReadMapHeader() },
		},
		// raw
		{
			data:  []byte{tagTrue},
			value: Raw{tagTrue},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagUint8, 0x00},
			value: Raw{tagUint8, 0x00},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagInt16, 0x00, 0x00},
			value: Raw{tagInt16, 0x00, 0x00},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagFloat32, 0x00, 0x00, 0x00, 0x00},
			value: Raw{tagFloat32, 0x00, 0x00, 0x00, 0x00},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagFloat64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			value: Raw{tagFloat64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagBin8, 0x03, ' ', ' ', ' '},
			value: Raw{tagBin8, 0x03, ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagBin16, 0x0, 0x03, ' ', ' ', ' '},
			value: Raw{tagBin16, 0x0, 0x03, ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagBin32, 0x0, 0x0, 0x0, 0x03, ' ', ' ', ' '},
			value: Raw{tagBin32, 0x0, 0x0, 0x0, 0x03, ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagArray16, 0x0, 0x03, tagNil, tagNil, tagNil},
			value: Raw{tagArray16, 0x0, 0x03, tagNil, tagNil, tagNil},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagArray32, 0x0, 0x0, 0x0, 0x03, tagNil, tagNil, tagNil},
			value: Raw{tagArray32, 0x0, 0x0, 0x0, 0x03, tagNil, tagNil, tagNil},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagMap16, 0x0, 0x03, tagNil, tagNil, tagNil, tagNil, tagNil, tagNil},
			value: Raw{tagMap16, 0x0, 0x03, tagNil, tagNil, tagNil, tagNil, tagNil, tagNil},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagMap32, 0x0, 0x0, 0x0, 0x03, tagNil, tagNil, tagNil, tagNil, tagNil, tagNil},
			value: Raw{tagMap32, 0x0, 0x0, 0x0, 0x03, tagNil, tagNil, tagNil, tagNil, tagNil, tagNil},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagFixExt1, 0x0d, ' '},
			value: Raw{tagFixExt1, 0x0d, ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagFixExt2, 0x0d, ' ', ' '},
			value: Raw{tagFixExt2, 0x0d, ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagFixExt4, 0x0d, ' ', ' ', ' ', ' '},
			value: Raw{tagFixExt4, 0x0d, ' ', ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagFixExt8, 0x0d, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			value: Raw{tagFixExt8, 0x0d, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagFixExt16, 0x0d, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			value: Raw{tagFixExt16, 0x0d, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagExt8, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' '},
			value: Raw{tagExt8, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagExt16, 0x0, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' '},
			value: Raw{tagExt16, 0x0, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{tagExt32, 0x0, 0x0, 0x0, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' '},
			value: Raw{tagExt32, 0x0, 0x0, 0x0, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{posFixintTag(7)},
			value: Raw{posFixintTag(7)},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{negFixintTag(-7)},
			value: Raw{negFixintTag(-7)},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{fixstrTag(3), ' ', ' ', ' '},
			value: Raw{fixstrTag(3), ' ', ' ', ' '},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{fixarrayTag(3), tagNil, tagNil, tagNil},
			value: Raw{fixarrayTag(3), tagNil, tagNil, tagNil},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		{
			data:  []byte{fixmapTag(3), tagNil, tagNil, tagNil, tagNil, tagNil, tagNil},
			value: Raw{fixmapTag(3), tagNil, tagNil, tagNil, tagNil, tagNil, tagNil},
			read:  func(r *Reader) (interface{}, error) { return r.ReadRaw(nil) },
		},
		// ext
		{
			data: []byte{tagFixExt1, 0xd, ' '},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0xd, nopBinaryMarshaler{}) },
		},
		{
			data: []byte{tagFixExt2, 0xd, ' ', ' '},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0xd, nopBinaryMarshaler{}) },
		},
		{
			data: []byte{tagFixExt4, 0xd, ' ', ' ', ' ', ' '},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0xd, nopBinaryMarshaler{}) },
		},
		{
			data: []byte{tagFixExt8, 0xd, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0xd, nopBinaryMarshaler{}) },
		},
		{
			data: []byte{tagFixExt16, 0xd, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0xd, nopBinaryMarshaler{}) },
		},
		{
			data: []byte{tagExt8, 0x03, 0xd, ' ', ' ', ' '},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0xd, nopBinaryMarshaler{}) },
		},
		{
			data: []byte{tagExt16, 0x0, 0x03, 0xd, ' ', ' ', ' '},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0xd, nopBinaryMarshaler{}) },
		},
		{
			data: []byte{tagExt32, 0x0, 0x0, 0x0, 0x03, 0xd, ' ', ' ', ' '},
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0xd, nopBinaryMarshaler{}) },
		},
		// time
		{
			data:  []byte{tagFixExt4, 0xff, 0x59, 0xca, 0x52, 0xa7},
			value: time.Date(2017, time.September, 26, 13, 14, 15, 0, time.UTC),
			read:  func(r *Reader) (interface{}, error) { return r.ReadTime() },
		},
		{
			data:  []byte{tagFixExt8, 0xff, 0x00, 0x00, 0x00, 0x40, 0x59, 0xca, 0x52, 0xa7},
			value: time.Date(2017, time.September, 26, 13, 14, 15, 16, time.UTC),
			read:  func(r *Reader) (interface{}, error) { return r.ReadTime() },
		},
		{
			data:  []byte{tagExt8, 12, 0xff, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x59, 0xca, 0x52, 0xa7},
			value: time.Date(2017, time.September, 26, 13, 14, 15, 16, time.UTC),
			read:  func(r *Reader) (interface{}, error) { return r.ReadTime() },
		},
	}

	for _, test := range tests {
		r := NewReader(bytes.NewReader(test.data))
		val, err := test.read(r)
		if err != nil {
			t.Errorf("unexpected read error: %v", err)
		} else if !reflect.DeepEqual(val, test.value) {
			t.Errorf("unexpected value for %s(%x): %v", tagType(test.data[0]), test.data, val)
		}

		if _, err = r.Peek(); err != io.EOF {
			t.Errorf("unexpected error for %s(%x): %v", tagType(test.data[0]), test.data, err)
		}
	}
}

func TestReaderPeek(t *testing.T) {
	tests := []struct {
		data []byte
		typ  Type
	}{
		// nil
		{
			data: []byte{tagNil},
			typ:  Nil,
		},
		// bool
		{
			data: []byte{tagFalse},
			typ:  Bool,
		},
		{
			data: []byte{tagTrue},
			typ:  Bool,
		},
		// int
		{
			data: []byte{tagInt8},
			typ:  Int,
		},
		{
			data: []byte{tagInt16},
			typ:  Int,
		},
		{
			data: []byte{tagInt32},
			typ:  Int,
		},
		{
			data: []byte{tagInt64},
			typ:  Int,
		},
		// uint
		{
			data: []byte{tagUint8},
			typ:  Uint,
		},
		{
			data: []byte{tagUint16},
			typ:  Uint,
		},
		{
			data: []byte{tagUint32},
			typ:  Uint,
		},
		{
			data: []byte{tagUint64},
			typ:  Uint,
		},
		// float
		{
			data: []byte{tagFloat32},
			typ:  Float,
		},
		{
			data: []byte{tagFloat64},
			typ:  Float,
		},
		// bytes
		{
			data: []byte{tagBin8},
			typ:  Bytes,
		},
		{
			data: []byte{tagBin16},
			typ:  Bytes,
		},
		{
			data: []byte{tagBin32},
			typ:  Bytes,
		},
		// string
		{
			data: []byte{tagStr8},
			typ:  String,
		},
		{
			data: []byte{tagStr16},
			typ:  String,
		},
		{
			data: []byte{tagStr32},
			typ:  String,
		},
		// array
		{
			data: []byte{tagArray16},
			typ:  Array,
		},
		{
			data: []byte{tagArray32},
			typ:  Array,
		},
		// map
		{
			data: []byte{tagMap16},
			typ:  Map,
		},
		{
			data: []byte{tagMap32},
			typ:  Map,
		},
		// ext
		{
			data: []byte{tagFixExt1, 0x0a},
			typ:  Ext,
		},
		{
			data: []byte{tagFixExt2, 0x0a},
			typ:  Ext,
		},
		{
			data: []byte{tagFixExt4, 0x0a},
			typ:  Ext,
		},
		{
			data: []byte{tagFixExt8, 0x0a},
			typ:  Ext,
		},
		{
			data: []byte{tagFixExt16, 0x0a},
			typ:  Ext,
		},
		{
			data: []byte{tagExt8, 0x0d, 0x0a},
			typ:  Ext,
		},
		{
			data: []byte{tagExt16, 0x00, 0x0d, 0x0a},
			typ:  Ext,
		},
		{
			data: []byte{tagExt32, 0x00, 0x00, 0x00, 0x0d, 0x0a},
			typ:  Ext,
		},
		// time
		{
			data: []byte{tagFixExt1, 0xff},
			typ:  Time,
		},
		{
			data: []byte{tagFixExt2, 0xff},
			typ:  Time,
		},
		{
			data: []byte{tagFixExt4, 0xff},
			typ:  Time,
		},
		{
			data: []byte{tagFixExt8, 0xff},
			typ:  Time,
		},
		{
			data: []byte{tagFixExt16, 0xff},
			typ:  Time,
		},
		{
			data: []byte{tagExt8, 0x0d, 0xff},
			typ:  Time,
		},
		{
			data: []byte{tagExt16, 0x00, 0x0d, 0xff},
			typ:  Time,
		},
		{
			data: []byte{tagExt32, 0x00, 0x00, 0x00, 0x0d, 0xff},
			typ:  Time,
		},
	}

	for _, test := range tests {
		r := NewReader(bytes.NewReader(test.data))
		typ, err := r.Peek()
		if err != nil {
			t.Errorf("unexpected peek error for %x: %v", test.data, err)
		} else if typ != test.typ {
			t.Errorf("unexpected type for %x: %s", test.data, typ)
		}
	}
}

func TestReaderSkip(t *testing.T) {
	stream := []byte{
		tagNil,
		tagFalse,
		tagTrue,
		tagInt8, 0x0,
		tagUint8, 0x0,
		tagInt16, 0x0, 0x0,
		tagUint16, 0x0, 0x0,
		tagInt32, 0x0, 0x0, 0x0, 0x0,
		tagUint32, 0x0, 0x0, 0x0, 0x0,
		tagInt64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		tagUint64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		tagStr8, 0x03, ' ', ' ', ' ',
		tagBin8, 0x03, ' ', ' ', ' ',
		tagStr16, 0x0, 0x03, ' ', ' ', ' ',
		tagBin16, 0x0, 0x03, ' ', ' ', ' ',
		tagStr32, 0x0, 0x0, 0x0, 0x03, ' ', ' ', ' ',
		tagBin32, 0x0, 0x0, 0x0, 0x03, ' ', ' ', ' ',
		tagArray16, 0x0, 0x03, tagNil, tagNil, tagNil,
		tagArray32, 0x0, 0x0, 0x0, 0x03, tagNil, tagNil, tagNil,
		tagMap16, 0x0, 0x03, tagNil, tagNil, tagNil, tagNil, tagNil, tagNil,
		tagMap32, 0x0, 0x0, 0x0, 0x03, tagNil, tagNil, tagNil, tagNil, tagNil, tagNil,
		tagFixExt1, 0x0d, ' ',
		tagFixExt2, 0x0d, ' ', ' ',
		tagFixExt4, 0x0d, ' ', ' ', ' ', ' ',
		tagFixExt8, 0x0d, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		tagFixExt16, 0x0d, ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		tagExt8, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' ',
		tagExt16, 0x0, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' ',
		tagExt32, 0x0, 0x0, 0x0, 0x05, 0x0d, ' ', ' ', ' ', ' ', ' ',
		posFixintTag(7),
		negFixintTag(-7),
		fixstrTag(3), ' ', ' ', ' ',
		fixarrayTag(3), tagNil, tagNil, tagNil,
		fixmapTag(3), tagNil, tagNil, tagNil, tagNil, tagNil, tagNil,
	}

	const n = 34
	r := NewReader(bytes.NewReader(stream))
	for i := 0; i < n; i++ {
		if err := r.Skip(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if err := r.Skip(); err != io.EOF {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReaderError(t *testing.T) {
	tests := []struct {
		data []byte
		err  string
		read func(*Reader) (interface{}, error)
	}{
		// nil
		{
			data: []byte{negFixintTag(-7)},
			err:  "unexpected type: int (expected nil)",
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadNil() },
		},
		// bool
		{
			data: []byte{negFixintTag(-7)},
			err:  "unexpected type: int (expected bool)",
			read: func(r *Reader) (interface{}, error) { return r.ReadBool() },
		},
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected bool)",
			read: func(r *Reader) (interface{}, error) { return r.ReadBool() },
		},
		// int8
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected int)",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt8() },
		},
		{
			data: []byte{tagInt16, 0x7f, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt8() },
		},
		// int16
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected int)",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt16() },
		},
		{
			data: []byte{tagInt32, 0x7f, 0xff, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt16() },
		},
		// int32
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected int)",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt32() },
		},
		{
			data: []byte{tagInt64, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt32() },
		},
		// int64
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected int)",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		{
			data: []byte{tagUint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		},
		// int
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected int)",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt() },
		},
		{
			data: []byte{tagUint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadInt() },
		},
		// uint8
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected uint)",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint8() },
		},
		{
			data: []byte{tagUint16, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint8() },
		},
		// uint16
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected uint)",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint16() },
		},
		{
			data: []byte{tagUint32, 0xff, 0xff, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint16() },
		},
		// uint32
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected uint)",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint32() },
		},
		{
			data: []byte{tagUint64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint32() },
		},
		// uint64
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected uint)",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data: []byte{tagInt8, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data: []byte{tagInt16, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data: []byte{tagInt32, 0xff, 0xff, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		{
			data: []byte{tagInt64, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:  "integer overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		},
		// uint
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected uint)",
			read: func(r *Reader) (interface{}, error) { return r.ReadUint() },
		},
		// float32
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected float)",
			read: func(r *Reader) (interface{}, error) { return r.ReadFloat32() },
		},
		{
			data: []byte{tagFloat64, 0x7f, 0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:  "floating-point overflow",
			read: func(r *Reader) (interface{}, error) { return r.ReadFloat32() },
		},
		// float64
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected float)",
			read: func(r *Reader) (interface{}, error) { return r.ReadFloat64() },
		},
		// bytes
		{
			data: []byte{tagInt8, 0x7f},
			err:  "unexpected type: int (expected bytes)",
			read: func(r *Reader) (interface{}, error) { return r.ReadBytes(nil) },
		},
		{
			data: []byte{tagBin8, 0x03, ' ', ' '},
			err:  io.ErrUnexpectedEOF.Error(),
			read: func(r *Reader) (interface{}, error) { return r.ReadBytes(nil) },
		},
		// string
		{
			data: []byte{tagInt8},
			err:  "unexpected type: int (expected string)",
			read: func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		{
			data: []byte{tagStr8, 0x03, ' ', ' '},
			err:  io.ErrUnexpectedEOF.Error(),
			read: func(r *Reader) (interface{}, error) { return r.ReadString() },
		},
		// array
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected array)",
			read: func(r *Reader) (interface{}, error) { return r.ReadArrayHeader() },
		},
		// array with size
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected array)",
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadArrayHeaderWithSize(2) },
		},
		{
			data: []byte{fixarrayTag(3)},
			err:  "invalid array header size 3",
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadArrayHeaderWithSize(2) },
		},
		// map
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected map)",
			read: func(r *Reader) (interface{}, error) { return r.ReadMapHeader() },
		},
		// ext
		{
			data: []byte{tagStr8},
			err:  "unexpected type: string (expected ext)",
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0x0d, nil) },
		},
		{
			data: []byte{tagFixExt1, 0x0e, ' '},
			err:  "invalid extension type 14",
			read: func(r *Reader) (interface{}, error) { return nil, r.ReadExt(0x0d, nil) },
		},
		// time
		{
			data: []byte{tagExt8, 16, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:  "invalid timestamp length 16",
			read: func(r *Reader) (interface{}, error) { return r.ReadTime() },
		},
	}

	for _, test := range tests {
		r := NewReader(bytes.NewReader(test.data))
		_, err := test.read(r)
		if err == nil {
			t.Error("expected error, got none")
		} else if err.Error() != test.err {
			t.Errorf("unexpected error message: %v", err)
		}
	}
}

type nopBinaryMarshaler struct{}

func (m nopBinaryMarshaler) UnmarshalBinary(p []byte) error {
	return nil
}
