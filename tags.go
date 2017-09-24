package msgpack

const (
	tagNil = 0xc0 // 11000000

	// bool
	tagFalse = 0xc2 // 11000010
	tagTrue  = 0xc3 // 11000011

	// int
	tagInt8  = 0xd0 // 11010000
	tagInt16 = 0xd1 // 11010001
	tagInt32 = 0xd2 // 11010010
	tagInt64 = 0xd3 // 11010011

	// uint
	tagUint8  = 0xcc // 11001100
	tagUint16 = 0xcd // 11001101
	tagUint32 = 0xce // 11001110
	tagUint64 = 0xcf // 11001111

	// float
	tagFloat32 = 0xca // 11001010
	tagFloat64 = 0xcb // 11001011

	// string
	tagStr8  = 0xd9 // 11011001
	tagStr16 = 0xda // 11011010
	tagStr32 = 0xdb // 11011011

	// binary
	tagBin8  = 0xc4 // 11000100
	tagBin16 = 0xc5 // 11000101
	tagBin32 = 0xc6 // 11000110

	// array
	tagArray16 = 0xdc // 11011100
	tagArray32 = 0xdd // 11011101

	// map
	tagMap16 = 0xde // 11011110
	tagMap32 = 0xdf // 11011111

	// ext
	tagExt8     = 0xc7 // 11000111
	tagExt16    = 0xc8 // 11001000
	tagExt32    = 0xc9 // 11001001
	tagFixExt1  = 0xd4 // 11010100
	tagFixExt2  = 0xd5 // 11010101
	tagFixExt4  = 0xd6 // 11010110
	tagFixExt8  = 0xd7 // 11010111
	tagFixExt16 = 0xd8 // 11011000
)

// positive fixint: 0xxx xxxx
func posFixintTag(i uint8) byte {
	return i & 0x7f
}

func isPosFixintTag(tag byte) bool {
	return (tag & 0x80) == 0
}

func readPosFixint(tag byte) uint8 {
	return tag & 0x7f
}

// negative fixint: 111x xxxx
func negFixintTag(i int8) byte {
	return 0xe0 | (uint8(i) & 0x1f)
}

func isNegFixintTag(tag byte) bool {
	return (tag & 0xe0) == 0xe0
}

func readNegFixint(tag byte) int8 {
	return int8(tag)
}

// fixstr: 101x xxxx
func fixstrTag(length int) byte {
	return 0xa0 | (uint8(length) & 0x1f)
}

func isFixstrTag(tag byte) bool {
	return (tag & 0xe0) == 0xa0
}

func readFixstr(tag byte) uint8 {
	return tag & 0x1f
}

// fixarray: 1001 xxxx
func fixarrayTag(length int) byte {
	return 0x90 | (uint8(length) & 0x0f)
}

func isFixarrayTag(tag byte) bool {
	return (tag & 0xf0) == 0x90
}

func readFixarray(tag byte) uint8 {
	return tag & 0x0f
}

// fixmap: 1000 xxxx
func fixmapTag(length int) byte {
	return 0x80 | (uint8(length) & 0x0f)
}

func isFixmapTag(tag byte) bool {
	return (tag & 0xf0) == 0x80
}

func readFixmap(tag byte) uint8 {
	return tag & 0x0f
}
