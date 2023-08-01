package utils

import (
	"fmt"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func GetListedEncoding(code CodePage) (*encoding.Decoder, error) {
	var decoder *encoding.Decoder
	err := error(nil)

	switch code {
	case CP_Unknown:
		decoder = charmap.CodePage037.NewDecoder()
	case CP_Windows1252:
		decoder = charmap.Windows1252.NewDecoder()
	default:
		err = fmt.Errorf("Unable to find decoder for Codepage: %d", code)
	}

	return decoder, err
}
