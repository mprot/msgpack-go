package msgpack

import (
	"bytes"
	"io"
)

// Encoder defines an interface for types which are able to encode
// themselves into an MessagePack encoding.
type Encoder interface {
	EncodeMsgpack(w *Writer) error
}

// Encode encodes v into the MessagePack encoding and writes it to w.
func Encode(w io.Writer, v Encoder) error {
	return v.EncodeMsgpack(NewWriter(w))
}

// Marshal encodes v into the MessagePack encoding and returns its encoding.
func Marshal(v Encoder) ([]byte, error) {
	var buf bytes.Buffer
	err := Encode(&buf, v)
	return buf.Bytes(), err
}
