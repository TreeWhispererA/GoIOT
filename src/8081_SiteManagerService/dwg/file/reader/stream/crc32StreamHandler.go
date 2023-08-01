package stream

import (
	"bytes"
	"io"
)

func GetCRC32StreamHandler(data []byte, seed uint32) (io.Reader, error) {
	randSeed := int32(1)
	for i := 0; i < len(data); i++ {
		randSeed *= 0x343fd
		randSeed += 0x269ec3

		values := byte(randSeed >> 0x10)
		data[i] ^= values
	}

	return bytes.NewReader(data), nil
}
