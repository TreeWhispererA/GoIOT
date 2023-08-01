package document

import (
	"tracio.com/sitemanagerservice/dwg/classes"
)

type CadDocument struct {
	Header      *CadHeader
	SummaryInfo *CadSummaryInfo
	Classes     map[int16]*classes.DxfClass
}

func NewCadDocument() *CadDocument {
	return &CadDocument{Classes: make(map[int16]*classes.DxfClass)}
}
