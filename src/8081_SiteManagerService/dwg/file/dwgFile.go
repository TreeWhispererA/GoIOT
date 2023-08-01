package file

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"tracio.com/sitemanagerservice/dwg/classes"
	"tracio.com/sitemanagerservice/dwg/document"
	"tracio.com/sitemanagerservice/dwg/file/builder"
	"tracio.com/sitemanagerservice/dwg/file/header"
	"tracio.com/sitemanagerservice/dwg/file/reader"
	"tracio.com/sitemanagerservice/dwg/file/reader/stream"
	"tracio.com/sitemanagerservice/dwg/file/section"
	"tracio.com/sitemanagerservice/dwg/file/version"
)

type DwgFile struct {
	Builder    *builder.DwgDocumentBuilder
	FileHeader *header.DwgFileHeader
	Document   *document.CadDocument
}

func (dwg *DwgFile) decryptDataSection(section *header.DwgLocalSectionMap, reader stream.IDwgStreamReader) (err error) {
	var position int64
	var value int32

	if position, err = reader.GetPosition(); err != nil {
		return err
	}

	secMask := uint32(position ^ 0x4164536B)

	//0x00	4	Section page type, since it’s always a data section: 0x4163043b
	if _, err = reader.ReadRawLong(); err != nil {
		return err
	} else {
		// pageType := value ^ secMask
	}
	//0x04	4	Section number
	if _, err = reader.ReadRawLong(); err != nil {
		return err
	} else {
		// sectionNumber := value ^ secMask
	}
	//0x08	4	Data size (compressed)
	if value, err = reader.ReadRawLong(); err != nil {
		return err
	} else {
		section.CompressedSize = uint32(value) ^ secMask
	}
	//0x0C	4	Page Size (decompressed)
	if value, err = reader.ReadInt(); err != nil {
		return err
	} else {
		section.PageSize = uint32(value) ^ secMask
	}
	//0x10	4	Start Offset (in the decompressed buffer)
	if value, err = reader.ReadInt(); err != nil {
		return err
	} else {
		section.Offset = uint64(uint32(value) ^ secMask)
	}
	//0x14	4	Page header Checksum (section page checksum calculated from unencoded header bytes, with the data checksum as seed)
	if value, err = reader.ReadInt(); err != nil {
		return err
	} else {
		section.Offset += uint64(uint32(value) ^ secMask)
	}

	//0x18	4	Data Checksum (section page checksum calculated from compressed data bytes, with seed 0)
	if value, err = reader.ReadInt(); err != nil {
		return err
	} else {
		section.Checksum = uint32(value) ^ secMask
	}
	//0x1C	4	Unknown (ODA writes a 0)
	if value, err = reader.ReadInt(); err != nil {
		return err
	} else {
		section.ODA = uint32(value) ^ secMask
	}

	return nil
}

func (dwg *DwgFile) getSectionBuffer18(r io.ReadSeeker, sectionName string) (sectionStream *stream.MemoryStream, err error) {
	descriptor, ok := dwg.FileHeader.Descriptors[sectionName]
	if !ok {
		return nil, nil
	}

	buf := make([]byte, int(descriptor.DecompressedSize)*len(descriptor.LocalSections))
	memoryStream := stream.NewMemoryStream(buf)

	for _, section := range descriptor.LocalSections {
		if section.IsEmpty {
			// Page is empty, fill the gap with 0s
			b := byte(0)
			for i := 0; i < int(section.DecompressedSize); i++ {
				binary.Write(memoryStream, binary.LittleEndian, b)
			}
		} else {
			//Get the page section header
			var sreader stream.IDwgStreamReader
			if sreader, err = stream.NewDwgStreamHandler(dwg.FileHeader.Version, r); err != nil {
				return nil, err
			}
			sreader.SetPosition(int64(section.Seeker))
			//Get the header data
			if err = dwg.decryptDataSection(section, sreader); err != nil {
				return nil, err
			}

			if descriptor.IsCompressed() {
				decompressor := stream.DwgLZ77AC18Decompressor{}
				if err = decompressor.Decompress(r, memoryStream); err != nil {
					return nil, err
				}
			} else {
				// Read the stream normally
				var buffer []byte
				if buffer, err = sreader.ReadBytes(int(section.CompressedSize)); err != nil {
					return nil, err
				} else if _, err = memoryStream.Write(buffer); err != nil {
					return nil, err
				}
			}
		}
	}

	// Reset the stream
	memoryStream.Seek(0, io.SeekStart)

	return memoryStream, nil
}

func (dwg *DwgFile) getSectionStream(r io.ReadSeeker, sectionDefinition section.DwgSectionDefinition) (stream.IDwgStreamReader, error) {
	var sectionStream io.ReadSeeker
	err := error(nil)
	switch dwg.FileHeader.Version {
	case version.AC1024, version.AC1027, version.AC1032:
		sectionStream, err = dwg.getSectionBuffer18(r, string(sectionDefinition))
	default:
		err = fmt.Errorf("Unsupported AutoCAD version: %d", dwg.FileHeader.Version)
	}

	// Section not found
	if sectionStream == nil {
		return nil, err
	}

	var streamHandler stream.IDwgStreamReader
	if streamHandler, err = stream.NewDwgStreamHandler(dwg.FileHeader.Version, sectionStream); err != nil {
		return nil, err
	}

	return streamHandler, nil
}

func (dwg *DwgFile) readHeader(r io.ReadSeeker) error {
	err := error(nil)

	dwg.FileHeader = header.NewDwgFileHeader()
	if err = dwg.FileHeader.Read(r); err != nil {
		return err
	}

	dwg.Document.Header = document.NewCadHeader()
	dwg.Document.Header.CodePage = dwg.FileHeader.MetaData.DrawingCodePage

	var sreader stream.IDwgStreamReader
	if sreader, err = dwg.getSectionStream(r, section.Header); err != nil {
		return err
	}

	var headerHandles *header.DwgHeaderHandlesCollection
	maintenanceVersion := dwg.FileHeader.MetaData.AcadMaintenanceVersion
	hReader := reader.NewDwgHeaderReader(dwg.FileHeader.Version, maintenanceVersion, sreader)
	if headerHandles, err = hReader.Read(dwg.Document.Header); err != nil {
		return err
	} else if dwg.Builder != nil {
		dwg.Builder.HeaderHandles = headerHandles
	}

	return nil
}

func (dwg *DwgFile) readClasses(r io.ReadSeeker) (map[int16]*classes.DxfClass, error) {
	err := error(nil)

	var sreader stream.IDwgStreamReader
	if sreader, err = dwg.getSectionStream(r, section.Classes); err != nil {
		return nil, err
	}

	var classes map[int16]*classes.DxfClass
	handler := reader.NewDwgClassesReader(dwg.FileHeader.Version,
		dwg.FileHeader.MetaData.AcadMaintenanceVersion, sreader)
	if classes, err = handler.Read(); err != nil {
		return nil, err
	}

	return classes, nil
}

func (dwg *DwgFile) readHandles(r io.ReadSeeker) (map[uint32]int32, error) {
	err := error(nil)

	var sreader stream.IDwgStreamReader
	if sreader, err = dwg.getSectionStream(r, section.Handles); err != nil {
		return nil, err
	}

	var handles map[uint32]int32
	handler := reader.NewDwgHandleReader(dwg.FileHeader.Version, sreader)
	if handles, err = handler.Read(); err != nil {
		return nil, err
	}

	return handles, nil
}

func (dwg *DwgFile) readObjects(r io.ReadSeeker) error {
	err := error(nil)

	var handles map[uint32]int32
	if handles, err = dwg.readHandles(r); err != nil {
		return err
	}
	if dwg.Document.Classes, err = dwg.readClasses(r); err != nil {
		return err
	}

	var sreader stream.IDwgStreamReader
	if dwg.FileHeader.Version <= version.AC1015 {
		return fmt.Errorf("Unsupported ACad version: %d", dwg.FileHeader.Version)
	} else {
		if sreader, err = dwg.getSectionStream(r, section.AcDbObjects); err != nil {
			return err
		}
	}

	var sectionReader *reader.DwgObjectReader
	sectionReader, err = reader.NewDwgObjectReader(dwg.FileHeader.Version, sreader, nil, handles, nil)

	// var objects map[uint32]int32
	if _, err = sectionReader.Read(); err != nil {
		return err
	}

	return nil
}

// func (dwg *DwgFile) readSummaryInfo(file *os.File) (err error) {
// 	// // Older versions than 2004 don't have summaryinfo in it's file
// 	// if dwg.FileHeader.Version < version.AC1018 {
// 	// 	return nil
// 	// }

// 	// var reader io.Reader
// 	// if reader, err = dwg.getSectionStream(file, section.SummaryInfo); err != nil {
// 	// 	return err
// 	// } else if reader != nil {
// 	// 	return nil
// 	// }

// 	return nil
// }

func (dwg *DwgFile) Open(fileName string) (*document.CadDocument, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	dwg.Document = document.NewCadDocument()
	dwg.Builder = builder.NewDwgDocumentBuilder(dwg.Document)

	// Read the file header
	if err = dwg.readHeader(file); err != nil {
		return nil, err
	}

	// Read all the objects in the file
	if err = dwg.readObjects(file); err != nil {
		return nil, err
	}

	return dwg.Document, nil
}
