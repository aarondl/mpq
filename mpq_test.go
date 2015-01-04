package mpq

import (
	"bytes"
	"os"
	"testing"
)

var m *MPQ

func TestMain(main *testing.M) {
	code := main.Run()
	if m != nil {
		m.Close()
	}
	os.Exit(code)
}

func setup() {
	var err error
	m, err = Open("Garden of Terror (72).StormReplay")
	if err != nil {
		panic("Could not open test replay: " + err.Error())
	}
}

func TestMPQHeader(t *testing.T) {
	setup()

	if m.Header.HeaderSize != 208 {
		t.Errorf("Incorrect Value for HeaderSize: %d", m.Header.HeaderSize)
	}
	if m.Header.Size != 0x6BFDA {
		t.Errorf("Incorrect Value for Size: %X", m.Header.Size)
	}
	if m.Header.FormatVersion != 3 {
		t.Errorf("Incorrect Value for FormatVersion: %d", m.Header.FormatVersion)
	}
	if m.Header.BlockSize != 5 {
		t.Errorf("Incorrect Value for BlockSize: %d", m.Header.BlockSize)
	}
	if m.Header.HashTablePos != 0x6BCFA {
		t.Errorf("Incorrect Value for HashTablePos: %X", m.Header.HashTablePos)
	}
	if m.Header.BlockTablePos != 0x6BEFA {
		t.Errorf("Incorrect Value for BlockTablePos: %X", m.Header.BlockTablePos)
	}
	if m.Header.HashTableSize != 32 {
		t.Errorf("Incorrect Value for HashTableSize: %d", m.Header.HashTableSize)
	}
	if m.Header.BlockTableSize != 14 {
		t.Errorf("Incorrect Value for BlockTableSize: %d", m.Header.BlockTableSize)
	}
	if m.Header.HiBlockTablePos != 0 {
		t.Errorf("Incorrect Value for HiBlockTablePos: %d", m.Header.HiBlockTablePos)
	}
	if m.Header.HashTablePosHi != 0 {
		t.Errorf("Incorrect Value for HashTablePosHi: %d", m.Header.HashTablePosHi)
	}
	if m.Header.BlockTablePosHi != 0 {
		t.Errorf("Incorrect Value for BlockTablePosHi: %d", m.Header.BlockTablePosHi)
	}
	if m.Header.ArchiveSize != 0x6BFDA {
		t.Errorf("Incorrect Value for ArchiveSize: %X", m.Header.ArchiveSize)
	}
	if m.Header.BETTablePos != 0x6BBBF {
		t.Errorf("Incorrect Value for BETTablePos: %X", m.Header.BETTablePos)
	}
	if m.Header.HETTablePos != 0x6BB68 {
		t.Errorf("Incorrect Value for HETTablePos: %X", m.Header.HETTablePos)
	}
	if m.Header.HashTableSize64 != 0x200 {
		t.Errorf("Incorrect Value for HashTableSize64: %X", m.Header.HashTableSize64)
	}
	if m.Header.BlockTableSize64 != 0xE0 {
		t.Errorf("Incorrect Value for BlockTableSize64: %X", m.Header.BlockTableSize64)
	}
	if m.Header.HiBlockTableSize64 != 0 {
		t.Errorf("Incorrect Value for HiBlockTableSize64: %X", m.Header.HiBlockTableSize64)
	}
	if m.Header.HETTableSize64 != 0x47 {
		t.Errorf("Incorrect Value for HETTableSize64: %X", m.Header.HETTableSize64)
	}
	if m.Header.BETTableSize64 != 0x12B {
		t.Errorf("Incorrect Value for BETTableSize64: %X", m.Header.BETTableSize64)
	}
	if m.Header.ChunkSize != 0x4000 {
		t.Errorf("Incorrect Value for ChunkSize: %X", m.Header.ChunkSize)
	}

	blockmd5 := []byte{0x59, 0xD6, 0x31, 0x75, 0xB2, 0xA1, 0x03, 0xEB, 0x5E, 0x81, 0x0C, 0x12, 0x30, 0xA4, 0x4D, 0x81}
	hashmd5 := []byte{0x13, 0x97, 0x99, 0x88, 0x5A, 0x78, 0xE5, 0xCB, 0x5C, 0x87, 0x07, 0x20, 0x24, 0x17, 0x7D, 0x70}
	hiblkmd5 := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	betmd5 := []byte{0x71, 0x37, 0x33, 0x51, 0x98, 0xC7, 0xF8, 0xAC, 0x51, 0xC0, 0x5D, 0xA5, 0x6F, 0x5F, 0x54, 0xA8}
	hetmd5 := []byte{0xC9, 0x01, 0x18, 0x4C, 0xE2, 0x06, 0x92, 0x91, 0xE2, 0xAD, 0x04, 0xE9, 0x6A, 0xA5, 0x34, 0xFF}
	headermd5 := []byte{0xE7, 0xAE, 0x9E, 0x92, 0xD2, 0x8F, 0xF3, 0x57, 0x27, 0x24, 0x01, 0xB6, 0x53, 0xD1, 0xDE, 0x3F}

	if bytes.Compare(m.Header.BlockTableMD5, blockmd5) != 0 {
		t.Errorf("\nExpected: % 02X\nGot     : % 02X", blockmd5, m.Header.BlockTableMD5)
	}
	if bytes.Compare(m.Header.HashTableMD5, hashmd5) != 0 {
		t.Errorf("\nExpected: % 02X\nGot     : % 02X", hashmd5, m.Header.HashTableMD5)
	}
	if bytes.Compare(m.Header.HiBlockTableMD5, hiblkmd5) != 0 {
		t.Errorf("\nExpected: % 02X\nGot     : % 02X", hiblkmd5, m.Header.HiBlockTableMD5)
	}
	if bytes.Compare(m.Header.BETTableMD5, betmd5) != 0 {
		t.Errorf("\nExpected: % 02X\nGot     : % 02X", betmd5, m.Header.BETTableMD5)
	}
	if bytes.Compare(m.Header.HETTableMD5, hetmd5) != 0 {
		t.Errorf("\nExpected: % 02X\nGot     : % 02X", hetmd5, m.Header.HETTableMD5)
	}
	if bytes.Compare(m.Header.MPQHeaderMD5, headermd5) != 0 {
		t.Errorf("\nExpected: % 02X\nGot     : % 02X", headermd5, m.Header.MPQHeaderMD5)
	}
}
