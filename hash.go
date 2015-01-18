package mpq

import (
	"encoding/binary"
	"strings"
)

const (
	blizzHashTableIndex = 0x000
	blizzHashNameA      = 0x100
	blizzHashNameB      = 0x200
	blizzHashFileKey    = 0x300
	blizzHashKey2Mix    = 0x400
)

func blizz(filename string, hashType uint32) uint32 {
	initCrypto()

	var seed1, seed2 uint32 = 0x7FED7FED, 0xEEEEEEEE

	filename = strings.ToUpper(strings.Replace(filename, `/`, `\`, -1))

	for i := 0; i < len(filename); i++ {
		seed1 = cryptTable[(hashType)+uint32(filename[i])] ^ (seed1 + seed2)
		seed2 = uint32(filename[i]) + seed1 + seed2 + (seed2 << 5) + 3
	}
	return seed1
}

func jenkins2(filename string) uint64 {
	var primary_hash uint32 = 1
	var secondary_hash uint32 = 2

	filename = strings.ToLower(strings.Replace(filename, `/`, `\`, -1))

	hashLittle2Str(filename, &secondary_hash, &primary_hash)

	return uint64(primary_hash)*uint64(0x100000000) + uint64(secondary_hash)
}

func hashLittle2Str(str string, pc, pb *uint32) {
	hashLittle2([]byte(str), pc, pb)
}

func hashLittle2(key []byte, pc, pb *uint32) {
	var a, b, c uint32
	a = 0xdeadbeef + uint32(len(key)) + *pc
	b, c = a, a
	length := len(key)

	c += *pb

	offset := 0
	// As I ported this I ignored this if statement which seems to be dealing with some sort of alignment
	// and ability to read larger chunks of data at the time. In Go everything is aligned and we never have
	// to worry about this since we don't use pointer arithmetic. But I had already ported the other parts
	// of the if statement even though they'll never be used. Educational purposes folks.
	if true || (key[0]&0x3) == 0 {
		for ; length > 12; length -= 12 {
			a += binary.LittleEndian.Uint32(key[offset : offset+4])
			offset += 4
			b += binary.LittleEndian.Uint32(key[offset : offset+4])
			offset += 4
			c += binary.LittleEndian.Uint32(key[offset : offset+4])
			offset += 4

			hashLittleMix(&a, &b, &c)
		}

		switch length {
		case 12:
			c += binary.LittleEndian.Uint32(key[offset+8 : offset+12])
			b += binary.LittleEndian.Uint32(key[offset+4 : offset+8])
			a += binary.LittleEndian.Uint32(key[offset : offset+4])
		case 11:
			c += uint32(key[offset+10]) << 16
			fallthrough
		case 10:
			c += uint32(key[offset+9]) << 8
			fallthrough
		case 9:
			c += uint32(key[offset+8])
			fallthrough
		case 8:
			b += binary.LittleEndian.Uint32(key[offset+4 : offset+8])
			a += binary.LittleEndian.Uint32(key[offset : offset+4])
		case 7:
			b += uint32(key[offset+6]) << 16
			fallthrough
		case 6:
			b += uint32(key[offset+5]) << 8
			fallthrough
		case 5:
			b += uint32(key[offset+4])
			fallthrough
		case 4:
			a += binary.LittleEndian.Uint32(key[offset : offset+4])
		case 3:
			a += uint32(key[offset+2]) << 16
			fallthrough
		case 2:
			a += uint32(key[offset+1]) << 8
			fallthrough
		case 1:
			a += uint32(key[offset])
		case 0:
			*pc = c
			*pb = b
			return
		}

	} else if (key[0] & 0x1) == 0 {
		for ; length > 12; length -= 12 {
			a += uint32(binary.LittleEndian.Uint16(key[offset:offset+2])) << 16
			offset += 2
			b += uint32(binary.LittleEndian.Uint16(key[offset:offset+2])) << 16
			offset += 2
			c += uint32(binary.LittleEndian.Uint16(key[offset:offset+2])) << 16
			offset += 2

			hashLittleMix(&a, &b, &c)
		}

		switch length {
		case 12:
			c += uint32(binary.LittleEndian.Uint16(key[offset+8:offset+10])) + (uint32((binary.LittleEndian.Uint16(key[offset+10 : offset+12]))) << 16)
			b += uint32(binary.LittleEndian.Uint16(key[offset+4:offset+6])) + (uint32((binary.LittleEndian.Uint16(key[offset+6 : offset+8]))) << 16)
			a += uint32(binary.LittleEndian.Uint16(key[offset:offset+2])) + (uint32((binary.LittleEndian.Uint16(key[offset+2 : offset+4]))) << 16)
		case 11:
			c += uint32(key[offset+10]) << 16
			fallthrough
		case 10:
			c += uint32(binary.LittleEndian.Uint16(key[offset+8 : offset+10]))
			b += uint32(binary.LittleEndian.Uint16(key[offset+4:offset+6])) + (uint32((binary.LittleEndian.Uint16(key[offset+6 : offset+8]))) << 16)
			a += uint32(binary.LittleEndian.Uint16(key[offset:offset+2])) + (uint32((binary.LittleEndian.Uint16(key[offset+2 : offset+4]))) << 16)
		case 9:
			c += uint32(key[offset+8])
			fallthrough
		case 8:
			b += uint32(binary.LittleEndian.Uint16(key[offset+4:offset+6])) + (uint32((binary.LittleEndian.Uint16(key[offset+6 : offset+8]))) << 16)
			a += uint32(binary.LittleEndian.Uint16(key[offset:offset+2])) + (uint32((binary.LittleEndian.Uint16(key[offset+2 : offset+4]))) << 16)
		case 7:
			b += uint32(key[offset+6]) << 16
			fallthrough
		case 6:
			b += uint32(binary.LittleEndian.Uint16(key[offset+4 : offset+6]))
			a += uint32(binary.LittleEndian.Uint16(key[offset:offset+2])) + (uint32((binary.LittleEndian.Uint16(key[offset+2 : offset+4]))) << 16)
		case 5:
			b += uint32(key[offset+4])
			fallthrough
		case 4:
			a += uint32(binary.LittleEndian.Uint16(key[offset:offset+2])) + (uint32((binary.LittleEndian.Uint16(key[offset+2 : offset+4]))) << 16)
		case 3:
			a += uint32(key[offset+2]) << 16
			fallthrough
		case 2:
			a += uint32(binary.LittleEndian.Uint16(key[offset : offset+2]))
		case 1:
			a += uint32(key[0])
		case 0:
			*pc = c
			*pb = b
			return
		}
	} else {
		for ; length > 12; length -= 12 {
			a += uint32(key[offset])
			a += uint32(key[offset+1]) << 8
			a += uint32(key[offset+2]) << 16
			a += uint32(key[offset+3]) << 24
			offset += 4
			b += uint32(key[offset])
			b += uint32(key[offset+1]) << 8
			b += uint32(key[offset+2]) << 16
			b += uint32(key[offset+3]) << 24
			offset += 4
			c += uint32(key[offset])
			c += uint32(key[offset+1]) << 8
			c += uint32(key[offset+2]) << 16
			c += uint32(key[offset+3]) << 24
			offset += 4

			hashLittleMix(&a, &b, &c)
		}

		switch length {
		case 12:
			c += uint32(key[offset+11]) << 24
		case 11:
			c += uint32(key[offset+10]) << 16
		case 10:
			c += uint32(key[offset+9]) << 8
		case 9:
			c += uint32(key[8])
		case 8:
			b += uint32(key[offset+7]) << 24
		case 7:
			b += uint32(key[offset+6]) << 16
		case 6:
			b += uint32(key[offset+5]) << 8
		case 5:
			b += uint32(key[4])
		case 4:
			a += uint32(key[offset+3]) << 24
		case 3:
			a += uint32(key[offset+2]) << 16
		case 2:
			a += uint32(key[offset+1]) << 8
		case 1:
			a += uint32(key[0])
		}
	}

	hashLittleFinal(&a, &b, &c)
	*pc = c
	*pb = b
}

func hashLittleMix(a, b, c *uint32) {
	*a -= *c
	*a ^= hashLittleRot(*c, 4)
	*c += *b

	*b -= *a
	*b ^= hashLittleRot(*a, 6)
	*a += *c

	*c -= *b
	*c ^= hashLittleRot(*b, 8)
	*b += *a

	*a -= *c
	*a ^= hashLittleRot(*c, 16)
	*c += *b

	*b -= *a
	*b ^= hashLittleRot(*a, 19)
	*a += *c

	*c -= *b
	*c ^= hashLittleRot(*b, 4)
	*b += *a
}

func hashLittleFinal(a, b, c *uint32) {
	*c ^= *b
	*c -= hashLittleRot(*b, 14)

	*a ^= *c
	*a -= hashLittleRot(*c, 11)

	*b ^= *a
	*b -= hashLittleRot(*a, 25)

	*c ^= *b
	*c -= hashLittleRot(*b, 16)

	*a ^= *c
	*a -= hashLittleRot(*c, 4)

	*b ^= *a
	*b -= hashLittleRot(*a, 14)

	*c ^= *b
	*c -= hashLittleRot(*b, 24)
}

func hashLittleRot(x, k uint32) uint32 {
	return (x << k) | (x >> (32 - k))
}
