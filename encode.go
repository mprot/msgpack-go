package msgpack

import (
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
	return AppendMarshal(v, nil)
}

// AppendMarshal encodes v into the MessagePack encoding and appends it to buf.
func AppendMarshal(v Encoder, buf []byte) ([]byte, error) {
	appender := &byteAppender{buf: buf}
	err := Encode(appender, v)
	return appender.buf, err
}

type byteAppender struct {
	buf []byte
}

func (a *byteAppender) Write(p []byte) (int, error) {
	a.buf = append(a.buf, p...)
	return len(p), nil
}
