package classes

type DxfClass struct {
	DxfName         string
	CppClassName    string
	ApplicationName string
	ProxyFlags      ProxyFlags
	InstanceCount   int32
	WasZombie       bool
	IsAnEntity      bool
	ClassNumber     int16
	ItemClassId     int16
}

func NewDxfClass() *DxfClass {
	return &DxfClass{
		ApplicationName: "ObjectDBX Classes",
	}
}
