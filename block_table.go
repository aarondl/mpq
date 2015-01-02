package mpq

import "io"

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
)

type BlockTable struct {
	FilePos        int
	CompressedSize int
	FileSize       int
	Flags          int
}

func (m *MPQ) readBlockTable(r io.Reader) error {
	return nil
}
