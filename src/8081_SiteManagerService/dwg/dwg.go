package dwg

import (
	"fmt"
	"image"

	"tracio.com/sitemanagerservice/dwg/file"
)

func LoadAsImage(fileName string) (*image.RGBA, error) {
	dwgFile := file.DwgFile{}
	if doc, err := dwgFile.Open(fileName); err != nil {
		return nil, err
	} else if doc == nil {
		return nil, fmt.Errorf("Empty Document")
	}

	return nil, nil
}
