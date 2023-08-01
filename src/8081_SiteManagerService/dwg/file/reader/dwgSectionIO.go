package reader

import (
	"fmt"

	"tracio.com/sitemanagerservice/dwg/file/reader/stream"
	"tracio.com/sitemanagerservice/dwg/file/section"
	"tracio.com/sitemanagerservice/dwg/file/version"
)

type DwgSectionIO struct {
	DwgVersion  version.ACadVersion
	SectionName section.DwgSectionDefinition
}

func CheckSentinel(actual []byte, expected []byte) bool {
	if len(actual) != len(expected) {
		return false
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			return false
		}
	}

	return true
}

func (this *DwgSectionIO) checkSentinel(sreader stream.IDwgStreamReader, expected []byte) (err error) {
	var sn []byte
	if sn, err = sreader.ReadSentinel(); err != nil {
		return err
	}

	if !CheckSentinel(sn, expected) {
		return fmt.Errorf("Invalid section sentinel found in %s", this.SectionName)
	}

	return nil
}

func (this *DwgSectionIO) IsR13_14Only() bool {
	return this.DwgVersion >= version.AC1012 && this.DwgVersion <= version.AC1014
}

func (this *DwgSectionIO) IsR13_15Only() bool {
	return this.DwgVersion >= version.AC1012 && this.DwgVersion <= version.AC1015
}

func (this *DwgSectionIO) IsR2004Pre() bool {
	return this.DwgVersion < version.AC1018
}

func (this *DwgSectionIO) IsR2000Plus() bool {
	return this.DwgVersion >= version.AC1015
}

func (this *DwgSectionIO) IsR2004Plus() bool {
	return this.DwgVersion >= version.AC1018
}

func (this *DwgSectionIO) IsR2007Plus() bool {
	return this.DwgVersion >= version.AC1021
}

func (this *DwgSectionIO) IsR2010Plus() bool {
	return this.DwgVersion >= version.AC1024
}

func (this *DwgSectionIO) IsR2013Plus() bool {
	return this.DwgVersion >= version.AC1027
}

func (this *DwgSectionIO) IsR2018Plus() bool {
	return this.DwgVersion >= version.AC1032
}
