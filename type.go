package msgpack

import "fmt"

const (
	extTime = -1
)

// Type represents an MessagePack data type.
type Type string

// All supported MessagePack types.
const (
	Nil    Type = "nil"
	Bool   Type = "bool"
	Int    Type = "int"
	Uint   Type = "uint"
	Float  Type = "float"
	String Type = "string"
	Bytes  Type = "bytes"
	Array  Type = "array"
	Map    Type = "map"
	Ext    Type = "ext"

	// extension types
	Time Type = "time"
)

func extType(typ int8) Type {
	if typ == extTime {
		return Time
	}
	return Ext
}

// tagType returns the type for a given tag. For every extension type
// Ext will be returned.
func tagType(tag byte) Type {
	switch tag {
	case tagNil:
		return Nil
	case tagFalse, tagTrue:
		return Bool
	case tagInt8, tagInt16, tagInt32, tagInt64:
		return Int
	case tagUint8, tagUint16, tagUint32, tagUint64:
		return Uint
	case tagFloat32, tagFloat64:
		return Float
	case tagBin8, tagBin16, tagBin32:
		return Bytes
	case tagStr8, tagStr16, tagStr32:
		return String
	case tagArray16, tagArray32:
		return Array
	case tagMap16, tagMap32:
		return Map
	case tagFixExt1, tagFixExt2, tagFixExt4, tagFixExt8, tagFixExt16, tagExt8, tagExt16, tagExt32:
		return Ext
	}

	switch {
	case isPosFixintTag(tag):
		return Uint
	case isNegFixintTag(tag):
		return Int
	case isFixstrTag(tag):
		return String
	case isFixarrayTag(tag):
		return Array
	case isFixmapTag(tag):
		return Map
	}

	return Type(fmt.Sprintf("%#02x", tag))
}
