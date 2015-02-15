package mpq

import (
	"encoding/binary"
	"io"
)

// BlockTable is the older style BETTable in the MPQ Header.
type BlockTable struct {
	EntryCount int
	Table      []byte

	entries []BlockTableEntry
}

// BlockTableEntry describes the attributes of a file in the BlockTable.
type BlockTableEntry struct {
	FilePosition   uint32
	CompressedSize uint32
	FileSize       uint32
	Flags          uint32
}

func (m *MPQ) readBlockTable(r io.Reader) error {
	b := &BlockTable{EntryCount: m.Header.BlockTableSize}

	size := uint64(m.Header.BlockTableSize)
	compressedSize := m.Header.BlockTableSize64

	var err error
	if b.Table, err = decryptDecompressTable(r, size, compressedSize, cryptKeyBlockTable); err != nil {
		return err
	}

	m.BlockTable = b
	return nil
}

// Entries retrieves all the hash table entries.
func (b *BlockTable) Entries() []BlockTableEntry {
	if b.entries != nil {
		return b.entries
	}
	offset := 0

	entries := make([]BlockTableEntry, b.EntryCount)
	for i := 0; i < b.EntryCount; i++ {
		entry := &entries[i]

		entry.FilePosition = binary.LittleEndian.Uint32(b.Table[offset : offset+4])
		offset += 4
		entry.CompressedSize = binary.LittleEndian.Uint32(b.Table[offset : offset+4])
		offset += 4
		entry.FileSize = binary.LittleEndian.Uint32(b.Table[offset : offset+4])
		offset += 4
		entry.Flags = binary.LittleEndian.Uint32(b.Table[offset : offset+4])
		offset += 4
	}

	b.entries = entries
	return entries
}
