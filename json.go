package msgpack

import (
	"bytes"
	"encoding/base64"
	"io"
	"strconv"
)

// CopyToJSON is a helper function for reading MessagePack encoded data from r,
// transforming it to JSON, and writing it to w.
func CopyToJSON(w io.Writer, r io.Reader) (written int, err error) {
	reader := NewReader(r)
	writer := jsonWriter{w}
	newline := [1]byte{'\n'}
	for {
		n, err := writer.WriteVal(reader, false)
		written += n
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return written, err
		}
		n, err = w.Write(newline[:])
		written += n
		if err != nil {
			return written, err
		}
	}
}

type jsonWriter struct {
	io.Writer
}

func (w *jsonWriter) WriteVal(r *Reader, quoted bool) (int, error) {
	typ, err := r.Peek()
	if err != nil {
		return 0, err
	}

	switch typ {
	case Nil:
		return w.writeNil(r, quoted)
	case Bool:
		return w.writeBool(r, quoted)
	case Int:
		return w.writeInt(r, quoted)
	case Uint:
		return w.writeUint(r, quoted)
	case Float:
		return w.writeFloat(r, quoted)
	case String:
		return w.writeString(r)
	case Bytes:
		return w.writeBytes(r)
	case Array:
		return w.writeArray(r, quoted)
	case Map:
		return w.writeMap(r, quoted)
	default:
		return 0, errorf("unsupported json type: %s", typ)
	}
}

func (w *jsonWriter) writeNil(r *Reader, quoted bool) (int, error) {
	if err := r.ReadNil(); err != nil {
		return 0, err
	}
	return w.puts("null", quoted)
}

func (w *jsonWriter) writeBool(r *Reader, quoted bool) (int, error) {
	b, err := r.ReadBool()
	if err != nil {
		return 0, err
	}

	s := "false"
	if b {
		s = "true"
	}
	return w.puts(s, quoted)
}

func (w *jsonWriter) writeInt(r *Reader, quoted bool) (int, error) {
	i, err := r.ReadInt64()
	if err != nil {
		return 0, err
	}
	return w.puts(strconv.FormatInt(i, 10), quoted)
}

func (w *jsonWriter) writeUint(r *Reader, quoted bool) (int, error) {
	ui, err := r.ReadUint64()
	if err != nil {
		return 0, err
	}
	return w.puts(strconv.FormatUint(ui, 10), quoted)
}

func (w *jsonWriter) writeFloat(r *Reader, quoted bool) (int, error) {
	f, err := r.ReadFloat64()
	if err != nil {
		return 0, err
	}
	return w.puts(strconv.FormatFloat(f, 'f', -1, 64), quoted)
}

func (w *jsonWriter) writeString(r *Reader) (int, error) {
	s, err := r.ReadString()
	if err != nil {
		return 0, err
	}
	return w.puts(s, true)
}

func (w *jsonWriter) writeBytes(r *Reader) (int, error) {
	b, err := r.ReadBytes(nil)
	if err != nil {
		return 0, err
	}

	var buf bytes.Buffer
	b64 := base64.NewEncoder(base64.StdEncoding, &buf)
	if _, err = b64.Write(b); err == nil {
		err = b64.Close()
	}
	if err != nil {
		return 0, err
	}
	return w.puts(buf.String(), true)
}

func (w *jsonWriter) writeArray(r *Reader, quoted bool) (int, error) {
	if quoted {
		return 0, errorString("cannot use array type as map key")
	}

	count, err := r.ReadArrayHeader()
	if err != nil {
		return 0, err
	}

	var n int

	if err = w.put('['); err != nil {
		return n, err
	}
	n++

	comma := false
	for i := 0; i < count; i++ {
		if comma {
			if err = w.put(','); err != nil {
				return n, err
			}
			n++
		}

		m, err := w.WriteVal(r, false)
		n += m
		if err != nil {
			return n, unexpectedEOF(err)
		}

		comma = true
	}

	if err = w.put(']'); err != nil {
		return n, err
	}
	n++
	return n, nil
}

func (w *jsonWriter) writeMap(r *Reader, quoted bool) (int, error) {
	if quoted {
		return 0, errorString("cannot use map type as map key")
	}

	count, err := r.ReadMapHeader()
	if err != nil {
		return 0, err
	}

	var n int

	if err = w.put('{'); err != nil {
		return n, err
	}
	n++

	comma := false
	for i := 0; i < count; i++ {
		if comma {
			if err = w.put(','); err != nil {
				return n, err
			}
			n++
		}

		m, err := w.WriteVal(r, true)
		n += m
		if err != nil {
			return n, unexpectedEOF(err)
		}

		if err = w.put(':'); err != nil {
			return n, err
		}
		n++

		m, err = w.WriteVal(r, false)
		n += m
		if err != nil {
			return n, unexpectedEOF(err)
		}

		comma = true
	}

	if err = w.put('}'); err != nil {
		return n, err
	}
	n++
	return n, nil
}

func (w *jsonWriter) put(c byte) error {
	if bw, ok := w.Writer.(io.ByteWriter); ok {
		return bw.WriteByte(c)
	}

	buf := [1]byte{c}
	_, err := w.Write(buf[:])
	return err
}

func (w *jsonWriter) puts(s string, quoted bool) (int, error) {
	if quoted {
		s = strconv.Quote(s)
	}
	return io.WriteString(w, s)
}

func unexpectedEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}
