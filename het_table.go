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

	Table        []byte
	FileIndicies []uint32
}

func (m *MPQ) readHETTable(r io.Reader) error {
	het := &HETTable{}

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

	offset := 32
	het.Table = make([]byte, het.TotalCount)
	copy(het.Table, buffer[offset:offset+het.TotalCount])
	offset += het.TotalCount

	// TODO: Read in file indicies
	// Array of file indexes. Bit size of each entry is taken from dwTotalIndexSize.
	// Table size is taken from dwHashTableSize.

	m.HETTable = het
	return nil
}

func decryptDecompressTable(r io.Reader, dataSize, compressedSize uint64, key uint32) ([]byte, error) {
	crypted := make([]byte, compressedSize)
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
