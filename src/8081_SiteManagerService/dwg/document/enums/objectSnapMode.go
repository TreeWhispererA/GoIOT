package enums

type ObjectSnapMode uint16

const (
	OSM_None                 ObjectSnapMode = 0
	OSM_EndPoint             ObjectSnapMode = 1
	OSM_MidPoint             ObjectSnapMode = 2
	OSM_Center               ObjectSnapMode = 4
	OSM_Node                 ObjectSnapMode = 8
	OSM_Quadrant             ObjectSnapMode = 0x0010
	OSM_Intersection         ObjectSnapMode = 0x0020
	OSM_Insertion            ObjectSnapMode = 0x0040
	OSM_Perpendicular        ObjectSnapMode = 0x0080
	OSM_Tangent              ObjectSnapMode = 0x0100
	OSM_Nearest              ObjectSnapMode = 0x0200
	OSM_ClearsAllObjectSnaps ObjectSnapMode = 0x0400
	OSM_ApparentIntersection ObjectSnapMode = 0x0800
	OSM_Extension            ObjectSnapMode = 0x1000
	OSM_Parallel             ObjectSnapMode = 0x2000
	OSM_AllModes             ObjectSnapMode = 0x3FFF
	OSM_SwitchedOff          ObjectSnapMode = 0x4000
)
