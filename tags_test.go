package msgpack

import "testing"

func TestPosFixint(t *testing.T) {
	tag := posFixintTag(7)
	if tag != 0x07 {
		t.Errorf("unexpected tag: %02x", tag)
	}
	if !isPosFixintTag(tag) {
		t.Error("expected positive fix integer")
	}
	if x := readPosFixint(tag); x != 7 {
		t.Errorf("unexpected value: %d", x)
	}
}

func TestNegFixint(t *testing.T) {
	tag := negFixintTag(-1)
	if tag != 0xff {
		t.Errorf("unexpected tag: %02x", tag)
	}
	if !isNegFixintTag(tag) {
		t.Error("expected negative fix integer")
	}
	if x := readNegFixint(tag); x != -1 {
		t.Errorf("unexpected value: %d", x)
	}
}

func TestFixstr(t *testing.T) {
	tag := fixstrTag(7)
	if tag != 0xa7 {
		t.Errorf("unexpected tag: %02x", tag)
	}
	if !isFixstrTag(tag) {
		t.Error("expected fix string")
	}
	if x := readFixstr(tag); x != 7 {
		t.Errorf("unexpected value: %d", x)
	}
}

func TestFixarray(t *testing.T) {
	tag := fixarrayTag(7)
	if tag != 0x97 {
		t.Errorf("unexpected tag: %02x", tag)
	}
	if !isFixarrayTag(tag) {
		t.Error("expected fix array")
	}
	if x := readFixarray(tag); x != 7 {
		t.Errorf("unexpected value: %d", x)
	}
}

func TestFixmap(t *testing.T) {
	tag := fixmapTag(7)
	if tag != 0x87 {
		t.Errorf("unexpected tag: %02x", tag)
	}
	if !isFixmapTag(tag) {
		t.Error("expected fix map")
	}
	if x := readFixmap(tag); x != 7 {
		t.Errorf("unexpected value: %d", x)
	}
}
