package mpq

import "io"

type HashTable struct {
	Name1 int
	Name2 int

	Locale   uint16
	Platform uint16

	BlockIndex int
}

func (m *MPQ) readHashTable(r io.Reader) error {
	return nil
}
