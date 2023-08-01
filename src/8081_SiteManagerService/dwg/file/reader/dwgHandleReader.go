package reader

import (
	"encoding/binary"
	"log"
	"math"

	"tracio.com/sitemanagerservice/dwg/file/reader/stream"
	"tracio.com/sitemanagerservice/dwg/file/section"
	"tracio.com/sitemanagerservice/dwg/file/version"
)

type DwgHandleReader struct {
	DwgSectionIO

	sreader stream.IDwgStreamReader
}

func (this *DwgHandleReader) Read() (map[uint32]int32, error) {
	objectMap := make(map[uint32]int32)
	err := error(nil)

	for true {
		//Set the "last handle" to all 0 and the "last loc" to 0L;
		lasthandle := uint32(0)
		lastloc := int32(0)

		//Short: size of this section. Note this is in BIGENDIAN order (MSB first)
		var size int16
		if err = binary.Read(this.sreader.GetStream(), binary.BigEndian, &size); err != nil {
			return nil, err
		} else if size == 2 {
			break
		}

		var startPos, curPos int64

		startPos, err = this.sreader.GetPosition()
		//Note that each section is cut off at a maximum length of 2032.
		maxSectionOffset := int64(math.Min(float64(size-2), 2032))
		lastPosition := startPos + maxSectionOffset

		//Repeat until out of data for this section:
		for true {
			if curPos, err = this.sreader.GetPosition(); err != nil {
				return nil, err
			} else if curPos >= lastPosition {
				break
			}

			//offset of this handle from last handle as modular char.
			var offset uint32
			if offset, err = this.sreader.ReadModularChar(); err != nil {
				return nil, err
			}
			lasthandle += offset

			//offset of location in file from last loc as modular char. (note
			//that location offsets can be negative, if the terminating byte
			//has the 4 bit set).
			var soffset int32
			if soffset, err = this.sreader.ReadSignedModularChar(); err != nil {
				return nil, err
			}
			lastloc += soffset

			if offset > 0 {
				objectMap[lasthandle] = lastloc
			} else {
				//0 offset, wrong reference
				log.Printf("Negative offset: {%d} for the handle: {%d}\n", offset, lasthandle)
			}
		}

		//CRC (most significant byte followed by least significant byte)
		this.sreader.ReadByte()
		this.sreader.ReadByte()
	}

	return objectMap, err
}

func NewDwgHandleReader(dwgVersion version.ACadVersion, sreader stream.IDwgStreamReader) *DwgHandleReader {
	return &DwgHandleReader{
		DwgSectionIO: DwgSectionIO{
			DwgVersion:  dwgVersion,
			SectionName: section.Handles,
		},
		sreader: sreader,
	}
}
