package msgpack

// Raw is a raw encoded MessagePack value. It can be used to delay
// MessagePack decoding.
type Raw []byte

// EncodeMsgpack writes the raw MessagePack value into w.
func (raw Raw) EncodeMsgpack(w *Writer) error {
	return w.WriteRaw(raw)
}

// DecodeMsgpack reads the raw MessagePack value from r.
func (raw *Raw) DecodeMsgpack(r *Reader) error {
	var err error
	*raw, err = r.ReadRaw(*raw)
	return err
}
