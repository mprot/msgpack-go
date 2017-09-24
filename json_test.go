package msgpack

import (
	"bytes"
	"strings"
	"testing"
)

func TestCopyToJSON(t *testing.T) {
	fatalErr := func(err error) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	data := []struct {
		json  string
		write func(*Writer)
	}{
		{`null`, func(w *Writer) { fatalErr(w.WriteNil()) }},
		{`true`, func(w *Writer) { fatalErr(w.WriteBool(true)) }},
		{`-13`, func(w *Writer) { fatalErr(w.WriteInt(-13)) }},
		{`13`, func(w *Writer) { fatalErr(w.WriteUint(13)) }},
		{`3.141592`, func(w *Writer) { fatalErr(w.WriteFloat64(3.141592)) }},
		{`"YmxvYg=="`, func(w *Writer) { fatalErr(w.WriteBytes([]byte("blob"))) }},
		{`"two\nlines"`, func(w *Writer) { fatalErr(w.WriteString("two\nlines")) }},

		{`[1,2,3]`, func(w *Writer) {
			fatalErr(w.WriteArrayHeader(3))
			fatalErr(w.WriteInt(1))
			fatalErr(w.WriteInt(2))
			fatalErr(w.WriteInt(3))
		}},

		{`{"3141592":"pi"}`, func(w *Writer) {
			fatalErr(w.WriteMapHeader(1))
			fatalErr(w.WriteInt(3141592))
			fatalErr(w.WriteString("pi"))
		}},
	}

	var encoded, json bytes.Buffer
	for _, d := range data {
		encoded.Reset()
		json.Reset()
		d.write(NewWriter(&encoded))

		n, err := CopyToJSON(&json, &encoded)
		if err != nil {
			t.Errorf("unexpected copy error: %v", err)
		} else if res := strings.TrimSuffix(json.String(), "\n"); res != d.json {
			t.Errorf("unexpected json result for %s: %s", d.json, res)
		} else if n != json.Len() {
			t.Errorf("unexpected number of bytes: %d", n)
		}
	}
}
