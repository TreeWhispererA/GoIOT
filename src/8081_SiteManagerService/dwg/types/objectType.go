package types

type ObjectType int16

const (
	Object_UNLISTED             = -999
	Object_INVALID              = -1
	Object_UNUSED               = 0
	Object_TEXT                 = 1
	Object_ATTRIB               = 2
	Object_ATTDEF               = 3
	Object_BLOCK                = 4
	Object_ENDBLK               = 5
	Object_SEQEND               = 6
	Object_INSERT               = 7
	Object_MINSERT              = 8
	Object_UNKNOW_9             = 9
	Object_VERTEX_2D            = 0x0A
	Object_VERTEX_3D            = 0x0B
	Object_VERTEX_MESH          = 0x0C
	Object_VERTEX_PFACE         = 0x0D
	Object_VERTEX_PFACE_FACE    = 0x0E
	Object_POLYLINE_2D          = 0x0F
	Object_POLYLINE_3D          = 0x10
	Object_ARC                  = 0x11
	Object_CIRCLE               = 0x12
	Object_LINE                 = 0x13
	Object_DIMENSION_ORDINATE   = 0x14
	Object_DIMENSION_LINEAR     = 0x15
	Object_DIMENSION_ALIGNED    = 0x16
	Object_DIMENSION_ANG_3_Pt   = 0x17
	Object_DIMENSION_ANG_2_Ln   = 0x18
	Object_DIMENSION_RADIUS     = 0x19
	Object_DIMENSION_DIAMETER   = 0x1A
	Object_POINT                = 0x1B
	Object_FACE3D               = 0x1C
	Object_POLYLINE_PFACE       = 0x1D
	Object_POLYLINE_MESH        = 0x1E
	Object_SOLID                = 0x1F
	Object_TRACE                = 0x20
	Object_SHAPE                = 0x21
	Object_VIEWPORT             = 0x22
	Object_ELLIPSE              = 0x23
	Object_SPLINE               = 0x24
	Object_REGION               = 0x25
	Object_SOLID3D              = 0x26
	Object_BODY                 = 0x27
	Object_RAY                  = 0x28
	Object_XLINE                = 0x29
	Object_DICTIONARY           = 0x2A
	Object_OLEFRAME             = 0x2B
	Object_MTEXT                = 0x2C
	Object_LEADER               = 0x2D
	Object_TOLERANCE            = 0x2E
	Object_MLINE                = 0x2F
	Object_BLOCK_CONTROL_OBJ    = 0x30
	Object_BLOCK_HEADER         = 0x31
	Object_LAYER_CONTROL_OBJ    = 0x32
	Object_LAYER                = 0x33
	Object_STYLE_CONTROL_OBJ    = 0x34
	Object_STYLE                = 0x35
	Object_UNKNOW_36            = 0x36
	Object_UNKNOW_37            = 0x37
	Object_LTYPE_CONTROL_OBJ    = 0x38
	Object_LTYPE                = 0x39
	Object_UNKNOW_3A            = 0x3A
	Object_UNKNOW_3B            = 0x3B
	Object_VIEW_CONTROL_OBJ     = 0x3C
	Object_VIEW                 = 0x3D
	Object_UCS_CONTROL_OBJ      = 0x3E
	Object_UCS                  = 0x3F
	Object_VPORT_CONTROL_OBJ    = 0x40
	Object_VPORT                = 0x41
	Object_APPID_CONTROL_OBJ    = 0x42
	Object_APPID                = 0x43
	Object_DIMSTYLE_CONTROL_OBJ = 0x44
	Object_DIMSTYLE             = 0x45
	Object_VP_ENT_HDR_CTRL_OBJ  = 0x46
	Object_VP_ENT_HDR           = 0x47
	Object_GROUP                = 0x48
	Object_MLINESTYLE           = 0x49
	Object_OLE2FRAME            = 0x4A
	Object_DUMMY                = 0x4B
	Object_LONG_TRANSACTION     = 0x4C
	Object_LWPOLYLINE           = 0x4D
	Object_HATCH                = 0x4E
	Object_XRECORD              = 0x4F
	Object_ACDBPLACEHOLDER      = 0x50
	Object_VBA_PROJECT          = 0x51
	Object_LAYOUT               = 0x52
	Object_ACAD_PROXY_ENTITY    = 0x1f2
	Object_ACAD_PROXY_OBJECT    = 0x1f3
)
