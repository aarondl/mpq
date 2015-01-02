package mpq

import (
	"encoding/binary"
	"io"
)

// UserData is additional data inside the MPQ.
type UserData struct {
	UserDataHeader
	Data []byte
}

type UserDataHeader struct {
	// Maximum size of user data
	MaxSize int
	// Offset of MPQ Header, relative to the begin of this header.
	HeaderOffset int
	// Size of the user data header?
	UserDataHeaderSize int
}

func (m *MPQ) readUserData(r io.Reader) error {
	buffer := make([]byte, 12)
	if _, err := r.Read(buffer); err != nil {
		return err
	}

	userData := &UserData{}

	userData.MaxSize = int(binary.LittleEndian.Uint32(buffer[:4]))
	userData.HeaderOffset = int(binary.LittleEndian.Uint32(buffer[4:8]))
	userData.UserDataHeaderSize = int(binary.LittleEndian.Uint32(buffer[8:12]))

	userData.Data = make([]byte, userData.MaxSize)
	if _, err := r.Read(userData.Data); err != nil {
		return err
	}

	m.UserData = userData
	return nil
}

func (u *UserData) decodeUserData(r io.Reader) error {
	return nil
}
