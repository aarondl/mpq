package mpq

import (
	"encoding/binary"
	"io"
)

const (
	cryptKeyHashTable  = 0xC3AF3770
	cryptKeyBlockTable = 0xEC83B3A3

	cryptTableSize = 0x500
)

var cryptTable []uint32

func init() {
	initCrypto()
}

func initCrypto() {
	if cryptTable != nil {
		return
	}
	cryptTable = make([]uint32, cryptTableSize)

	var seed, index1, index2, i uint32 = 0x00100001, 0, 0, 0

	for index1 = 0; index1 < 0x100; index1++ {
		for i, index2 = 0, index1; i < 5; i, index2 = i+1, index2+0x100 {
			var tmp1, tmp2 uint32

			seed = (seed*125 + 3) % 0x2AAAAB
			tmp1 = (seed & 0xFFFF) << 0x10

			seed = (seed*125 + 3) % 0x2AAAAB
			tmp2 = (seed & 0xFFFF)

			cryptTable[index2] = tmp1 | tmp2
		}
	}
}

type decryptReader struct {
	reader            io.Reader
	key1, key2, value uint32
}

func newDecryptReader(reader io.Reader, key1 uint32) *decryptReader {
	return &decryptReader{
		reader: reader,
		key1:   key1,
		key2:   0xEEEEEEEE,
	}
}

func (d *decryptReader) Read(buf []byte) (n int, err error) {
	n, err = d.reader.Read(buf)

	length := n >> 2
	for i := 0; i < length*4; i += 4 {
		d.key2 += cryptTable[0x400+(d.key1&0xFF)]

		value := binary.LittleEndian.Uint32(buf[i:])
		value ^= (d.key1 + d.key2)
		binary.LittleEndian.PutUint32(buf[i:], value)

		d.key1 = ((^d.key1 << 0x15) + 0x11111111) | (d.key1 >> 0x0B)
		d.key2 = value + d.key2 + (d.key2 << 5) + 3
	}

	return n, err
}

func decryptBlock(block []byte, length int, key1 uint32) {
	var value uint32
	var key2 uint32 = 0xEEEEEEEE

	length >>= 2

	for i := 0; i < length*4; i += 4 {
		key2 += cryptTable[0x400+(key1&0xFF)]

		value = binary.LittleEndian.Uint32(block[i:])
		value ^= (key1 + key2)
		binary.LittleEndian.PutUint32(block[i:], value)

		key1 = ((^key1 << 0x15) + 0x11111111) | (key1 >> 0x0B)
		key2 = value + key2 + (key2 << 5) + 3
	}
}
