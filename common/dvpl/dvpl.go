package dvpl

import (
	"encoding/binary"
	"errors"
	"hash/crc32"

	"github.com/pierrec/lz4/v4"
)

func EncryptDVPL(inputBuf []byte) ([]byte, error) {
	inputSize := len(inputBuf)

	compressedBuf := make([]byte, lz4.CompressBlockBound(inputSize))
	actualCompressedSize, err := lz4.CompressBlock(inputBuf, compressedBuf, nil)
	if err != nil {
		return nil, errors.New("failed to compress lz4")
	}

	compressedBuf = compressedBuf[:actualCompressedSize]

	var outputBuf []byte

	// Skip compression when compressed size > input size
	if inputSize < actualCompressedSize {
		outputBuf = binary.LittleEndian.AppendUint32(inputBuf, uint32(inputSize))
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, uint32(inputSize))

		crc32DataSum := crc32.ChecksumIEEE(inputBuf)
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, crc32DataSum)

		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, 0)
	} else {
		outputBuf = binary.LittleEndian.AppendUint32(compressedBuf, uint32(inputSize))
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, uint32(actualCompressedSize))

		crc32DataSum := crc32.ChecksumIEEE(compressedBuf)
		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, crc32DataSum)

		outputBuf = binary.LittleEndian.AppendUint32(outputBuf, 1)
	}
	outputBuf = append(outputBuf, []byte("DVPL")...)
	return outputBuf, nil
}

func DecryptDVPL(inputBuf []byte) ([]byte, error) {
	dataBuf := inputBuf[:len(inputBuf)-20]
	footerBuf := inputBuf[len(inputBuf)-20:]

	originalSize := binary.LittleEndian.Uint32(footerBuf[:4])
	compressedSize := binary.LittleEndian.Uint32(footerBuf[4:8])
	if int(compressedSize) != len(dataBuf) {
		return nil, errors.New("invalid compressed data length")
	}

	crc32DataSum := binary.LittleEndian.Uint32(footerBuf[8:12])
	if crc32DataSum != crc32.ChecksumIEEE(dataBuf) {
		return nil, errors.New("invalid crc32 sum")
	}

	compressType := binary.LittleEndian.Uint32(footerBuf[12:16])
	outputBuf := make([]byte, originalSize)
	if compressType == 0 {
		outputBuf = dataBuf
	} else {
		actualOutputSize, err := lz4.UncompressBlock(dataBuf, outputBuf)
		if err != nil {
			return nil, errors.New("failed to uncompressed lz4")
		}
		outputBuf = outputBuf[:actualOutputSize]
	}
	return outputBuf, nil
}
