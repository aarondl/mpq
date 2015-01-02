package mpq

import "encoding/binary"

const (
	cryptKeyHashTable  = 0xC3AF3770
	cryptKeyBlockTable = 0xEC83B3A3

	cryptTableSize = 0x500
)

var cryptTable []uint32

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

func decryptBlock(block []byte, length int, key1 uint32) {
	initCrypto()

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

/*
void prepareCryptTable()
{
    unsigned long seed = 0x00100001, index1 = 0, index2 = 0, i;

    for(index1 = 0; index1 < 0x100; index1++)
    {
        for(index2 = index1, i = 0; i < 5; i++, index2 += 0x100)
        {
            unsigned long temp1, temp2;

            seed = (seed * 125 + 3) % 0x2AAAAB;
            temp1 = (seed & 0xFFFF) << 0x10;

            seed = (seed * 125 + 3) % 0x2AAAAB;
            temp2 = (seed & 0xFFFF);

            cryptTable[index2] = (temp1 | temp2);
        }
    }
}

void DecryptBlock(void *block, long length, unsigned long key)
{
    unsigned long seed = 0xEEEEEEEE, unsigned long ch;
    unsigned long *castBlock = (unsigned long *)block;

    // Round to longs
    length >>= 2;

    while(length-- > 0)
    {
        seed += stormBuffer[0x400 + (key & 0xFF)];
        ch = *castBlock ^ (key + seed);

        key = ((~key << 0x15) + 0x11111111) | (key >> 0x0B);
        seed = ch + seed + (seed << 5) + 3;
        *castBlock++ = ch;
    }
}
*/
