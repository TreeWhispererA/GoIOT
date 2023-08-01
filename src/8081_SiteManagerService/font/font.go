package font

import (
	"io/ioutil"

	"github.com/flopp/go-findfont"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

func getFontFileName(fontData draw2d.FontData) string {
	fontFileName := fontData.Name
	if fontData.Style&draw2d.FontStyleBold != 0 {
		fontFileName += "b"
	}

	if fontData.Style&draw2d.FontStyleItalic != 0 {
		fontFileName += "i"
	}

	fontFileName += ".ttf"
	return fontFileName
}

// FolderFontCache can Load font from folder
type FolderFontCache struct {
	fonts map[string]*truetype.Font
}

// Load a font from cache if exists otherwise it will load the font from file
func (cache *FolderFontCache) Load(fontData draw2d.FontData) (font *truetype.Font, err error) {
	fontFile := getFontFileName(fontData)

	if font = cache.fonts[fontFile]; font != nil {
		return font, nil
	}

	fontPath, err := findfont.Find(fontFile)
	if err != nil {
		return nil, err
	}

	var data []byte
	if data, err = ioutil.ReadFile(fontPath); err != nil {
		return
	}

	if font, err = truetype.Parse(data); err != nil {
		return
	}

	cache.fonts[fontFile] = font
	return
}

// Store a font to this cache
func (cache *FolderFontCache) Store(fontData draw2d.FontData, font *truetype.Font) {
	cache.fonts[getFontFileName(fontData)] = font
}

func UpdateDraw2dFontSettings() {
	draw2d.SetFontCache(&FolderFontCache{
		fonts: make(map[string]*truetype.Font),
	})
}
