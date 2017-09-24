package msgpack

import "testing"

func TestTagType(t *testing.T) {
	tests := []struct {
		typ Type
		tag byte
	}{
		{Nil, tagNil},
		{Bool, tagTrue},
		{Bool, tagFalse},
		{Int, negFixintTag(-7)},
		{Int, tagInt8},
		{Int, tagInt16},
		{Int, tagInt32},
		{Int, tagInt64},
		{Uint, posFixintTag(7)},
		{Uint, tagUint8},
		{Uint, tagUint16},
		{Uint, tagUint32},
		{Uint, tagUint64},
		{Float, tagFloat32},
		{Float, tagFloat64},
		{Bytes, tagBin8},
		{Bytes, tagBin16},
		{Bytes, tagBin32},
		{String, fixstrTag(7)},
		{String, tagStr8},
		{String, tagStr16},
		{String, tagStr32},
		{Array, fixarrayTag(7)},
		{Array, tagArray16},
		{Array, tagArray32},
		{Map, fixmapTag(7)},
		{Map, tagMap16},
		{Map, tagMap32},
		{Ext, tagFixExt1},
		{Ext, tagFixExt2},
		{Ext, tagFixExt4},
		{Ext, tagFixExt8},
		{Ext, tagFixExt16},
		{Ext, tagExt8},
		{Ext, tagExt16},
		{Ext, tagExt32},
	}

	for _, test := range tests {
		typ := tagType(test.tag)
		if typ != test.typ {
			t.Errorf("unexpected type for tag %02x: %s (expected %s)", test.tag, typ, test.typ)
		}
	}
}
