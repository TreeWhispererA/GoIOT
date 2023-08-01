package builder

import (
	"tracio.com/sitemanagerservice/dwg/document"
	"tracio.com/sitemanagerservice/dwg/file/header"
)

type DwgDocumentBuilder struct {
	documentToBuilder *document.CadDocument
	// cadObjects map[uint32]*CadObject
	// templates map[uint32]*CadTemplate
	// tableTemplates map[uint32]*CadTableTemplate
	HeaderHandles *header.DwgHeaderHandlesCollection
}

func NewDwgDocumentBuilder(document *document.CadDocument) *DwgDocumentBuilder {
	return &DwgDocumentBuilder{
		documentToBuilder: document,
	}
}
