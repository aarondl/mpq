package mpq

import (
	"bytes"
	"testing"
)

func TestHETTable(t *testing.T) {
	setup()

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
	if m.HETTable.HashTableSize != 18 {
		t.Errorf("Incorrect Value for HashTableSize: %X", m.HETTable.HashTableSize)
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

	hashes := []byte{0x80, 0xDE, 0xCE, 0x00, 0x81, 0xE9, 0x00, 0x00, 0xC4, 0x84, 0x83, 0xBB, 0x00, 0xD7, 0xC6, 0x87, 0x89, 0xE4}
	indicies := []uint{0x5, 0x6, 0x9, 0xF, 0x0, 0xD, 0xF, 0xF, 0x4, 0x7, 0x3, 0xC, 0xF, 0x1, 0x2, 0xB, 0x8, 0xA}

	if bytes.Compare(m.HETTable.Hashes, hashes) != 0 {
		t.Errorf("Incorrect Value for Hashes: % 02X", m.HETTable.Hashes)
	}

	indexes, err := m.HETTable.Indexes()
	if err != nil {
		t.Error("Failed to read HET Table Indexes.")
	}
	for i, v := range indexes {
		if indicies[i] != v {
			t.Errorf("Incorrect Value for Index %d: %02X, expected: %02X", i, v, indicies[i])
		}
	}
}
