package reader

import (
	"io"

	"tracio.com/sitemanagerservice/dwg/classes"
	"tracio.com/sitemanagerservice/dwg/file/reader/stream"
	"tracio.com/sitemanagerservice/dwg/file/section"
	"tracio.com/sitemanagerservice/dwg/file/version"
)

type DwgObjectReader struct {
	DwgSectionIO

	crcStream  io.ReadSeeker
	crcReader  stream.IDwgStreamReader
	sreader    stream.IDwgStreamReader
	handles    []uint32
	handleMap  map[uint32]int32
	classesMap map[string]*classes.DxfClass
}

func (this *DwgObjectReader) Read() (map[uint32]int32, error) {
	return nil, nil
}

func NewDwgObjectReader(dwgVersion version.ACadVersion,
	sreader stream.IDwgStreamReader, handles []uint32, handleMap map[uint32]int32,
	classesMap map[string]*classes.DxfClass) (*DwgObjectReader, error) {
	reader := &DwgObjectReader{
		DwgSectionIO: DwgSectionIO{
			DwgVersion:  dwgVersion,
			SectionName: section.AcDbObjects,
		},
		sreader:    sreader,
		handles:    handles,
		handleMap:  handleMap,
		classesMap: classesMap,
	}

	err := error(nil)

	//Initialize the crc stream
	//RS : CRC for the data section, starting after the sentinel. Use 0xC0C1 for the initial value
	// if this._builder.Configuration.CrcCheck {
	// 	this._crcStream = new CRC8StreamHandler(this._reader.Stream, 0xC0C1);
	// } else {
	if reader.crcStream, err = stream.CloneStream(sreader.GetStream()); err != nil {
		return reader, err
	}
	// }

	if _, err = reader.crcStream.Seek(0, io.SeekStart); err != nil {
		return reader, err
	}

	//Setup the entity handler
	reader.crcReader, err = stream.NewDwgStreamHandler(reader.DwgVersion, reader.crcStream)

	return reader, err
}
