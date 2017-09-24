package msgpack

import "io"

// Encoder defines an interface for types which are able to encode
// themselves into an MessagePack encoding.
type Encoder interface {
	EncodeMsgpack(w *Writer) error
}

// Encode encodes v into the MessagePack encoding and writes it to w.
func Encode(w io.Writer, v Encoder) error {
	return v.EncodeMsgpack(NewWriter(w))
}
