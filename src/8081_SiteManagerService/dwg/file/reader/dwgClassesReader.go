package reader

import (
	"io"

	"tracio.com/sitemanagerservice/dwg/classes"
	"tracio.com/sitemanagerservice/dwg/file/reader/stream"
	"tracio.com/sitemanagerservice/dwg/file/section"
	"tracio.com/sitemanagerservice/dwg/file/version"
)

type DwgClassesReader struct {
	DwgSectionIO

	maintenanceVersion byte
	sreader            stream.IDwgStreamReader
}

func (this *DwgClassesReader) getCurrPos(sreader stream.IDwgStreamReader) (int64, error) {
	if this.IsR2007Plus() {
		pos, err := this.sreader.GetPositionInBits()
		return int64(pos), err
	} else {
		return this.sreader.GetPosition()
	}
}

func (this *DwgClassesReader) Read() (map[int16]*classes.DxfClass, error) {
	classesMap := make(map[int16]*classes.DxfClass)
	err := error(nil)

	//SN : 0x8D 0xA1 0xC4 0xB8 0xC4 0xA9 0xF8 0xC5 0xC0 0xDC 0xF4 0x5F 0xE7 0xCF 0xB6 0x8A
	this.checkSentinel(this.sreader, section.StartSentienls[this.SectionName])

	//RL : size of class data area
	var position int64
	var size int32
	if size, err = this.sreader.ReadRawLong(); err != nil {
		return nil, err
	} else if position, err = this.sreader.GetPosition(); err != nil {
		return nil, err
	}
	endSection := position + int64(size)

	//R2010+ (only present if the maintenance version is greater than 3!)
	if (this.DwgVersion >= version.AC1024 && this.maintenanceVersion > 3) ||
		this.DwgVersion > version.AC1027 {
		//RL : unknown, possibly the high 32 bits of a 64-bit size?
		this.sreader.ReadRawLong()
	}

	var flagPos int32
	//+R2007 Only:
	if this.IsR2007Plus() {
		var offset, savedOffset int32
		// Setup readers
		flagPos, err = this.sreader.GetPositionInBits()
		offset, err = this.sreader.ReadRawLong()
		flagPos += offset - 1

		savedOffset, err = this.sreader.GetPositionInBits()
		offset, err = this.sreader.SetPositionByFlag(int32(flagPos))
		endSection = int64(offset)

		this.sreader.SetPositionInBits(savedOffset)

		//Setup the text reader for versions 2007 and above
		var textReader stream.IDwgStreamReader
		//Create a copy of the stream
		var newStream io.ReadSeeker
		newStream, err = stream.CloneStream(this.sreader.GetStream())
		textReader, err = stream.NewDwgStreamHandler(this.DwgVersion, newStream)

		//Set the position and use the flag
		textReader.SetPositionInBits(int32(endSection))

		this.sreader = stream.NewDwgMergedReader(this.sreader, textReader, nil)

		//BL: 0x00
		this.sreader.ReadBitLong()
		//B : flag - to find the data string at the end of the section
		this.sreader.ReadBit()
	}

	if this.DwgVersion == version.AC1018 {
		//BS : Maxiumum class number
		this.sreader.ReadBitShort()
		//RC: 0x00
		this.sreader.ReadRawChar()
		//RC: 0x00
		this.sreader.ReadRawChar()
		//B : true
		this.sreader.ReadBit()
	}

	//We read sets of these until we exhaust the data.
	for true {
		if position, err = this.getCurrPos(this.sreader); err != nil {
			return nil, err
		} else if position >= endSection {
			break
		}

		dxfClass := classes.NewDxfClass()
		//BS : classnum
		dxfClass.ClassNumber, err = this.sreader.ReadBitShort()
		//BS : version – in R14, becomes a flag indicating whether objects can be moved, edited, etc.
		var proxyFlags int16
		if proxyFlags, err = this.sreader.ReadBitShort(); err == nil {
			dxfClass.ProxyFlags = classes.ProxyFlags(proxyFlags)
		}

		//TV : appname
		dxfClass.ApplicationName, err = this.sreader.ReadVariableText()
		//TV: cplusplusclassname
		dxfClass.CppClassName, err = this.sreader.ReadVariableText()
		//TV : classdxfname
		dxfClass.DxfName, err = this.sreader.ReadVariableText()

		//B : wasazombie
		dxfClass.WasZombie, err = this.sreader.ReadBit()
		//BS : itemclassid -- 0x1F2 for classes which produce entities, 0x1F3 for classes which produce objects.
		dxfClass.ItemClassId, err = this.sreader.ReadBitShort()

		if this.DwgVersion == version.AC1018 {
			//BL : Number of objects created of this type in the current DB(DXF 91).
			this.sreader.ReadBitLong()
			//BS : Dwg Version
			this.sreader.ReadBitShort()
			//BS : Maintenance release version.
			this.sreader.ReadBitShort()
			//BL : Unknown(normally 0L)
			this.sreader.ReadBitLong()
			//BL : Unknown(normally 0L)
			this.sreader.ReadBitLong()
		} else if this.DwgVersion > version.AC1018 {
			//BL : Number of objects created of this type in the current DB(DXF 91).
			dxfClass.InstanceCount, err = this.sreader.ReadBitLong()

			//BS : Dwg Version
			this.sreader.ReadBitLong()
			//BS : Maintenance release version.
			this.sreader.ReadBitLong()
			//BL : Unknown(normally 0L)
			this.sreader.ReadBitLong()
			//BL : Unknown(normally 0L)
			this.sreader.ReadBitLong()
		}

		classesMap[dxfClass.ClassNumber] = dxfClass
	}

	if this.IsR2007Plus() {
		this.sreader.SetPositionInBits(flagPos + 1)
	}

	//RS: CRC
	this.sreader.ResetShift()

	//0x72,0x5E,0x3B,0x47,0x3B,0x56,0x07,0x3A,0x3F,0x23,0x0B,0xA0,0x18,0x30,0x49,0x75
	this.checkSentinel(this.sreader, section.EndSentinels[this.SectionName])

	return classesMap, err
}

func NewDwgClassesReader(dwgVersion version.ACadVersion, maintenanceVersion byte,
	sreader stream.IDwgStreamReader) *DwgClassesReader {
	return &DwgClassesReader{
		DwgSectionIO: DwgSectionIO{
			DwgVersion:  dwgVersion,
			SectionName: section.Classes,
		},
		maintenanceVersion: maintenanceVersion,
		sreader:            sreader,
	}
}
