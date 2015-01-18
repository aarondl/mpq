package mpq

import (
	"errors"
	"fmt"
)

type bitArray struct {
	bits       []byte
	byteOffset int
	bitOffset  uint
}

func newBitArray(b []byte) *bitArray {
	return &bitArray{bits: b}
}

func (b *bitArray) next(nBits int) (int64, error) {
	var val int64
	var bitOffset uint
	for nBits > 0 {
		if b.bitOffset == 8 {
			b.bitOffset = 0
			b.byteOffset++

			if b.byteOffset >= len(b.bits) && nBits > 0 {
				return 0, errors.New("Ran out of bits.")
			}
		}

		curByte := b.bits[b.byteOffset]

		maskSize := uint(nBits)
		if maskSize > (8 - b.bitOffset) {
			maskSize = 8 - b.bitOffset
		}

		var mask byte = ((1 << maskSize) - 1) << b.bitOffset

		val |= (int64(mask&curByte) >> b.bitOffset) << bitOffset
		bitOffset += maskSize
		b.bitOffset += maskSize
		nBits -= int(maskSize)
	}
	return val, nil
}

func (b *bitArray) nextBytes(dst []byte, nBits int) error {
	var byteOffset int
	var bitOffset uint

	for nBits > 0 {
		if b.bitOffset == 8 {
			b.bitOffset = 0
			b.byteOffset++

			if b.byteOffset >= len(b.bits) && nBits > 0 {
				return errors.New("Ran out of bits.")
			}
		}
		src := b.bits[b.byteOffset]

		if bitOffset == 8 {
			bitOffset = 0
			byteOffset++

			if byteOffset >= len(dst) {
				return nil
			}
		}

		maskSize := uint(nBits)
		if maskSize > (8 - bitOffset) {
			maskSize = 8 - bitOffset
		}
		if maskSize > (8 - b.bitOffset) {
			maskSize = 8 - b.bitOffset
		}
		if maskSize > 8 {
			maskSize = 8
		}

		var mask byte = ((1 << maskSize) - 1) << b.bitOffset

		/*fmt.Println("nBits:", nBits, "maskSize:", maskSize)
		fmt.Println("B.BitOffset:", b.bitOffset, "B.ByteOffset:", b.byteOffset)
		fmt.Println("BitOffset:", bitOffset, "ByteOffset:", byteOffset)
		fmt.Print("Src:  ")
		printBin(src)
		fmt.Print("Mask: ")
		printBin(mask)
		fmt.Print("Dst1: ")
		printBin(dst[byteOffset])*/

		dst[byteOffset] |= ((mask & src) >> b.bitOffset) << bitOffset
		/*fmt.Print("Dst2: ")
		printBin(dst[byteOffset])
		fmt.Println()*/
		bitOffset += maskSize
		b.bitOffset += maskSize
		nBits -= int(maskSize)
	}

	return nil
}

func printBin(b byte) {
	for j := uint(8); j > 0; j-- {
		if j%4 == 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%d", ((1<<(j-1))&b)>>(j-1))
	}
	fmt.Println()
}

func printBinUint(b uint) {
	for j := uint(32); j > 0; j-- {
		if j%4 == 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%d", ((1<<(j-1))&b)>>(j-1))
	}
	fmt.Println()
}

func printBinUint64(b uint64) {
	for j := uint(64); j > 0; j-- {
		if j%4 == 0 {
			fmt.Print(" ")
		}
		fmt.Printf("%d", ((1<<(j-1))&b)>>(j-1))
	}
	fmt.Println()
}
