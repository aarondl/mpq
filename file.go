package mpq

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
)

const (
	fileFlagImplode    = 0x00000100 // File is compressed using PKWARE Data compression library
	fileFlagCompress   = 0x00000200 // File is compressed using combination of compression methods
	fileFlagEncrypted  = 0x00010000 // The file is encrypted
	fileFlagFixKey     = 0x00020000 // The decryption key for the file is altered according to the position of the file in the archive
	fileFlagPatchFile  = 0x00100000 // The file contains incremental patch for an existing file in base MPQ
	fileFlagSingleUnit = 0x01000000 // Instead of being divided to 0x1000-bytes blocks, the file is stored as single unit
	fileFlagDelete     = 0x02000000 // File is a deletion marker, indicating that the file no longer exists. This is used to allow patch archives to delete files present in lower-priority archives in the search chain. The file usually has length of 0 or 1 byte and its name is a hash
	fileFlagSectorCRC  = 0x04000000 // File has checksums for each sector (explained in the File Data section). Ignored if file is not compressed or imploded.
	fileFlagExists     = 0x80000000 // Set if file exists, reset when the file was deleted

	fileCompressedMask = 0x0000FF00 // Mask for a file being compressed
)

// These errors are possible return values (along with others) from the mpq.Open() call.
var (
	// ErrFileNotFound occurs on mpq.Open() when the filename given is
	// not in the archive.
	ErrFileNotFound = errors.New("File not found in archive")
	// ErrFileDeleted occurs when the filename given is present in the BET/Block tables
	// but has been flagged as deleted.
	ErrFileDeleted = errors.New("File has been removed from the archive")
	// ErrFileEmpty occurs when the file is of size 0 bytes inside the archive.
	ErrFileEmpty = errors.New("File is empty")
)

// File represents a file in the MPQ archive.
type File struct {
	Name   string
	Locale uint16

	FileSize       uint64
	CompressedSize uint64
	Position       uint64

	Flags uint32
}

// Open the file for reading.
func (m *MPQ) Open(filename string) (io.Reader, error) {
	var file *File
	var ok bool
	if file, ok = m.FileList[filename]; !ok {
		return nil, ErrFileNotFound
	}

	return m.open(file)
}

func (m *MPQ) open(file *File) (io.Reader, error) {
	if file.Position == 0 || file.FileSize == 0 || file.CompressedSize == 0 {
		return nil, ErrFileEmpty
	}

	var err error
	if _, err = m.reader.Seek(m.offset+int64(file.Position), 0); err != nil {
		return nil, fmt.Errorf("Failed to seek to file position: %v", err)
	}

	reader := io.LimitReader(m.reader, int64(file.CompressedSize))

	if file.Flags&fileFlagExists == 0 {
		return nil, ErrFileDeleted
	}

	if file.Flags&fileFlagSingleUnit == 0 {
		return nil, errors.New("Cannot process multi-unit files")
	}

	if file.Flags&fileFlagEncrypted != 0 {
		key := blizz(file.Name, blizzHashFileKey)
		if file.Flags&fileFlagFixKey != 0 {
			key = (key + uint32(m.offset)) ^ uint32(file.FileSize)
		}
		reader = newDecryptReader(reader, key)
	}

	if file.Flags&fileCompressedMask != 0 && file.FileSize != file.CompressedSize {
		if file.Flags&fileFlagCompress != 0 {
			if m.Header.FormatVersion >= mpqFormatVersion2 {
				reader, err = newDecompressReader(reader, file.FileSize)
			} else {
				err = errors.New("Oldschool MPQ multiple compression is not supported")
			}
		} else if file.Flags&fileFlagImplode != 0 {
			err = errors.New("PKWARE Implode Compression not supported")
		}
	}

	return reader, err
}

// buildFileList attempts to use the read in structures to create a file listing.
func (m *MPQ) buildFileList() error {
	var err error
	var file *File

	listInfo, err := m.FileInfo("(listfile)")
	if err != nil {
		return err
	}

	list, err := m.open(listInfo)
	if err != nil {
		return err
	}

	m.fileNames = make([]string, 0)
	scanner := bufio.NewScanner(list)
	for scanner.Scan() {
		m.fileNames = append(m.fileNames, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Make sure to fetch special file info.
	m.fileNames = append(m.fileNames, []string{"(attributes)", "(userdata)"}...)

	for i := 0; i < len(m.fileNames); i++ {
		fileName := m.fileNames[i]
		if file, err = m.FileInfo(fileName); err != nil {
			if i+1 != len(m.fileNames) {
				m.fileNames[i], m.fileNames[i+1] = m.fileNames[i+1], m.fileNames[i]
			}
			m.fileNames = m.fileNames[:len(m.fileNames)-1]
		} else if err == nil {
			m.FileList[fileName] = file
		}
	}

	// Add the pre-done list info stuff.
	m.fileNames = append(m.fileNames, listInfo.Name)
	m.FileList[listInfo.Name] = listInfo

	sort.Strings(m.fileNames)

	return nil
}

// FileInfo attempts to get the file information for a filename.
func (m *MPQ) FileInfo(name string) (*File, error) {
	if m.HETTable != nil && m.BETTable != nil {
		return m.findFromHETAndBET(name)
	} else if m.HashTable != nil && m.BlockTable != nil {
		return m.findFromHashAndBlock(name)
	}

	return nil, errors.New("HET, BET, Hash and Block tables are all unavailable")
}

func (m *MPQ) findFromHETAndBET(name string) (*File, error) {
	hash := (jenkins2(name) & m.HETTable.AndMask) | m.HETTable.OrMask
	hetHash := byte(hash >> uint(m.HETTable.HashEntrySize-8))
	betHash := hash & (m.HETTable.AndMask >> 0x08)

	indexes, err := m.HETTable.Indexes()
	if err != nil {
		return nil, err
	}

	files, err := m.BETTable.Entries()
	if err != nil {
		return nil, err
	}

	var betEntry *BETTableEntry
	for i := int(hash % uint64(m.HETTable.HashTableSize)); i < len(m.HETTable.Hashes); i++ {
		nameHash1 := m.HETTable.Hashes[i]
		if nameHash1 == 0 {
			break
		}

		if hetHash != nameHash1 {
			continue
		}

		betEntry = &files[int(indexes[i])]
		if betHash == betEntry.NameHash2 {
			break
		}
		betEntry = nil
	}

	if betEntry == nil {
		return nil, ErrFileNotFound
	}

	return &File{
		Name:           name,
		FileSize:       betEntry.FileSize,
		CompressedSize: betEntry.CompressedSize,
		Position:       betEntry.FilePosition,
		Flags:          betEntry.Flags,
	}, nil
}

func (m *MPQ) findFromHashAndBlock(name string) (*File, error) {
	start := blizz(name, blizzHashTableIndex) & uint32(m.Header.HashTableSize-1)
	name1 := blizz(name, blizzHashNameA)
	name2 := blizz(name, blizzHashNameB)

	hashTableEntries := m.HashTable.Entries()
	blockTableEntries := m.BlockTable.Entries()

	var blockEntry *BlockTableEntry
	for i := int(start); i < len(hashTableEntries); i++ {
		entry := &hashTableEntries[i]
		if entry.BlockIndex == 0xFFFFFFFF {
			break
		}

		if name1 == entry.Name1 && name2 == entry.Name2 {
			blockEntry = &blockTableEntries[entry.BlockIndex]
		}
	}

	if blockEntry == nil {
		return nil, ErrFileNotFound
	}

	return &File{
		Name:           name,
		FileSize:       uint64(blockEntry.FileSize),
		CompressedSize: uint64(blockEntry.CompressedSize),
		Position:       uint64(blockEntry.FilePosition),
		Flags:          uint32(blockEntry.Flags),
	}, nil
}
