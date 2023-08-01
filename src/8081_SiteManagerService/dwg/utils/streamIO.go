package utils

import (
	"io"
)

func ReadString(stream io.Reader, length int, codePage CodePage) (value string, err error) {
	value = ""

	if length <= 0 {
		return value, nil
	}

	buf := make([]byte, length)
	if _, err = stream.Read(buf); err != nil {
		return value, err
	}

	// var decoder *encoding.Decoder
	// if decoder, err = GetListedEncoding(codePage); err != nil {
	// 	return value, err
	// }

	return string(value), nil
}

func BytesToString(bytes []byte) string {
	n := len(bytes)
	for i := 0; i < n; i++ {
		if bytes[i] == 0 {
			n = i
			break
		}
	}
	return string(bytes[:n])
}
