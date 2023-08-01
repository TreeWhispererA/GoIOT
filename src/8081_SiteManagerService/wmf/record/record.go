package record

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
)

type WmfRecordType uint16

const (
	META_EOF                   WmfRecordType = 0x0000
	META_REALIZEPALETTE        WmfRecordType = 0x0035
	META_SETPALENTRIES         WmfRecordType = 0x0037
	META_SETBKMODE             WmfRecordType = 0x0102
	META_SETMAPMODE            WmfRecordType = 0x0103
	META_SETROP2               WmfRecordType = 0x0104
	META_SETRELABS             WmfRecordType = 0x0105
	META_SETPOLYFILLMODE       WmfRecordType = 0x0106
	META_SETSTRETCHBLTMODE     WmfRecordType = 0x0107
	META_SETTEXTCHAREXTRA      WmfRecordType = 0x0108
	META_RESTOREDC             WmfRecordType = 0x0127
	META_RESIZEPALETTE         WmfRecordType = 0x0139
	META_DIBCREATEPATTERNBRUSH WmfRecordType = 0x0142
	META_SETLAYOUT             WmfRecordType = 0x0149
	META_SETBKCOLOR            WmfRecordType = 0x0201
	META_SETTEXTCOLOR          WmfRecordType = 0x0209
	META_OFFSETVIEWPORTORG     WmfRecordType = 0x0211
	META_LINETO                WmfRecordType = 0x0213
	META_MOVETO                WmfRecordType = 0x0214
	META_OFFSETCLIPRGN         WmfRecordType = 0x0220
	META_FILLREGION            WmfRecordType = 0x0228
	META_SETMAPPERFLAGS        WmfRecordType = 0x0231
	META_SELECTPALETTE         WmfRecordType = 0x0234
	META_POLYGON               WmfRecordType = 0x0324
	META_POLYLINE              WmfRecordType = 0x0325
	META_SETTEXTJUSTIFICATION  WmfRecordType = 0x020A
	META_SETWINDOWORG          WmfRecordType = 0x020B
	META_SETWINDOWEXT          WmfRecordType = 0x020C
	META_SETVIEWPORTORG        WmfRecordType = 0x020D
	META_SETVIEWPORTEXT        WmfRecordType = 0x020E
	META_OFFSETWINDOWORG       WmfRecordType = 0x020F
	META_SCALEWINDOWEXT        WmfRecordType = 0x0410
	META_SCALEVIEWPORTEXT      WmfRecordType = 0x0412
	META_EXCLUDECLIPRECT       WmfRecordType = 0x0415
	META_INTERSECTCLIPRECT     WmfRecordType = 0x0416
	META_ELLIPSE               WmfRecordType = 0x0418
	META_FLOODFILL             WmfRecordType = 0x0419
	META_FRAMEREGION           WmfRecordType = 0x0429
	META_ANIMATEPALETTE        WmfRecordType = 0x0436
	META_TEXTOUT               WmfRecordType = 0x0521
	META_POLYPOLYGON           WmfRecordType = 0x0538
	META_EXTFLOODFILL          WmfRecordType = 0x0548
	META_RECTANGLE             WmfRecordType = 0x041B
	META_SETPIXEL              WmfRecordType = 0x041F
	META_ROUNDRECT             WmfRecordType = 0x061C
	META_PATBLT                WmfRecordType = 0x061D
	META_SAVEDC                WmfRecordType = 0x001E
	META_PIE                   WmfRecordType = 0x081A
	META_STRETCHBLT            WmfRecordType = 0x0B23
	META_ESCAPE                WmfRecordType = 0x0626
	META_INVERTREGION          WmfRecordType = 0x012A
	META_PAINTREGION           WmfRecordType = 0x012B
	META_SELECTCLIPREGION      WmfRecordType = 0x012C
	META_SELECTOBJECT          WmfRecordType = 0x012D
	META_SETTEXTALIGN          WmfRecordType = 0x012E
	META_ARC                   WmfRecordType = 0x0817
	META_CHORD                 WmfRecordType = 0x0830
	META_BITBLT                WmfRecordType = 0x0922
	META_EXTTEXTOUT            WmfRecordType = 0x0a32
	META_SETDIBTODEV           WmfRecordType = 0x0d33
	META_DIBBITBLT             WmfRecordType = 0x0940
	META_DIBSTRETCHBLT         WmfRecordType = 0x0b41
	META_STRETCHDIB            WmfRecordType = 0x0f43
	META_DELETEOBJECT          WmfRecordType = 0x01f0
	META_CREATEPALETTE         WmfRecordType = 0x00f7
	META_CREATEPATTERNBRUSH    WmfRecordType = 0x01F9
	META_CREATEPENINDIRECT     WmfRecordType = 0x02FA
	META_CREATEFONTINDIRECT    WmfRecordType = 0x02FB
	META_CREATEBRUSHINDIRECT   WmfRecordType = 0x02FC
	META_CREATEREGION          WmfRecordType = 0x06FF
)

type WmfMapMode uint16

const (
	MM_TEXT        WmfMapMode = 0x0001
	MM_LOMETRIC    WmfMapMode = 0x0002
	MM_HIMETRIC    WmfMapMode = 0x0003
	MM_LOENGLISH   WmfMapMode = 0x0004
	MM_HIENGLISH   WmfMapMode = 0x0005
	MM_TWIPS       WmfMapMode = 0x0006
	MM_ISOTROPIC   WmfMapMode = 0x0007
	MM_ANISOTROPIC WmfMapMode = 0x0008
)

type WmfMixMode uint16

const (
	TRANSPARENT WmfMixMode = 0x0001
	OPAQUE      WmfMixMode = 0x0002
)

type WmfFgMixMode uint16

const (
	R2_BLACK       WmfFgMixMode = 1
	R2_NOTMERGEPEN WmfFgMixMode = 2
	R2_MASKNOTPEN  WmfFgMixMode = 3
	R2_NOTCOPYPEN  WmfFgMixMode = 4
	R2_MASKPENNOT  WmfFgMixMode = 5
	R2_NOT         WmfFgMixMode = 6
	R2_XORPEN      WmfFgMixMode = 7
	R2_NOTMASKPEN  WmfFgMixMode = 8
	R2_MASKPEN     WmfFgMixMode = 9
	R2_NOTXORPEN   WmfFgMixMode = 10
	R2_NOP         WmfFgMixMode = 11
	R2_MERGENOTPEN WmfFgMixMode = 12
	R2_COPYPEN     WmfFgMixMode = 13
	R2_MERGEPENNOT WmfFgMixMode = 14
	R2_MERGEPEN    WmfFgMixMode = 15
	R2_WHITE       WmfFgMixMode = 16
)

type BrushStyle uint16

const (
	BS_SOLID         BrushStyle = 0x0000
	BS_NULL          BrushStyle = 0x0001
	BS_HATCHED       BrushStyle = 0x0002
	BS_PATTERN       BrushStyle = 0x0003
	BS_INDEXED       BrushStyle = 0x0004
	BS_DIBPATTERN    BrushStyle = 0x0005
	BS_DIBPATTERNPT  BrushStyle = 0x0006
	BS_PATTERN8X8    BrushStyle = 0x0007
	BS_DIBPATTERN8X8 BrushStyle = 0x0008
	BS_MONOPATTERN   BrushStyle = 0x0009
)

type HatchStyle uint16

const (
	HS_HORIZONTAL HatchStyle = 0x0000
	HS_VERTICAL   HatchStyle = 0x0001
	HS_FDIAGONAL  HatchStyle = 0x0002
	HS_BDIAGONAL  HatchStyle = 0x0003
	HS_CROSS      HatchStyle = 0x0004
	HS_DIAGCROSS  HatchStyle = 0x0005
)

type WmfBrush struct {
	Color      color.RGBA
	Style      BrushStyle
	HatchStyle HatchStyle
}

type PenStyle uint16

const (
	PS_SOLID       PenStyle = 0x0000
	PS_DASH        PenStyle = 0x0001
	PS_DOT         PenStyle = 0x0002
	PS_DASHDOT     PenStyle = 0x0003
	PS_DASHDOTDOT  PenStyle = 0x0004
	PS_NULL        PenStyle = 0x0005
	PS_INSIDEFRAME PenStyle = 0x0006
	PS_USERSTYLE   PenStyle = 0x0007
	PS_ALTERNATE   PenStyle = 0x0008
)

type WmfPen struct {
	Color color.RGBA
	Style PenStyle
	Width image.Point
}

type PaletteEntryFlag byte

const (
	PC_RESERVED   PaletteEntryFlag = 0x01
	PC_EXPLICIT   PaletteEntryFlag = 0x02
	PC_NOCOLLAPSE PaletteEntryFlag = 0x04
)

type PaletteEntry struct {
	Flag  PaletteEntryFlag
	Color color.RGBA
}

type Palette struct {
	Start   int16
	Entries []PaletteEntry
}

type CharacterSet byte

const (
	ANSI_CHARSET        CharacterSet = 0x00000000
	DEFAULT_CHARSET     CharacterSet = 0x00000001
	SYMBOL_CHARSET      CharacterSet = 0x00000002
	MAC_CHARSET         CharacterSet = 0x0000004D
	SHIFTJIS_CHARSET    CharacterSet = 0x00000080
	HANGUL_CHARSET      CharacterSet = 0x00000081
	JOHAB_CHARSET       CharacterSet = 0x00000082
	GB2312_CHARSET      CharacterSet = 0x00000086
	CHINESEBIG5_CHARSET CharacterSet = 0x00000088
	GREEK_CHARSET       CharacterSet = 0x000000A1
	TURKISH_CHARSET     CharacterSet = 0x000000A2
	VIETNAMESE_CHARSET  CharacterSet = 0x000000A3
	HEBREW_CHARSET      CharacterSet = 0x000000B1
	ARABIC_CHARSET      CharacterSet = 0x000000B2
	BALTIC_CHARSET      CharacterSet = 0x000000BA
	RUSSIAN_CHARSET     CharacterSet = 0x000000CC
	THAI_CHARSET        CharacterSet = 0x000000DE
	EASTEUROPE_CHARSET  CharacterSet = 0x000000EE
	OEM_CHARSET         CharacterSet = 0x000000FF
)

type OutPrecision byte

const (
	OUT_DEFAULT_PRECIS        OutPrecision = 0x00000000
	OUT_STRING_PRECIS         OutPrecision = 0x00000001
	OUT_CHARACTER_PRECIS      OutPrecision = 0x00000002
	OUT_STROKE_PRECIS         OutPrecision = 0x00000003
	OUT_TT_PRECIS             OutPrecision = 0x00000004
	OUT_DEVICE_PRECIS         OutPrecision = 0x00000005
	OUT_RASTER_PRECIS         OutPrecision = 0x00000006
	OUT_TT_ONLY_PRECIS        OutPrecision = 0x00000007
	OUT_OUTLINE_PRECIS        OutPrecision = 0x00000008
	OUT_SCREEN_OUTLINE_PRECIS OutPrecision = 0x00000009
	OUT_PS_ONLY_PRECIS        OutPrecision = 0x0000000A
)

type ClipPrecision byte

const (
	CLIP_DEFAULT_PRECIS   ClipPrecision = 0x00000000
	CLIP_CHARACTER_PRECIS ClipPrecision = 0x00000001
	CLIP_STROKE_PRECIS    ClipPrecision = 0x00000002
	CLIP_LH_ANGLES        ClipPrecision = 0x00000010
	CLIP_TT_ALWAYS        ClipPrecision = 0x00000020
	CLIP_DFA_DISABLE      ClipPrecision = 0x00000040
	CLIP_EMBEDDED         ClipPrecision = 0x00000080
)

type FontQuality byte

const (
	DEFAULT_QUALITY        FontQuality = 0x00
	DRAFT_QUALITY          FontQuality = 0x01
	PROOF_QUALITY          FontQuality = 0x02
	NONANTIALIASED_QUALITY FontQuality = 0x03
	ANTIALIASED_QUALITY    FontQuality = 0x04
	CLEARTYPE_QUALITY      FontQuality = 0x05
)

type PitchFont byte

const (
	DEFAULT_PITCH  PitchFont = 0
	FIXED_PITCH    PitchFont = 1
	VARIABLE_PITCH PitchFont = 2
)

type FamilyFont byte

const (
	FF_DONTCARE   FamilyFont = 0x00
	FF_ROMAN      FamilyFont = 0x01
	FF_SWISS      FamilyFont = 0x02
	FF_MODERN     FamilyFont = 0x03
	FF_SCRIPT     FamilyFont = 0x04
	FF_DECORATIVE FamilyFont = 0x05
)

type WmfFont struct {
	Height        int16
	Width         int16
	Escapement    int16
	Orientation   int16
	Weight        int16
	Italic        bool
	Underline     bool
	StrikeOut     bool
	CharSet       CharacterSet
	OutPrecision  OutPrecision
	ClipPrecision ClipPrecision
	Quality       FontQuality
	Pitch         PitchFont
	Family        FamilyFont
	FaceName      string
}

type WmfTextAlignmentMode uint16

const (
	TA_NOUPDATECP WmfTextAlignmentMode = 0x0000
	TA_UPDATECP   WmfTextAlignmentMode = 0x0001
	TA_LEFT       WmfTextAlignmentMode = 0x0000
	TA_RIGHT      WmfTextAlignmentMode = 0x0002
	TA_CENTER     WmfTextAlignmentMode = 0x0006
	TA_TOP        WmfTextAlignmentMode = 0x0000
	TA_BOTTOM     WmfTextAlignmentMode = 0x0008
	TA_BASELINE   WmfTextAlignmentMode = 0x0018
	TA_RTLREADING WmfTextAlignmentMode = 0x0100
)

type ExtTextOutOptions uint16

const (
	ETO_OPAQUE         ExtTextOutOptions = 0x0002
	ETO_CLIPPED        ExtTextOutOptions = 0x0004
	ETO_GLYPH_INDEX    ExtTextOutOptions = 0x0010
	ETO_RTLREADING     ExtTextOutOptions = 0x0080
	ETO_NUMERICSLOCAL  ExtTextOutOptions = 0x0400
	ETO_NUMERICSLATIN  ExtTextOutOptions = 0x0800
	ETO_IGNORELANGUAGE ExtTextOutOptions = 0x1000
	ETO_PDY            ExtTextOutOptions = 0x2000
)

type WmfExtTextout struct {
	X         int16
	Y         int16
	Options   uint16
	Rectangle image.Rectangle
	Value     string
}

type WmfRecord struct {
	RecordSize      uint32
	RecordSizeBytes uint16
	RecordType      WmfRecordType
	BkColor         color.Color
	BkMode          WmfMixMode
	FgMode          WmfFgMixMode
	Brush           WmfBrush
	Font            WmfFont
	Pen             WmfPen
	MapMode         WmfMapMode
	Palette         Palette
	SelectedObject  uint16
	Points          []image.Point
	TextAlign       WmfTextAlignmentMode
	ExtTextout      WmfExtTextout
	WindowOrigin    image.Point
	WindowExt       image.Point
}

func readColor(reader io.Reader) (color.RGBA, error) {
	c := color.RGBA{}
	err := error(nil)

	var r, g, b, a byte
	err = binary.Read(reader, binary.LittleEndian, &r)
	err = binary.Read(reader, binary.LittleEndian, &g)
	err = binary.Read(reader, binary.LittleEndian, &b)
	err = binary.Read(reader, binary.LittleEndian, &a)
	if err == nil {
		c = color.RGBA{R: r, G: g, B: b, A: 0xff}
	}

	return c, err
}

func readPoint(reader io.Reader) (image.Point, error) {
	var x, y int16
	var pt image.Point
	err := error(nil)

	err = binary.Read(reader, binary.LittleEndian, &x)
	if err == nil {
		err = binary.Read(reader, binary.LittleEndian, &y)
	}
	if err == nil {
		pt = image.Point{X: int(x), Y: int(y)}
	}

	return pt, err
}

func readRectangle(reader io.Reader) (image.Rectangle, error) {
	var rectangle image.Rectangle
	err := error(nil)

	var left, top, right, bottom int16
	err = binary.Read(reader, binary.LittleEndian, &left)
	err = binary.Read(reader, binary.LittleEndian, &top)
	err = binary.Read(reader, binary.LittleEndian, &right)
	err = binary.Read(reader, binary.LittleEndian, &bottom)
	if err == nil {
		rectangle = image.Rectangle{
			Min: image.Point{X: int(left), Y: int(bottom)},
			Max: image.Point{X: int(right), Y: int(top)},
		}
	}

	return rectangle, err
}

func bytesToString(bytes []byte) string {
	n := len(bytes)
	for i := 0; i < n; i++ {
		if bytes[i] == 0 {
			n = i
			break
		}
	}
	return string(bytes[:n])
}

func Read(file *os.File) (*WmfRecord, error) {
	err := error(nil)
	var begin int64
	var end int64

	begin, err = file.Seek(0, io.SeekCurrent)

	var recordSize uint32
	var recordType WmfRecordType
	err = binary.Read(file, binary.LittleEndian, &recordSize)
	if err != nil {
		return nil, err
	}
	err = binary.Read(file, binary.LittleEndian, &recordType)
	if err != nil {
		return nil, err
	}

	record := WmfRecord{RecordSize: recordSize, RecordType: recordType}

	switch record.RecordType {
	case META_CREATEPENINDIRECT:
		record.Pen = WmfPen{}

		err = binary.Read(file, binary.LittleEndian, &record.Pen.Style)
		if err == nil {
			record.Pen.Width, err = readPoint(file)
		}
		if err == nil {
			record.Pen.Color, err = readColor(file)
		}
	case META_CREATEBRUSHINDIRECT:
		record.Brush = WmfBrush{}

		err = binary.Read(file, binary.LittleEndian, &record.Brush.Style)
		if err == nil {
			record.Brush.Color, err = readColor(file)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.Brush.HatchStyle)
		}
	case META_CREATEFONTINDIRECT:
		record.Font = WmfFont{}

		var b byte
		err = binary.Read(file, binary.LittleEndian, &record.Font.Height)
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.Font.Width)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.Font.Escapement)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.Font.Orientation)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.Font.Weight)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &b)
		}
		if err == nil {
			record.Font.Italic = b == 1
			err = binary.Read(file, binary.LittleEndian, &b)
		}
		if err == nil {
			record.Font.Underline = b == 1
			err = binary.Read(file, binary.LittleEndian, &b)
		}
		if err == nil {
			record.Font.StrikeOut = b == 1
			err = binary.Read(file, binary.LittleEndian, &record.Font.CharSet)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.Font.OutPrecision)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.Font.ClipPrecision)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.Font.Quality)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &b)
			record.Font.Pitch = PitchFont(b & 0x0f)
			record.Font.Family = FamilyFont(b >> 4)
		}
		if err == nil {
			left := 2*record.RecordSize - 6 /* RecordSize and RecordFunction */ - 18 /* from Height to PitchAndFamily */
			if left > 0 {
				buffer := make([]byte, left)
				err = binary.Read(file, binary.LittleEndian, buffer)
				if err == nil {
					record.Font.FaceName = bytesToString(buffer)
				}
			} else {
				record.Font.FaceName = ""
			}
		}

	case META_CREATEPALETTE:
		var numEntries int16
		record.Palette = Palette{}

		err = binary.Read(file, binary.LittleEndian, &record.Palette.Start)
		if err == nil {
			numEntries = 0
			err = binary.Read(file, binary.LittleEndian, &numEntries)
		}
		for i := int16(0); i < numEntries; i++ {
			var flag PaletteEntryFlag
			var r, g, b byte
			err = binary.Read(file, binary.LittleEndian, &flag)
			err = binary.Read(file, binary.LittleEndian, &b)
			err = binary.Read(file, binary.LittleEndian, &g)
			err = binary.Read(file, binary.LittleEndian, &r)
			if err == nil {
				record.Palette.Entries = append(record.Palette.Entries, PaletteEntry{
					Flag: flag, Color: color.RGBA{R: r, G: g, B: b, A: 0xff}})
			} else {
				break
			}
		}
	case META_POLYGON, META_POLYLINE:
		record.Points = []image.Point{}

		var numPoints int16
		err = binary.Read(file, binary.LittleEndian, &numPoints)
		if err != nil {
			numPoints = 0
		}

		var pt image.Point
		for i := int16(0); i < numPoints; i++ {
			pt, err = readPoint(file)
			if err == nil {
				record.Points = append(record.Points, pt)
			} else {
				break
			}
		}
	case META_SELECTOBJECT, META_DELETEOBJECT:
		err = binary.Read(file, binary.LittleEndian, &record.SelectedObject)
	case META_SELECTPALETTE:
		err = binary.Read(file, binary.LittleEndian, &record.SelectedObject)
	case META_SETBKCOLOR:
		record.BkColor, err = readColor(file)
	case META_SETBKMODE:
		err = binary.Read(file, binary.LittleEndian, &record.BkMode)
	case META_SETMAPMODE:
		err = binary.Read(file, binary.LittleEndian, &record.MapMode)
	case META_SETROP2:
		err = binary.Read(file, binary.LittleEndian, &record.FgMode)
	case META_SETTEXTALIGN:
		err = binary.Read(file, binary.LittleEndian, &record.TextAlign)
	case META_SETTEXTCOLOR:
		record.Pen = WmfPen{}
		record.Pen.Color, err = readColor(file)
	case META_SETPIXEL:
		record.Pen = WmfPen{}
		record.Points = []image.Point{}

		record.Pen.Color, err = readColor(file)
		var x, y int16
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &y)
			err = binary.Read(file, binary.LittleEndian, &x)
			if err == nil {
				record.Points = append(record.Points, image.Point{X: int(x), Y: int(y)})
			}
		}

	case META_SETWINDOWORG:
		record.WindowOrigin, err = readPoint(file)
	case META_SETWINDOWEXT:
		var ext image.Point
		ext, err = readPoint(file)
		if err == nil {
			ext = image.Point{X: int(ext.Y), Y: int(ext.X)}
		}
		record.WindowExt = ext
	case META_EXTTEXTOUT:
		record.ExtTextout = WmfExtTextout{}

		var stringLength int32

		err = binary.Read(file, binary.LittleEndian, &record.ExtTextout.Y)
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &record.ExtTextout.X)
		}
		if err == nil {
			err = binary.Read(file, binary.LittleEndian, &stringLength)
		} else {
			stringLength = 0
		}
		if stringLength > 0 {
			bytesValue := make([]byte, stringLength)
			err = binary.Read(file, binary.LittleEndian, &bytesValue)
			if err == nil {
				record.ExtTextout.Value = bytesToString(bytesValue)
			}
		}
		if err == nil {
			record.ExtTextout.Rectangle, err = readRectangle(file)
		}
	case META_REALIZEPALETTE:
	case META_EOF:
	default:
		err = fmt.Errorf("Unsupported record type: %X", record.RecordType)
	}

	end, _ = file.Seek(0, io.SeekCurrent)

	rlen := end - begin
	excess := 2*int64(record.RecordSize) - rlen
	if excess > 0 {
		file.Seek(excess, io.SeekCurrent)
	}

	return &record, err
}
