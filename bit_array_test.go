package mpq

import (
	"bytes"
	"testing"
)

func toBinInt(s string) int64 {
	var val int64
	var offset uint
	for i := len(s) - 1; i >= 0; i-- {
		switch s[i] {
		case '1':
			val |= 1 << offset
			offset++
		case '0':
			offset++
		case ' ':
		default:
			panic("Not a valid binary number.")
		}
	}

	return val
}

func toBin(s string) byte {
	return byte(toBinInt(s))
}

func TestBitArray_Next(t *testing.T) {
	data := []byte{toBin("0000 1111"), toBin("1010 0101"), toBin("1111 0000")}

	b := newBitArray(data)
	var val int64
	var err error

	if val, err = b.next(5); err != nil {
		t.Error("Unexpected Error:", err)
	} else if val != toBinInt("01111") {
		t.Errorf("Wrong Value: % 02X", val)
	}
	if val, err = b.next(6); err != nil {
		t.Error("Unexpected Error:", err)
	} else if val != toBinInt("101000") {
		t.Errorf("Wrong Value: % 02X", val)
	}
	if val, err = b.next(1); err != nil {
		t.Error("Unexpected Error:", err)
	} else if val != toBinInt("0") {
		t.Errorf("Wrong Value: % 02X", val)
	}
	if val, err = b.next(3); err != nil {
		t.Error("Unexpected Error:", err)
	} else if val != toBinInt("010") {
		t.Errorf("Wrong Value: % 02X", val)
	}
	if val, err = b.next(2); err != nil {
		t.Error("Unexpected Error:", err)
	} else if val != toBinInt("01") {
		t.Errorf("Wrong Value: % 02X", val)
	}
	if val, err = b.next(7); err != nil {
		t.Error("Unexpected Error:", err)
	} else if val != toBinInt("1111 000") {
		t.Errorf("Wrong Value: % 02X", val)
	}
	if val, err = b.next(4); err == nil {
		t.Error("Expected an error to occur on overflow.")
	} else if val != 0 {
		t.Error("Expected 0 value.")
	}
}

func TestBitArray_NextBytes(t *testing.T) {
	data := []byte{0x00, 0xF0, 0xFF, 0x0F, 0x00}

	b := newBitArray(data)
	var err error

	val := make([]byte, 2)
	if err = b.nextBytes(val, 12); err != nil {
		t.Error("Unexpected Error:", err)
	} else if bytes.Compare([]byte{0x00, 0x00}, val) != 0 {
		t.Errorf("Wrong Value: % 02X", val)
	}

	val = make([]byte, 3)
	if err = b.nextBytes(val, 22); err != nil {
		t.Error("Unexpected Error:", err)
	} else if bytes.Compare([]byte{0xFF, 0xFF, 0x00}, val) != 0 {
		t.Errorf("Wrong Value: % 02X", val)
	}

	val = make([]byte, 1)
	if err = b.nextBytes(val, 6); err != nil {
		t.Error("Unexpected Error:", err)
	} else if bytes.Compare([]byte{0x00}, val) != 0 {
		t.Errorf("Wrong Value: % 02X", val)
	}

	if err = b.nextBytes(val, 6); err == nil {
		t.Error("Expected an overflow error.")
	}
}
