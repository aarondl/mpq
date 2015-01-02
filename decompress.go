package mpq

import "errors"

const (
	compressionHuffman     = 0x01       // Huffmann compression (used on WAVE files only)
	compressionZlib        = 0x02       // ZLIB compression
	compressionPkware      = 0x08       // PKWARE DCL compression
	compressionBzip2       = 0x10       // BZIP2 compression (added in Warcraft III)
	compressionSparse      = 0x20       // Sparse compression (added in Starcraft 2)
	compressionADPCMono    = 0x40       // IMA ADPCM compression (mono)
	compressionADPCMStereo = 0x80       // IMA ADPCM compression (stereo)
	compressionLZMA        = 0x12       // LZMA compression. Added in Starcraft 2. This value is NOT a combination of flags.
	compressionNextSame    = 0xFFFFFFFF // Same compression
)

func decompress(dest []byte, src []byte) error {
	if len(dest) == len(src) {
		copy(dest, src)
		return nil
	}

	offset := 0
	compressionMethod := src[offset]
	offset++

	switch compressionMethod {
	case compressionHuffman:
		return errors.New("Huffman compression not supported.")
	case compressionLZMA:
		return errors.New("LZMA compression not supported.")
	case compressionZlib:
		return errors.New("Zlib compression not supported.")
	case compressionBzip2:
		return errors.New("Bzip2 compression not supported.")
	case compressionPkware:
		return errors.New("PKWare compression not supported.")
	case compressionSparse:
		return errors.New("Sparse compression not supported.")
	case compressionSparse | compressionZlib:
		return errors.New("Sparse+Zlib compression not supported.")
	case compressionSparse | compressionBzip2:
		return errors.New("Sparse+Bzip2 compression not supported.")
	case compressionADPCMono | compressionHuffman:
		return errors.New("ADPCMMono+Huffman compression not supported.")
	case compressionADPCMStereo | compressionHuffman:
		return errors.New("ADPCMStereo+Huffman compression not supported.")
	}

	return nil
}
