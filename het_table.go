package mpq

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/aarondl/bitstream"
)

const (
	extTableHeaderSize = 12
)

var (
	headerHETTable = []byte("HET\x1A")
)

// HETTable from the MPQ Header.
type HETTable struct {
	Version  int
	DataSize int

	TableSize      int
	EntryCount     int
	HashTableSize  int
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

	buffer, err := decryptDecompressExtTable(r, uint64(het.DataSize), m.Header.HETTableSize64, cryptKeyHashTable)
	if err != nil {
		return err
	}

	het.TableSize = int(binary.LittleEndian.Uint32(buffer[0:4]))
	het.EntryCount = int(binary.LittleEndian.Uint32(buffer[4:8]))
	het.HashTableSize = int(binary.LittleEndian.Uint32(buffer[8:12]))
	het.HashEntrySize = int(binary.LittleEndian.Uint32(buffer[12:16]))
	het.IndexSizeTotal = int(binary.LittleEndian.Uint32(buffer[16:20]))
	het.IndexSizeExtra = int(binary.LittleEndian.Uint32(buffer[20:24]))
	het.IndexSize = int(binary.LittleEndian.Uint32(buffer[24:28]))
	het.BlockTableSize = int(binary.LittleEndian.Uint32(buffer[28:32]))

	// Read Table Information
	if het.HashEntrySize != 0x40 {
		het.AndMask = uint64(1) << uint(het.HashEntrySize)
	}
	het.AndMask--
	het.OrMask = uint64(1) << uint(het.HashEntrySize-1)

	het.count = het.HashTableSize
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
func (h *HETTable) Indexes() ([]uint, error) {
	ret := make([]uint, h.count)
	b := bitstream.New(bytes.NewBuffer(h.Indicies))

	var val uint64
	var err error
	for i := 0; i < h.count; i++ {
		if val, err = b.Bits(h.bitCount); err != nil {
			return nil, errors.New("HET Table ended unexpectedly")
		}

		ret[i] = uint(val)
	}

	return ret, nil
}

func decryptDecompressExtTable(r io.Reader, dataSize, compressedSize uint64, key uint32) ([]byte, error) {
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
