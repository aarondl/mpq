package mpq

import (
	"encoding/binary"
	"io"
)

const (
	hashTableEmpty   = 0xFFFFFFFF
	hashTableDeleted = 0xFFFFFFFE

	hashTableEntrySize = 16
)

type HashTable struct {
	EntryCount int
	Table      []byte
}

type HashTableEntry struct {
	Name1 uint32
	Name2 uint32

	Locale   uint16
	Platform uint16

	BlockIndex uint32
}

func (m *MPQ) readHashTable(r io.Reader) error {
	h := &HashTable{EntryCount: m.Header.HashTableSize}

	size := uint64(m.Header.HashTableSize)
	compressedSize := m.Header.HashTableSize64

	var err error
	if h.Table, err = decryptDecompressTable(r, size, compressedSize, cryptKeyHashTable); err != nil {
		return err
	}

	m.HashTable = h
	return nil
}

// Entries retrieves all the hash table entries.
func (h *HashTable) Entries() []HashTableEntry {
	offset := 0

	entries := make([]HashTableEntry, h.EntryCount)
	for i := 0; i < h.EntryCount; i++ {
		entry := &entries[i]

		entry.Name1 = binary.LittleEndian.Uint32(h.Table[offset : offset+4])
		offset += 4
		entry.Name2 = binary.LittleEndian.Uint32(h.Table[offset : offset+4])
		offset += 4

		entry.Locale = binary.LittleEndian.Uint16(h.Table[offset : offset+2])
		offset += 2
		entry.Platform = binary.LittleEndian.Uint16(h.Table[offset : offset+2])
		offset += 2

		entry.BlockIndex = binary.LittleEndian.Uint32(h.Table[offset : offset+4])
		offset += 4
	}

	return entries
}

func decryptDecompressTable(r io.Reader, dataSize, compressedSize uint64, key uint32) ([]byte, error) {
	crypted := make([]byte, compressedSize)
	if _, err := r.Read(crypted); err != nil {
		return nil, err
	}

	decryptBlock(crypted, int(compressedSize), key)

	if dataSize+extTableHeaderSize <= compressedSize {
		return crypted, nil
	}

	decompressed := make([]byte, dataSize)
	if err := decompress(decompressed, crypted); err != nil {
		return nil, err
	}

	return decompressed, nil
}
