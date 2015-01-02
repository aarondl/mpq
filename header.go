package mpq

import (
	"encoding/binary"
	"io"
)

const (
	mpqFormatVersion1 uint16 = iota // Up to Burning Crusade
	mpqFormatVersion2               // Burning Crusade and Newer
	mpqFormatVersion3               // WoW Cata Beta
	mpqFormatVersion4               // WoW Cata and newer

	headerArchive  = 0x1A
	headerUserData = 0x1B

	digestSize = 16
)

var (
	headerMPQ = []byte("MPQ")
)

// Header is the MPQ Archive's header.
type Header struct {
	HeaderSize    int
	Size          int
	FormatVersion uint16

	BlockSize      uint16
	HashTablePos   int
	BlockTablePos  int
	HashTableSize  int
	BlockTableSize int

	// MPQ Header v2
	HiBlockTablePos uint64
	HashTablePosHi  uint16
	BlockTablePosHi uint16

	// MPQ Header v3
	ArchiveSize uint64
	BETTablePos uint64
	HETTablePos uint64

	// MPQ Header v4
	HashTableSize64    uint64
	BlockTableSize64   uint64
	HiBlockTableSize64 uint64
	HETTableSize64     uint64
	BETTableSize64     uint64
	ChunkSize          int

	BlockTableMD5   []byte
	HashTableMD5    []byte
	HiBlockTableMD5 []byte
	BETTableMD5     []byte
	HETTableMD5     []byte
	MPQHeaderMD5    []byte
}

func (m *MPQ) readArchiveHeader(r io.Reader) error {
	var err error

	header := &Header{}

	buffer := make([]byte, 256)
	if _, err = r.Read(buffer[:28]); err != nil {
		return err
	}
	header.HeaderSize = int(binary.LittleEndian.Uint32(buffer[:4]))
	header.Size = int(binary.LittleEndian.Uint32(buffer[4:8]))
	header.FormatVersion = binary.LittleEndian.Uint16(buffer[8:10])

	header.BlockSize = binary.LittleEndian.Uint16(buffer[10:12])
	header.HashTablePos = int(binary.LittleEndian.Uint32(buffer[12:16]))
	header.BlockTablePos = int(binary.LittleEndian.Uint32(buffer[16:20]))
	header.HashTableSize = int(binary.LittleEndian.Uint32(buffer[20:24]))
	header.BlockTableSize = int(binary.LittleEndian.Uint32(buffer[24:28]))

	if header.FormatVersion >= 1 { // Version >= 2
		if _, err = r.Read(buffer[:12]); err != nil {
			return err
		}
		header.HiBlockTablePos = binary.LittleEndian.Uint64(buffer[:8])
		header.HashTablePosHi = binary.LittleEndian.Uint16(buffer[8:10])
		header.BlockTablePosHi = binary.LittleEndian.Uint16(buffer[10:12])
	}

	if header.FormatVersion >= 2 { // Version >= 3
		if _, err = r.Read(buffer[:24]); err != nil {
			return err
		}
		header.ArchiveSize = binary.LittleEndian.Uint64(buffer[:8])
		header.BETTablePos = binary.LittleEndian.Uint64(buffer[8:16])
		header.HETTablePos = binary.LittleEndian.Uint64(buffer[16:24])
	}

	if header.FormatVersion >= 3 { // Version >= 4
		if _, err = r.Read(buffer[:140]); err != nil {
			return err
		}

		header.HashTableSize64 = binary.LittleEndian.Uint64(buffer[:8])
		header.BlockTableSize64 = binary.LittleEndian.Uint64(buffer[8:16])
		header.HiBlockTableSize64 = binary.LittleEndian.Uint64(buffer[16:24])
		header.HETTableSize64 = binary.LittleEndian.Uint64(buffer[24:32])
		header.BETTableSize64 = binary.LittleEndian.Uint64(buffer[32:40])
		header.ChunkSize = int(binary.LittleEndian.Uint32(buffer[40:44]))

		header.BlockTableMD5 = make([]byte, digestSize)
		header.HashTableMD5 = make([]byte, digestSize)
		header.HiBlockTableMD5 = make([]byte, digestSize)
		header.BETTableMD5 = make([]byte, digestSize)
		header.HETTableMD5 = make([]byte, digestSize)
		header.MPQHeaderMD5 = make([]byte, digestSize)

		copy(header.BlockTableMD5, buffer[44:60])
		copy(header.HashTableMD5, buffer[60:76])
		copy(header.HiBlockTableMD5, buffer[76:92])
		copy(header.BETTableMD5, buffer[92:108])
		copy(header.HETTableMD5, buffer[108:124])
		copy(header.MPQHeaderMD5, buffer[124:140])
	}

	m.Header = header
	return nil
}
