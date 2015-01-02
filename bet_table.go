package mpq

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

var (
	headerBETTable = []byte("BET\x1A")
)

type BETTable struct {
	Version  int
	DataSize int

	TableSize      int
	EntryCount     int
	Unknown08      int // 0x10
	TableEntrySize int

	BitIndexFilePos   int
	BitIndexFileSize  int
	BitIndexCmpSize   int
	BitIndexFlagIndex int
	BitIndexUnknown   int

	BitCountFilePos   int
	BitCountFileSize  int
	BitCountCmpSize   int
	BitCountFlagIndex int
	BitCountUnknown   int

	HashSizeTotal int
	HashSizeExtra int
	HashSize      int
	HashArraySize int

	FlagCount int
	Flags     []uint32

	// File table. Size of each entry is taken from dwTableEntrySize.
	// Size of the table is (dwTableEntrySize * dwMaxFileCount), round up to 8.

	// Array of BET hashes. Table size is taken from dwMaxFileCount from HET table
}

func (m *MPQ) readBETTable(r io.Reader) error {
	bet := &BETTable{}

	header := make([]byte, extTableHeaderSize)
	if _, err := r.Read(header); err != nil {
		return err
	}

	if bytes.Compare(headerBETTable, header[:4]) != 0 {
		return fmt.Errorf("BET Table header not found, got: %02X", header[:4])
	}
	bet.Version = int(binary.LittleEndian.Uint32(header[4:8]))
	bet.DataSize = int(binary.LittleEndian.Uint32(header[8:12]))

	buffer, err := decryptDecompressTable(r, uint64(bet.DataSize), m.Header.BETTableSize64, cryptKeyBlockTable)
	if err != nil {
		return err
	}

	bet.TableSize = int(binary.LittleEndian.Uint32(buffer[0:4]))
	bet.EntryCount = int(binary.LittleEndian.Uint32(buffer[4:8]))
	bet.Unknown08 = int(binary.LittleEndian.Uint32(buffer[8:12]))
	bet.TableEntrySize = int(binary.LittleEndian.Uint32(buffer[12:16]))
	bet.BitIndexFilePos = int(binary.LittleEndian.Uint32(buffer[16:20]))
	bet.BitIndexFileSize = int(binary.LittleEndian.Uint32(buffer[20:24]))
	bet.BitIndexCmpSize = int(binary.LittleEndian.Uint32(buffer[24:28]))
	bet.BitIndexFlagIndex = int(binary.LittleEndian.Uint32(buffer[28:32]))
	bet.BitIndexUnknown = int(binary.LittleEndian.Uint32(buffer[32:36]))
	bet.BitCountFilePos = int(binary.LittleEndian.Uint32(buffer[36:40]))
	bet.BitCountFileSize = int(binary.LittleEndian.Uint32(buffer[40:44]))
	bet.BitCountCmpSize = int(binary.LittleEndian.Uint32(buffer[44:48]))
	bet.BitCountFlagIndex = int(binary.LittleEndian.Uint32(buffer[48:52]))
	bet.BitCountUnknown = int(binary.LittleEndian.Uint32(buffer[52:56]))
	bet.HashSizeTotal = int(binary.LittleEndian.Uint32(buffer[56:60]))
	bet.HashSizeExtra = int(binary.LittleEndian.Uint32(buffer[60:64]))
	bet.HashSize = int(binary.LittleEndian.Uint32(buffer[64:68]))
	bet.HashArraySize = int(binary.LittleEndian.Uint32(buffer[68:72]))
	bet.FlagCount = int(binary.LittleEndian.Uint32(buffer[72:76]))

	m.BETTable = bet
	return nil
}
