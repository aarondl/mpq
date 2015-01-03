package mpq

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	extTableHeaderSize = 12
)

var (
	headerHETTable = []byte("HET\x1A")
)

type HETTable struct {
	Version  int
	DataSize int

	TableSize      int
	EntryCount     int
	TotalCount     int
	HashEntrySize  int
	IndexSizeTotal int
	IndexSizeExtra int
	IndexSize      int
	BlockTableSize int

	count    int
	bitCount int
	Hashes   []byte
	Indicies []byte

	AndMask uint64
	OrMask  uint64
}

func (m *MPQ) readHETTable(r io.Reader) error {
	het := &HETTable{}

	// Read Header Information
	header := make([]byte, extTableHeaderSize)
	if _, err := r.Read(header); err != nil {
		return err
	}

	if bytes.Compare(headerHETTable, header[:4]) != 0 {
		return fmt.Errorf("HET Table header not found, got: %02X", header[:4])
	}
	het.Version = int(binary.LittleEndian.Uint32(header[4:8]))
	het.DataSize = int(binary.LittleEndian.Uint32(header[8:12]))

	buffer, err := decryptDecompressTable(r, uint64(het.DataSize), m.Header.HETTableSize64, cryptKeyHashTable)
	if err != nil {
		return err
	}

	het.TableSize = int(binary.LittleEndian.Uint32(buffer[0:4]))
	het.EntryCount = int(binary.LittleEndian.Uint32(buffer[4:8]))
	het.TotalCount = int(binary.LittleEndian.Uint32(buffer[8:12]))
	het.HashEntrySize = int(binary.LittleEndian.Uint32(buffer[12:16]))
	het.IndexSizeTotal = int(binary.LittleEndian.Uint32(buffer[16:20]))
	het.IndexSizeExtra = int(binary.LittleEndian.Uint32(buffer[20:24]))
	het.IndexSize = int(binary.LittleEndian.Uint32(buffer[24:28]))
	het.BlockTableSize = int(binary.LittleEndian.Uint32(buffer[28:32]))

	// Read Table Information
	if het.HashEntrySize != 0x40 {
		het.AndMask = 1 << uint(het.HashEntrySize)
	}
	het.AndMask -= 1
	het.OrMask = 1 << uint(het.HashEntrySize-1)

	het.count = het.TotalCount
	if het.count == 0 {
		het.count = (het.EntryCount * 4) / 3
	}

	maxValue := het.EntryCount
	for maxValue > 0 {
		maxValue >>= 1
		het.bitCount++
	}

	het.Hashes = make([]byte, het.count)
	het.Indicies = make([]byte, (het.count*het.bitCount+7)/8)
	for i := 0; i < len(het.Indicies); i++ {
		het.Indicies[i] = 0xFF
	}

	offset := 32
	copy(het.Hashes, buffer[offset:offset+het.count])
	offset += het.count
	copy(het.Indicies, buffer[offset:offset+len(het.Indicies)])

	m.HETTable = het
	return nil
}

// Indexes reads the bit array from the het.Indicies and turns it into a uint array.
func (h *HETTable) Indexes() []uint {
	ret := make([]uint, h.count)

	bitOffset := uint(8)
	var curByte byte
	byteIndex := 0
	for i := 0; i < h.count; i++ {
		for j := uint(0); j < uint(h.bitCount); j++ {
			if bitOffset == 8 {
				bitOffset = 0
				curByte = h.Indicies[byteIndex]
				byteIndex++
			}

			ret[i] |= (uint(curByte&(1<<bitOffset)) >> bitOffset) << uint(j)
			bitOffset++
		}
	}

	return ret
}

func decryptDecompressTable(r io.Reader, dataSize, compressedSize uint64, key uint32) ([]byte, error) {
	crypted := make([]byte, compressedSize-extTableHeaderSize)
	if _, err := r.Read(crypted); err != nil {
		return nil, err
	}

	decryptBlock(crypted, int(compressedSize-extTableHeaderSize), key)

	if dataSize+extTableHeaderSize <= compressedSize {
		return crypted, nil
	}

	decompressed := make([]byte, dataSize)
	if err := decompress(decompressed, crypted); err != nil {
		return nil, err
	}

	return decompressed, nil
}
