package header

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"tracio.com/sitemanagerservice/dwg/file/reader/stream"
	"tracio.com/sitemanagerservice/dwg/file/version"
	"tracio.com/sitemanagerservice/dwg/utils"
)

type DwgFileMetaData struct {
	AcadMaintenanceVersion byte
	PreviewAddress         int32
	DwgVersion             byte
	AppReleaseVersion      byte
	DrawingCodePage        utils.CodePage
	SecurityType           int32
	SummaryInfoAddr        int32
	VbaProjectAddr         int32
	AppInfoAddress         int32
}

type DwgFileHeaderData struct {
	FileId               [12]byte // "AcFssFcAJMB\0"
	Unknown1             int32    // 0x00
	Unknown2             int32    // 0x6C
	Unknown3             int32    // 0x04
	RootTreeNodeGap      int32    // Root tree node gap
	LeftGap              int32    // Lowermost left tree node gap
	RightGap             int32    // Lowermost right tree node gap
	Unknown4             int32    // Unknown long(ODA writes 1)
	LastPageId           int32    // Last section page Id
	LastSectionAddr      uint64   // Last section page end address
	SecondHeaderAddr     uint64   // Second header data address pointing to the repeated header data at the end of the file
	GapAmount            uint32   // Gap amount
	SectionAmount        uint32   // Section page amount
	Unknown5             int32    // 0x20
	Unknown6             int32    // 0x80
	Unknown7             int32    // 0x40
	SectionPageMapId     uint32   // Section Page Map Id
	PageMapAddress       uint64   // Section Page Map address(add 0x100 to this value)
	SectionMapId         uint32   // Section Map Id
	SectionArrayPageSize uint32   // Section page array size
	GapArraySize         uint32   // Gap array size
	CRCSeed              uint32   // CRC32
}

type DwgSectionLocatorRecord struct {
	Number int32
	Size   int32
	Seeker int32
}

type DwgLocalSectionMap struct {
	IsEmpty          bool
	PageNumber       int32
	CompressedSize   uint32
	Offset           uint64
	DecompressedSize uint32
	Seeker           int32
	Size             int32
	Checksum         uint32
	CRC              uint32
	PageSize         uint32
	ODA              uint32
	SectionMap       int32
}

type DwgSectionDescriptor struct {
	CompressedSize   uint64
	PageCount        int32
	DecompressedSize uint32
	Unknown1         int32
	CompressedCode   int32
	SectionId        int32
	Encrypted        int32
	Name             string
	LocalSections    []*DwgLocalSectionMap
}

func (descriptor *DwgSectionDescriptor) IsCompressed() bool {
	return descriptor.CompressedCode == 2
}

type PageHeaderData struct {
	SectionType      int32
	DecompressedSize int32
	CompressedSize   int32
	CompressionType  int32
	Checksum         uint32
}

type DwgFileHeader struct {
	Version     version.ACadVersion
	MetaData    DwgFileMetaData
	HeaderData  DwgFileHeaderData
	Records     map[int32]*DwgSectionLocatorRecord
	Descriptors map[string]*DwgSectionDescriptor
}

func (header *DwgFileHeader) Read(r io.ReadSeeker) (err error) {
	// Seek to start
	r.Seek(0, io.SeekStart)

	var value [6]byte
	if err = binary.Read(r, binary.LittleEndian, &value); err != nil {
		return err
	}

	var dwgVersion version.ACadVersion
	if dwgVersion, err = version.GetVersionFromName(string(value[:])); err != nil {
		return err
	}
	header.Version = dwgVersion

	var sreader stream.IDwgStreamReader
	if sreader, err = stream.NewDwgStreamHandler(dwgVersion, r); err != nil {
		return err
	}

	switch dwgVersion {
	case version.AC1024, version.AC1027, version.AC1032:
		err = header.readFileHeaderAC18(sreader)
	default:
		err = fmt.Errorf("Unsupported AutoCAD version: %d", dwgVersion)
	}

	return err
}

func (header *DwgFileHeader) readFileMetadata(sreader stream.IDwgStreamReader) error {
	err := error(nil)

	// 5 bytes of 0x00
	err = sreader.Advance(5)

	// 0x0B	1	Maintenance release version
	header.MetaData.AcadMaintenanceVersion, err = sreader.ReadByte()
	//0x0C	1	Byte 0x00, 0x01, or 0x03
	sreader.Advance(1)
	//0x0D	4	Preview address(long), points to the image page + page header size(0x20).
	header.MetaData.PreviewAddress, err = sreader.ReadRawLong()
	//0x11	1	Dwg version (Acad version that writes the file)
	header.MetaData.DwgVersion, err = sreader.ReadByte()
	//0x12	1	Application maintenance release version(Acad maintenance version that writes the file)
	header.MetaData.AppReleaseVersion, err = sreader.ReadByte()
	//0x13	2	Codepage
	codePage, _ := sreader.ReadShort()
	header.MetaData.DrawingCodePage = utils.GetCodePage(int(codePage))
	// this._encoding = sreader.Encoding = getListedEncoding((int)fileheader.DrawingCodePage);

	//Advance empty bytes
	//0x15	3	3 0x00 bytes
	sreader.Advance(3)

	//0x18	4	SecurityType (long), see R2004 meta data, the definition is the same, paragraph 4.1.
	header.MetaData.SecurityType, err = sreader.ReadRawLong()
	//0x1C	4	Unknown long
	sreader.ReadRawLong()
	//0x20	4	Summary info Address in stream
	header.MetaData.SummaryInfoAddr, err = sreader.ReadRawLong()
	//0x24	4	VBA Project Addr(0 if not present)
	header.MetaData.VbaProjectAddr, err = sreader.ReadRawLong()

	//0x28	4	0x00000080
	sreader.ReadRawLong()

	//0x2C	4	App info Address in stream
	sreader.ReadRawLong()

	//Get to offset 0x80
	//0x30	0x80	0x00 bytes
	sreader.Advance(80)

	return err
}

func (header *DwgFileHeader) readFileHeaderAC18(sreader stream.IDwgStreamReader) error {
	var err error

	if err = header.readFileMetadata(sreader); err != nil {
		return err
	}

	// # Read Encrypted HeaderData
	var encryptedData []byte
	if encryptedData, err = sreader.ReadBytes(0x6C); err != nil {
		return err
	}
	var headerStream io.Reader
	if headerStream, err = stream.GetCRC32StreamHandler(encryptedData, 0); err != nil {
		return err
	}
	// headerStream.Encoding = utils.GetListedEncoding(utils.CP_Windows1252)

	if _, err = sreader.ReadBytes(20); err != nil {
		return err
	}

	// Read header encrypted data
	err = binary.Read(headerStream, binary.LittleEndian, &header.HeaderData)
	if err != nil {
		return err
	} else if string(header.HeaderData.FileId[:]) != "AcFssFcAJMB\x00" {
		return errors.New("File validation failed, invalid id")
	}
	header.HeaderData.PageMapAddress += 0x100

	// # Read page map of the file
	if _, err = sreader.SetPosition(int64(header.HeaderData.PageMapAddress)); err != nil {
		return err
	}

	// Get the page size
	var pageHeaderData *PageHeaderData
	if pageHeaderData, err = header.getPageHeaderData(sreader); err != nil {
		return err
	}

	// Get the decompressed stream to read the records
	var decompressed *stream.MemoryStream
	if decompressed, err = stream.DecompressLZ77AC18(sreader.GetStream(), int(pageHeaderData.DecompressedSize)); err != nil {
		return err
	}

	// Section size
	total := int32(0x100)
	for decompressed.Position() < decompressed.Length() {
		record := DwgSectionLocatorRecord{}
		//0x00	4	Section page number, starts at 1, page numbers are unique per file.
		if err = binary.Read(decompressed, binary.LittleEndian, &record.Number); err != nil {
			return err
		}
		//0x04	4	Section size
		if err = binary.Read(decompressed, binary.LittleEndian, &record.Size); err != nil {
			return err
		}

		if record.Number >= 0 {
			record.Seeker = total
			header.Records[record.Number] = &record
		} else {
			//If the section number is negative, this represents a gap in the sections (unused data).
			//For a negative section number, the following data will be present after the section size:

			//0x00	4	Parent
			decompressed.Seek(4, io.SeekCurrent)
			//0x04	4	Left
			decompressed.Seek(4, io.SeekCurrent)
			//0x08	4	Right
			decompressed.Seek(4, io.SeekCurrent)
			//0x0C	4	0x00
			decompressed.Seek(4, io.SeekCurrent)
		}

		total += record.Size
	}

	// # Read the data section map

	// Set the position of the map
	sectionMapRecord := header.Records[int32(header.HeaderData.SectionMapId)]
	if _, err = sreader.SetPosition(int64(sectionMapRecord.Seeker)); err != nil {
		return err
	}

	// Get the page size
	if pageHeaderData, err = header.getPageHeaderData(sreader); err != nil {
		return err
	}
	if decompressed, err = stream.DecompressLZ77AC18(sreader.GetStream(), int(pageHeaderData.DecompressedSize)); err != nil {
		return err
	}

	//0x00	4	Number of section descriptions(NumDescriptions)
	var nDescriptions int32
	if err = binary.Read(decompressed, binary.LittleEndian, &nDescriptions); err != nil {
		return err
	}

	//0x04	4	0x02 (long)
	decompressed.Seek(4, io.SeekCurrent)
	//0x08	4	0x00007400 (long)
	decompressed.Seek(4, io.SeekCurrent)
	//0x0C	4	0x00 (long)
	decompressed.Seek(4, io.SeekCurrent)
	//0x10	4	Unknown (long), ODA writes NumDescriptions here.
	decompressed.Seek(4, io.SeekCurrent)

	var name [64]byte
	for i := int32(0); i < nDescriptions; i++ {
		descriptor := DwgSectionDescriptor{}
		//0x00	8	Size of section(OdUInt64)
		binary.Read(decompressed, binary.LittleEndian, &descriptor.CompressedSize)
		/*0x08	4	Page count(PageCount). Note that there can be more pages than PageCount,
		as PageCount is just the number of pages written to file.
		If a page contains zeroes only, that page is not written to file.
		These “zero pages” can be detected by checking if the page’s start
		offset is bigger than it should be based on the sum of previously read pages
		decompressed size(including zero pages).After reading all pages, if the total
		decompressed size of the pages is not equal to the section’s size, add more zero
		pages to the section until this condition is met.
		*/
		binary.Read(decompressed, binary.LittleEndian, &descriptor.PageCount)
		//0x0C	4	Max Decompressed Size of a section page of this type(normally 0x7400)
		binary.Read(decompressed, binary.LittleEndian, &descriptor.DecompressedSize)
		//0x10	4	Unknown(long)
		decompressed.Seek(4, io.SeekCurrent)
		//0x14	4	Compressed(1 = no, 2 = yes, normally 2)
		binary.Read(decompressed, binary.LittleEndian, &descriptor.CompressedCode)
		//0x18	4	Section Id(starts at 0). The first section(empty section) is numbered 0, consecutive sections are numbered descending from(the number of sections – 1) down to 1.
		binary.Read(decompressed, binary.LittleEndian, &descriptor.SectionId)
		//0x1C	4	Encrypted(0 = no, 1 = yes, 2 = unknown)
		binary.Read(decompressed, binary.LittleEndian, &descriptor.Encrypted)
		//0x20	64	Section Name(string)
		binary.Read(decompressed, binary.LittleEndian, &name)
		descriptor.Name = utils.BytesToString(name[:])

		//Following this, the following (local) section page map data will be present
		for j := int32(0); j < descriptor.PageCount; j++ {
			localmap := DwgLocalSectionMap{}
			//0x00	4	Page number(index into SectionPageMap), starts at 1
			binary.Read(decompressed, binary.LittleEndian, &localmap.PageNumber)
			//0x04	4	Data size for this page(compressed size).
			binary.Read(decompressed, binary.LittleEndian, &localmap.CompressedSize)
			//0x08	8	Start offset for this page(OdUInt64).If this start offset is smaller than the sum of the decompressed size of all previous pages, then this page is to be preceded by zero pages until this condition is met.
			binary.Read(decompressed, binary.LittleEndian, &localmap.Offset)

			//same decompressed size and seeker (temporal values)
			localmap.DecompressedSize = descriptor.DecompressedSize
			localmap.Seeker = header.Records[localmap.PageNumber].Seeker

			//Maximum section page size appears to be 0x7400 bytes in the normal case.
			//If a logical section of the file (the database objects, for example) exceeds this size, then it is broken up into pages of size 0x7400.

			descriptor.LocalSections = append(descriptor.LocalSections, &localmap)
		}

		//Get the final size for the local section
		sizeLeft := uint32(descriptor.CompressedSize % uint64(descriptor.DecompressedSize))
		if sizeLeft > 0 && len(descriptor.LocalSections) > 0 {
			descriptor.LocalSections[len(descriptor.LocalSections)-1].DecompressedSize = sizeLeft
		}

		header.Descriptors[descriptor.Name] = &descriptor
	}

	return nil
}

func (header *DwgFileHeader) getPageHeaderData(sreader stream.IDwgStreamReader) (headerData *PageHeaderData, err error) {
	pageHeaderData := PageHeaderData{}

	//0x00	4	Section page type:
	//Section page map: 0x41630e3b
	//Section map: 0x4163003b
	pageHeaderData.SectionType, err = sreader.ReadRawLong()
	//0x04	4	Decompressed size of the data that follows
	pageHeaderData.DecompressedSize, err = sreader.ReadRawLong()
	//0x08	4	Compressed size of the data that follows(CompDataSize)
	pageHeaderData.CompressedSize, err = sreader.ReadRawLong()

	//0x0C	4	Compression type(0x02)
	pageHeaderData.CompressionType, err = sreader.ReadRawLong()
	//0x10	4	Section page checksum
	var checksum int32
	if checksum, err = sreader.ReadRawLong(); err == nil {
		pageHeaderData.Checksum = uint32(checksum)
	}

	return &pageHeaderData, err
}

func NewDwgFileHeader() *DwgFileHeader {
	return &DwgFileHeader{
		Records:     make(map[int32]*DwgSectionLocatorRecord),
		Descriptors: make(map[string]*DwgSectionDescriptor),
	}
}
