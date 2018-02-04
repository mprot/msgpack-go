# msgpack-go
msgpack-go is a [MessagePack](http://msgpack.org/) implementation for the Go programming language.

## Encoding
To encode values into the binary MessagePack format, an `Encode` function is provided:
```Go
func Encode(w io.Writer, v Encoder) error
```
`Encoder` describes an encodable type which fulfills the interface
```Go
type Encoder interface {
	EncodeMsgpack(w *Writer) error
}
```
and writes its value to the given writer. The `Writer` type provides a `Write*` method for each supported data type. For a complete overview of the `Writer` type, see the [documentation](https://godoc.org/github.com/mprot/msgpack-go#Writer).

## Decoding
To decode binary MessagePack data into values, a `Decode` function is provided:
```Go
func Decode(r io.Reader, v Decoder) error
```
`Decoder` describes a decodable type which fulfills the interface
```Go
type Decoder interface {
	DecodeMsgpack(r *Reader) error
}
```
and reads the necessary values from the given reader. The `Reader` type provides a `Read*` method for each supported data type. Furthermore, a reader provides the following extra methods:
* [Peek](https://godoc.org/github.com/mprot/msgpack-go#Reader.Peek) for looking up the next type in the stream without moving the read pointer, and
* [Skip](https://godoc.org/github.com/mprot/msgpack-go#Reader.Skip) for skipping any value which comes next in the stream.

For a complete overview of the `Reader` type, see the [documentation](https://godoc.org/github.com/mprot/msgpack-go#Reader).

## Example
```Go
package main

import (
	"bytes"
	"log"

	msgpack "github.com/mprot/msgpack-go"
)

type Custom struct {
	Number  int
	Boolean bool
}

func (c *Custom) EncodeMsgpack(w *msgpack.Writer) error {
	if err := w.WriteInt(c.Number); err != nil {
		return err
	}
	return w.WriteBool(c.Boolean)
}

func (c *Custom) DecodeMsgpack(r *msgpack.Reader) (err error) {
	if c.Number, err = r.ReadInt(); err != nil {
		return err
	}
	c.Boolean, err = r.ReadBool()
	return err
}

func main() {
	var buf bytes.Buffer

	encoded := &Custom{Number: 7, Boolean: true}
	if err := msgpack.Encode(&buf, encoded); err != nil {
		log.Fatal(err)
	}

	decoded := &Custom{}
	if err := msgpack.Decode(&buf, decoded); err != nil {
		log.Fatal(err)
	}

	if decoded.Number != encoded.Number || decoded.Boolean != encoded.Boolean {
		log.Fatal("something went wrong")
	}
}

```
