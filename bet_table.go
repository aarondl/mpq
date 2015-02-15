package mpq

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/aarondl/bitstream"
)

var (
	headerBETTable = []byte("BET\x1A")

	errorBETTableBounds = errors.New("BET Table ended unexpectedly")
)

// BETTable from the MPQ Header.
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

	TableEntries []byte
	Hashes       []byte

	entries []BETTableEntry
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

	buffer, err := decryptDecompressExtTable(r, uint64(bet.DataSize), m.Header.BETTableSize64, cryptKeyBlockTable)
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

	offset := 76
	bet.Flags = make([]uint32, bet.FlagCount)
	for i := 0; i < bet.FlagCount; i++ {
		bet.Flags[i] = binary.LittleEndian.Uint32(buffer[offset : offset+4])
		offset += 4
	}

	bet.TableEntries = make([]byte, (bet.TableEntrySize*bet.EntryCount+7)/8)
	copy(bet.TableEntries, buffer[offset:offset+len(bet.TableEntries)])
	offset += len(bet.TableEntries)

	bet.Hashes = make([]byte, (bet.HashSizeTotal*bet.EntryCount+7)/8)
	copy(bet.Hashes, buffer[offset:offset+len(bet.Hashes)])

	m.BETTable = bet
	return nil
}

// BETTableEntry is a table entry.
type BETTableEntry struct {
	NameHash2      uint64
	FilePosition   uint64
	FileSize       uint64
	CompressedSize uint64
	FlagIndex      uint32
	Flags          uint32
}

// Entries parses the TableEntries and Hashes bit arrays into an array of BETTableEntry.
func (b *BETTable) Entries() ([]BETTableEntry, error) {
	if b.entries != nil {
		return b.entries, nil
	}
	entries := make([]BETTableEntry, b.EntryCount)
	barr := bitstream.New(bytes.NewBuffer(b.TableEntries))

	var val uint64
	var err error
	for i := 0; i < b.EntryCount; i++ {
		entry := &entries[i]

		if val, err = barr.Bits(b.BitCountFilePos); err != nil {
			return nil, errorBETTableBounds
		}
		entry.FilePosition = val

		if val, err = barr.Bits(b.BitCountFileSize); err != nil {
			return nil, errorBETTableBounds
		}
		entry.FileSize = val

		if val, err = barr.Bits(b.BitCountCmpSize); err != nil {
			return nil, errorBETTableBounds
		}
		entry.CompressedSize = val

		if val, err = barr.Bits(b.BitCountFlagIndex); err != nil {
			return nil, errorBETTableBounds
		}
		entry.FlagIndex = uint32(val)

		entry.Flags = b.Flags[entry.FlagIndex]
	}

	barr = bitstream.New(bytes.NewBuffer(b.Hashes))
	for i := 0; i < b.EntryCount; i++ {
		if val, err = barr.Bits(b.HashSizeTotal); err != nil {
			return nil, errorBETTableBounds
		}

		entries[i].NameHash2 = val
	}

	b.entries = entries
	return entries, nil
}
