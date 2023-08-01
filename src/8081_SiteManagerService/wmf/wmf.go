package wmf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"

	"tracio.com/sitemanagerservice/wmf/record"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

type WmfMetafileType uint16

const (
	MEMORYMETAFILE WmfMetafileType = 0x0001
	DISKMETAFILE   WmfMetafileType = 0x0002
)

type WmfMetafileVersion uint16

const (
	METAVERSION100 WmfMetafileVersion = 0x0100
	METAVERSION300 WmfMetafileVersion = 0x0300
)

type WmfFormat struct {
	Handle        uint16
	Left          int16
	Top           int16
	Right         int16
	Bottom        int16
	PixelsPerInch uint16
	Reserved      uint32
	Checksum      uint16
}

type WmfHeader struct {
	FileType   WmfMetafileType
	HeaderSize uint16
	Version    WmfMetafileVersion
	FileSize   uint32
	NumObjects uint16
	MaxRecord  uint32
	NumMembers uint16
}

type WmfFile struct {
	Format  WmfFormat
	Header  WmfHeader
	records []*record.WmfRecord
}

func loadWmfFileHeader(file io.Reader, wmf *WmfFile) error {
	err := binary.Read(file, binary.LittleEndian, &wmf.Format)
	err = binary.Read(file, binary.LittleEndian, &wmf.Header)

	return err
}

func loadWmfFile(fileName string) (*WmfFile, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	wmf := WmfFile{}

	var key uint32
	binary.Read(file, binary.LittleEndian, &key)
	if key != 0x9ac6cdd7 {
		return nil, errors.New("WMF key does not match the pattern")
	}

	loadWmfFileHeader(file, &wmf)

	for true {
		rec, err := record.Read(file)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				fmt.Println(err)
			}
		}

		wmf.records = append(wmf.records, rec)
		if rec.RecordType == record.META_EOF {
			break
		}
	}

	return &wmf, nil
}

func LoadAsImage(fileName string, use_default_dpi bool) (*image.RGBA, error) {
	wmfFile, err := loadWmfFile(fileName)
	if err != nil {
		return nil, err
	}

	scale := 1.0
	dpi := 96
	if use_default_dpi {
		dpi = int(wmfFile.Format.PixelsPerInch)
	} else {
		scale = float64(dpi) / float64(wmfFile.Format.PixelsPerInch)
	}
	DEFAULT_LINE_WIDTH := 1 / scale

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: int(math.Floor(float64(wmfFile.Format.Left) * scale)),
			Y: int(math.Floor(float64(wmfFile.Format.Top) * scale)),
		},
		Max: image.Point{
			X: int(math.Ceil(float64(wmfFile.Format.Right) * scale)),
			Y: int(math.Ceil(float64(wmfFile.Format.Bottom) * scale)),
		},
	})
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetDPI(dpi)
	gc.SetMatrixTransform(draw2d.NewScaleMatrix(scale, scale))

	WHITE_COLOR := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	BLANK_COLOR := color.RGBA{R: 0, G: 0, B: 0, A: 0}

	// Clear background
	gc.SetFillColor(WHITE_COLOR)
	gc.Clear()

	resources := []*record.WmfRecord{}

	for _, rec := range wmfFile.records {
		switch rec.RecordType {
		case record.META_SETMAPMODE:
			continue
		case record.META_SETWINDOWORG:
			// Set origin
			continue
		case record.META_SETWINDOWEXT:
			// Set Ext
			continue
		case record.META_SETBKCOLOR:
			continue
		case record.META_SETBKMODE:
			// fmt.Printf("SetBkMode(%d)\n", rec.BkMode)
		case record.META_SETROP2:
			// fmt.Printf("SetROP2(%d)\n", rec.FgMode)
		case record.META_CREATEPALETTE, record.META_CREATEBRUSHINDIRECT, record.META_CREATEPENINDIRECT, record.META_CREATEFONTINDIRECT:
			var i int
			for i = 0; i < len(resources); i++ {
				if resources[i] == nil {
					resources[i] = rec
					break
				}
			}
			if i == len(resources) {
				resources = append(resources, rec)
			}
		case record.META_SELECTPALETTE:
			continue
		case record.META_SELECTOBJECT:
			rec = resources[rec.SelectedObject]
			if rec.RecordType == record.META_CREATEPENINDIRECT {
				gc.SetStrokeColor(rec.Pen.Color)
				if rec.Pen.Width.X == 0 {
					gc.SetLineWidth(DEFAULT_LINE_WIDTH)
				} else {
					gc.SetLineWidth(float64(rec.Pen.Width.X))
				}
			} else if rec.RecordType == record.META_CREATEBRUSHINDIRECT {
				if rec.Brush.Style == record.BS_NULL {
					gc.SetFillColor(BLANK_COLOR)
				} else if rec.Brush.Style == record.BS_SOLID {
					gc.SetFillColor(rec.Brush.Color)
				} else if rec.Brush.Style == record.BS_HATCHED {
					fmt.Printf("BrushStyle is not implemented: %d\n", int(rec.Brush.Style))
				}
			} else if rec.RecordType == record.META_CREATEFONTINDIRECT {
				fontStyle := draw2d.FontStyleNormal
				if rec.Font.Italic == true {
					fontStyle = draw2d.FontStyleItalic
				}
				fontData := draw2d.FontData{Name: rec.Font.FaceName, Family: draw2d.FontFamilySerif, Style: fontStyle}
				gc.SetFontData(fontData)
				gc.SetFontSize(math.Abs(float64(rec.Font.Height)))
			}
		case record.META_DELETEOBJECT:
			resources[rec.SelectedObject] = nil
		case record.META_POLYGON, record.META_POLYLINE:
			for i, pt := range rec.Points {
				if i == 0 {
					gc.MoveTo(float64(pt.X), float64(pt.Y))
				} else {
					gc.LineTo(float64(pt.X), float64(pt.Y))
				}
			}

			if rec.RecordType == record.META_POLYGON {
				gc.Close()
				gc.FillStroke()
			} else {
				gc.Stroke()
			}
		case record.META_SETPIXEL:
			pt := rec.Points[0]
			gc.MoveTo(float64(pt.X), float64(pt.Y))
			gc.LineTo(float64(pt.X), float64(pt.Y))
			gc.Stroke()
		case record.META_SETTEXTALIGN:
			// fmt.Printf("SetTextAlign(%d)\n", rec.TextAlign)
		case record.META_SETTEXTCOLOR:
			gc.SetStrokeColor(rec.Pen.Color)
		case record.META_EXTTEXTOUT:
			gc.FillStringAt(rec.ExtTextout.Value,
				float64(rec.ExtTextout.X), float64(rec.ExtTextout.Y))
		case record.META_REALIZEPALETTE:
		case record.META_EOF:
		default:
			fmt.Printf("RecordType %d is not implemented\n", rec.RecordType)
		}
	}

	return img, nil
}
