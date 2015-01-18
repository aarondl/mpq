package mpq

import (
	"encoding/binary"
	"io"
)

const userDataHeaderSize = 16

// UserData is additional data inside the MPQ.
type UserData struct {
	UserDataHeader
}

type UserDataHeader struct {
	// offset off the user data.
	offset int64
	// MaxSize of user data
	MaxSize int
	// Offset of MPQ Header, relative to the begin of this header.
	HeaderOffset int
	// Size of the user data header?
	UserDataHeaderSize int
}

func (m *MPQ) readUserData(r io.Reader, offset int64) error {
	buffer := make([]byte, 12)
	if _, err := r.Read(buffer); err != nil {
		return err
	}

	userData := &UserData{}

	userData.MaxSize = int(binary.LittleEndian.Uint32(buffer[:4]))
	userData.HeaderOffset = int(binary.LittleEndian.Uint32(buffer[4:8]))
	userData.UserDataHeaderSize = int(binary.LittleEndian.Uint32(buffer[8:12]))

	userData.offset = offset + userDataHeaderSize // Offset 4 bytes to avoid header.

	m.UserData = userData
	return nil
}

// OpenUserData returns a reader that can be used to read the user data.
func (m *MPQ) OpenUserData() (io.Reader, error) {
	_, err := m.reader.Seek(m.UserData.offset, 0)
	if err != nil {
		return nil, err
	}

	return io.LimitReader(m.reader, int64(m.UserData.MaxSize)), nil
}
