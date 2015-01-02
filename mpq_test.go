package mpq

import (
	"bytes"
	"testing"
)

var m *MPQ

func init() {
	var err error
	m, err = Open("Garden of Terror (72).StormReplay")
	if err != nil {
		panic("Could not open test replay: " + err.Error())
	}
	defer m.Close()
}

func TestMPQHeader(t *testing.T) {
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

func TestUserData(t *testing.T) {
	if m.UserData.MaxSize != 512 {
		t.Errorf("Incorrect Value for MaxSize: %d", m.UserData.MaxSize)
	}
	if m.UserData.HeaderOffset != 1024 {
		t.Errorf("Incorrect Value for HeaderOffset: %d", m.UserData.HeaderOffset)
	}
	if m.UserData.UserDataHeaderSize != 97 {
		t.Errorf("Incorrect Value for UserDataHeaderSize: %d", m.UserData.UserDataHeaderSize)
	}

	userData := []byte{
		0x05, 0x0E, 0x00, 0x02, 0x3A, 0x48, 0x65, 0x72, 0x6F, 0x65, 0x73, 0x20, 0x6F, 0x66, 0x20, 0x74,
		0x68, 0x65, 0x20, 0x53, 0x74, 0x6F, 0x72, 0x6D, 0x20, 0x72, 0x65, 0x70, 0x6C, 0x61, 0x79, 0x1B,
		0x31, 0x31, 0x02, 0x05, 0x0C, 0x00, 0x09, 0x02, 0x02, 0x09, 0x00, 0x04, 0x09, 0x0E, 0x06, 0x09,
		0x04, 0x08, 0x09, 0x92, 0x89, 0x04, 0x0A, 0x09, 0x92, 0x89, 0x04, 0x04, 0x09, 0x04, 0x06, 0x09,
		0xB8, 0xD7, 0x02, 0x08, 0x06, 0x00, 0x0A, 0x05, 0x02, 0x02, 0x02, 0x20, 0x54, 0x82, 0x5A, 0xE1,
		0x9E, 0x9C, 0xA9, 0xA0, 0x8A, 0x25, 0xDD, 0xA5, 0x95, 0xA5, 0xA8, 0xAE, 0x0C, 0x09, 0x92, 0x89,
		0x04,
	}

	if bytes.Compare(userData, m.UserData.Data[:len(userData)]) != 0 {
		t.Error("User data contains wrong value.")
	}
	for i := len(userData); i < m.UserData.MaxSize; i++ {
		if m.UserData.Data[i] != 0 {
			t.Error("Expected the rest of user data to be blank.")
		}
	}
}

func TestHETTable(t *testing.T) {
	if m.HETTable.Version != 1 {
		t.Errorf("Incorrect Value for Version: %d", m.HETTable.Version)
	}
	if m.HETTable.DataSize != 0x3B {
		t.Errorf("Incorrect Value for DataSize: %X", m.HETTable.DataSize)
	}
	if m.HETTable.TableSize != 0x3B {
		t.Errorf("Incorrect Value for TableSize: %X", m.HETTable.TableSize)
	}
	if m.HETTable.EntryCount != 14 {
		t.Errorf("Incorrect Value for EntryCount: %d", m.HETTable.EntryCount)
	}
	if m.HETTable.TotalCount != 18 {
		t.Errorf("Incorrect Value for TotalCount: %X", m.HETTable.TotalCount)
	}
	if m.HETTable.HashEntrySize != 0x40 {
		t.Errorf("Incorrect Value for HashEntrySize: %X", m.HETTable.HashEntrySize)
	}
	if m.HETTable.IndexSizeTotal != 4 {
		t.Errorf("Incorrect Value for TotalIndexSize: %d", m.HETTable.IndexSizeTotal)
	}
	if m.HETTable.IndexSizeExtra != 0 {
		t.Errorf("Incorrect Value for IndexSizeExtra: %d", m.HETTable.IndexSizeExtra)
	}
	if m.HETTable.IndexSize != 4 {
		t.Errorf("Incorrect Value for IndexSize: %d", m.HETTable.IndexSize)
	}
	if m.HETTable.BlockTableSize != 9 {
		t.Errorf("Incorrect Value for BlockTableSize: %d", m.HETTable.BlockTableSize)
	}
}

func TestBETTable(t *testing.T) {
	if m.BETTable.Version != 1 {
		t.Errorf("Incorrect Value for Version: %d", m.BETTable.Version)
	}
	if m.BETTable.DataSize != 0x11F {
		t.Errorf("Incorrect Value for DataSize: %X", m.BETTable.DataSize)
	}
	if m.BETTable.TableSize != 0x11F {
		t.Errorf("Incorrect Value for TableSize: %X", m.BETTable.TableSize)
	}
	if m.BETTable.EntryCount != 14 {
		t.Errorf("Incorrect Value for EntryCount: %d", m.BETTable.EntryCount)
	}
	if m.BETTable.Unknown08 != 0x10 {
		t.Errorf("Incorrect Value for Unknown08: %X", m.BETTable.Unknown08)
	}
	if m.BETTable.TableEntrySize != 60 {
		t.Errorf("Incorrect Value for TableEntrySize: %d", m.BETTable.TableEntrySize)
	}
	if m.BETTable.BitIndexFilePos != 0 {
		t.Errorf("Incorrect Value for BitIndexFilePos: %d", m.BETTable.BitIndexFilePos)
	}
	if m.BETTable.BitIndexFileSize != 19 {
		t.Errorf("Incorrect Value for BitIndexFileSize: %d", m.BETTable.BitIndexFileSize)
	}
	if m.BETTable.BitIndexCmpSize != 39 {
		t.Errorf("Incorrect Value for BitIndexCmpSize: %d", m.BETTable.BitIndexCmpSize)
	}
	if m.BETTable.BitIndexFlagIndex != 58 {
		t.Errorf("Incorrect Value for BitIndexFlagIndex: %d", m.BETTable.BitIndexFlagIndex)
	}
	if m.BETTable.BitIndexUnknown != 0x3C {
		t.Errorf("Incorrect Value for BitIndexUnknown: %X", m.BETTable.BitIndexUnknown)
	}
	if m.BETTable.BitCountFilePos != 19 {
		t.Errorf("Incorrect Value for BitCountFilePos: %d", m.BETTable.BitCountFilePos)
	}
	if m.BETTable.BitCountFileSize != 20 {
		t.Errorf("Incorrect Value for BitCountFileSize: %d", m.BETTable.BitCountFileSize)
	}
	if m.BETTable.BitCountCmpSize != 19 {
		t.Errorf("Incorrect Value for BitCountCmpSize: %d", m.BETTable.BitCountCmpSize)
	}
	if m.BETTable.BitCountFlagIndex != 2 {
		t.Errorf("Incorrect Value for BitCountFlagIndex: %d", m.BETTable.BitCountFlagIndex)
	}
	if m.BETTable.BitCountUnknown != 0 {
		t.Errorf("Incorrect Value for BitCountUnknown: %d", m.BETTable.BitCountUnknown)
	}
	if m.BETTable.HashSizeTotal != 56 {
		t.Errorf("Incorrect Value for HashSizeTotal: %d", m.BETTable.HashSizeTotal)
	}
	if m.BETTable.HashSizeExtra != 0 {
		t.Errorf("Incorrect Value for HashSizeExtra: %d", m.BETTable.HashSizeExtra)
	}
	if m.BETTable.HashSize != 56 {
		t.Errorf("Incorrect Value for HashSize: %d", m.BETTable.HashSize)
	}
	if m.BETTable.HashArraySize != 98 {
		t.Errorf("Incorrect Value for HashArraySize: %d", m.BETTable.HashArraySize)
	}
	if m.BETTable.FlagCount != 2 {
		t.Errorf("Incorrect Value for FlagCount: %d", m.BETTable.FlagCount)
	}
}
