/*
Package mpq provides read-only access to an MPQ file.

Although this package may seem
complete and wonderful I assure you it is not. There are a great many MPQ files with protections
in them (typically mangled by third party tools not related to the creators of MPQ files) that
cannot be handled by this package.

Furthermore there are a number of strange special cases and legacy situations that can arise
for MPQ files, and as such these special cases and odd MPQ files may also fail to load.
*/
package mpq

import (
	"bytes"
	"errors"
	"io"
	"os"
)

// MPQ represents a single MPQ file and allows access to all fields
// and contained files.
type MPQ struct {
	reader   io.ReadSeeker
	Header   *Header
	UserData *UserData

	BETTable     *BETTable
	HETTable     *HETTable
	HashTable    *HashTable
	BlockTable   *BlockTable
	HiBlockTable *HiBlockTable

	offset int64

	fileNames []string
	FileList  map[string]*File
}

// Open an MPQ File for reading.
func Open(filename string) (*MPQ, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return OpenReader(f)
}

// OpenReader opens a stream that contains an MPQ file for reading.
func OpenReader(reader io.ReadSeeker) (*MPQ, error) {
	var buffer [4]byte

	m := &MPQ{reader: reader, FileList: make(map[string]*File)}

	var err error
	readHeader := false
	for !readHeader {
		if _, err = reader.Read(buffer[:]); err != nil {
			return nil, err
		}

		if bytes.Compare(buffer[:3], headerMPQ) == 0 {
			if buffer[3] == headerArchive {
				if err = m.readArchiveHeader(reader); err != nil {
					return nil, err
				}
				readHeader = true
				break
			} else if buffer[3] == headerUserData {
				if err = m.readUserData(reader, m.offset); err != nil {
					return nil, err
				}
			}
		}

		m.offset += 512
		reader.Seek(m.offset, 0)
	}

	if !readHeader {
		return nil, errors.New("Could not find MPQ header.")
	}

	if m.Header.HETTablePos != 0 {
		hetOffset := m.offset + int64(m.Header.HETTablePos)
		if _, err = reader.Seek(hetOffset, 0); err != nil {
			return nil, err
		}
		if err = m.readHETTable(reader); err != nil {
			return nil, err
		}
	}

	if m.Header.BETTablePos != 0 {
		betOffset := m.offset + int64(m.Header.BETTablePos)
		if _, err = reader.Seek(betOffset, 0); err != nil {
			return nil, err
		}
		if err = m.readBETTable(reader); err != nil {
			return nil, err
		}
	}

	if m.Header.HashTablePos != 0 || m.Header.HashTablePosHi != 0 {
		pos := m.offset + ((int64(m.Header.HashTablePosHi) << 32) | int64(m.Header.HashTablePos))
		if _, err = reader.Seek(pos, 0); err != nil {
			return nil, err
		}
		if err = m.readHashTable(reader); err != nil {
			return nil, err
		}
	}

	if m.Header.BlockTablePos != 0 || m.Header.BlockTablePosHi != 0 {
		pos := m.offset + ((int64(m.Header.BlockTablePosHi) << 32) | int64(m.Header.BlockTablePos))
		if _, err = reader.Seek(pos, 0); err != nil {
			return nil, err
		}
		if err = m.readBlockTable(reader); err != nil {
			return nil, err
		}
	}

	if m.Header.HiBlockTablePos != 0 {
		pos := m.offset + int64(m.Header.HiBlockTablePos)
		if _, err = reader.Seek(pos, 0); err != nil {
			return nil, err
		}
		if err = m.readHiBlockTable(reader); err != nil {
			return nil, err
		}
	}

	if err = m.buildFileList(); err != nil {
		return nil, err
	}

	return m, nil
}

// Files in the archive.
func (m *MPQ) Files() ([]string, error) {
	if m.fileNames == nil {
		if err := m.buildFileList(); err != nil {
			return nil, err
		}
	}

	files := make([]string, len(m.fileNames))
	copy(files, m.fileNames)

	return files, nil
}

// Close attempts to close the MPQ file handle if the given stream has a close.
func (m *MPQ) Close() error {
	if closer, ok := m.reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
