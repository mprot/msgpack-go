package msgpack

import (
	"io"
)

// Decoder defines an interface for types which are able to decode
// themselves from an MessagePack encoding.
type Decoder interface {
	DecodeMsgpack(r *Reader) error
}

// Decode decodes the MessagePack encoding provided by r into v.
func Decode(r io.Reader, v Decoder) error {
	reader := NewReader(r)
	err := v.DecodeMsgpack(reader)
	releaseReader(reader)
	return err
}
