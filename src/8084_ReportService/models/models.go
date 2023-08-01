package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Object struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id"`
	Name              string             `json:"name" bson:"name"`
	DisplayName       string             `json:"displayName" bson:"displayName"`
	Description       string             `json:"description" bson:"description"`
	ExpirationTime    int                `json:"expirationtime" bson:"expirationtime"`
	IsAutoDisappear   bool               `json:"isautodisappear" bson:"isautodisappear"`
	AutoDisappearMode int                `json:"autodisappearmode" bson:"autodisappearmode"`
	Color             string             `json:"color" bson:"color"`
	Icon              string             `json:"icon" bson:"icon"`
	ClearMode         int                `json:"clearmode" bson:"clearmode"`
	ClearInDay        int                `json:"clearinday" bson:"clearinday"`
	Filter            ObjectFilter       `json:"objectfilter" bson:"objectfilter"`
	Attribute         []ObjectAttribute  `json:"objectattribute" bson:"objectattribute"`
	Ancestor          ObjectAncestor     `json:"objectancestor" bson:"objectancestor"`
}

type ObjectFilter struct {
	IsTagOffset bool `json:"istagoffset" bson:"istagoffset"`
	OffsetX     int  `json:"offsetx" bson:"offsetx"`
	OffsetY     int  `json:"offsety" bson:"offsety"`
	OffsetZ     int  `json:"offsetz" bson:"offsetz"`

	IsZoneLockDown bool `json:"iszonelockdown" bson:"iszonelockdown"`
	ZoneClock      int  `json:"zoneclock" bson:"zoneclock"`

	IsMoveMedian bool `json:"ismovemedian" bson:"ismovemedian"`
	MoveRadius   int  `json:"moveradius" bson:"moveradius"`
	MoveCount    int  `json:"movecount" bson:"movecount"`
	MoveMedian   int  `json:"movemedian" bson:"movemedian"`
	MedianCount  int  `json:"mediancount" bson:"mediancount"`
	MedianCountZ int  `json:"mediancountz" bson:"mediancountz"`

	IsRate      bool `json:"israte" bson:"israte"`
	RateRadius  int  `json:"rateradius" bson:"rateradius"`
	RateMaxTime int  `json:"ratemaxtime" bson:"ratemaxtime"`
	RateMinTime int  `json:"ratemintime" bson:"ratemintime"`

	IsZoneChange  bool `json:"iszonechange" bson:"iszonechange"`
	ZoneNumBlinks int  `json:"zonenumblinks" bson:"zonenumblinks"`
}

type ObjectAttribute struct {
	Attributes []ObjectAttributeSet `json:"attributes" bson:"attributes"`
}

type ObjectAttributeSet struct {
	Name         string   `json:"name" bson:"name"`
	DisplayName  string   `json:"displayname" bson:"displayname"`
	Description  string   `json:"description" bson:"description"`
	DataType     int      `json:"datatype" bson:"datatype"`
	IsSearchable bool     `json:"issearchable" bson:"issearchable"`
	IsLookuup    bool     `json:"islookuup" bson:"islookuup"`
	LookupText   []string `json:"lookuptext" bson:"lookuptext"`
}

type ObjectAncestor struct {
	ParentType      int `json:"parenttype" bson:"parenttype"`
	GroupUpdateMode int `json:"groupupdatemode" bson:"groupupdatemode"`
}

type ObjectReport struct {
	ObjectID        string    `json:"objectid" bson:"parenttype"`
	Icon            string    `json:"icon" bson:"icon"`
	ObjectType      string    `json:"objecttype" bson:"objecttype"`
	MapData         MapData   `json:"mapdata" bson:"mapdata"`
	TagID           string    `json:"tagid" bson:"tagid"`
	BlinkTimeStamp  time.Time `json:"blinktimestamp" bson:"blinktimestamp"`
	X_COR           float64   `json:"x_cor" bson:"x_cor"`
	Y_COR           float64   `json:"y_cor" bson:"y_cor"`
	Z_COR           float64   `json:"z_cor" bson:"z_cor"`
	Latitude        float64   `json:"latitude" bson:"latitude"`
	Longitude       float64   `json:"longitude" bson:"longitude"`
	ManualTimestamp time.Time `json:"manualtimestamp" bson:"manualtimestamp"`
	EPC_URI         string    `json:"epc_uri" bson:"epc_uri"`
	EPC_Company     string    `json:"epc_company" bson:"epc_company"`
	EPC_Reference   string    `json:"epc_reference" bson:"epc_reference"`
	EPC_Serial      string    `json:"epc_serial" bson:"epc_serial"`
}

type MapData struct {
	ZoneGroupName string `json:"zonegroupname" bson:"zonegroupname"`
	ZoneName      string `json:"zonename" bson:"zonename"`
	SiteName      string `json:"sitename" bson:"sitename"`
	MapName       string `json:"mapname" bson:"mapname"`
	SiteGroupName string `json:"sitegroupname" bson:"sitegroupname"`
}

type ObjectHistoryReport struct {
	X_COR      float64   `json:"x_cor" bson:"x_cor"`
	Y_COR      float64   `json:"y_cor" bson:"y_cor"`
	ReportTime time.Time `json:"reporttime" bson:"reporttime"`
}

type TagReport struct {
	TagID                 string `json:"tagid" bson:"tagid"`
	IsRegistered          bool   `json:"isregistered" bson:"isregistered"`
	IsLocated             bool   `json:"islocated" bson:"islocated"`
	ObjectType            string `json:"objecttype" bson:"objecttype"`
	TechnologySource      string `json:"technologysource" bson:"technologysource"`
	ObjectID              string `json:"objectid" bson:"objectid"`
	DeviceID              string `json:"deviceid" bson:"deviceid"`
	Protocol              string `json:"protocol" bson:"protocol"`
	IsBlinking            bool   `json:"isblinking" bson:"isblinking"`
	IsAlerting            bool   `json:"isalerting" bson:"isalerting"`
	LowBatteryTime        int    `json:"lowbatterytime" bson:"lowbatterytime"`
	LowBatteryElapsedTime int    `json:"lowbatteryelapsedtime" bson:"lowbatteryelapsed"`
}

type PrfidResponse struct {
	Message string             `json:"message" bson"message"`
	Data    []BriefPrfidDevice `json:"data", bson:"data"`
}

type BriefPrfidDevice struct {
	MAC         string      `json:"mac" bson:"mac"`
	AntennaData AntennaData `json:"antennadata" bson:"antennadata"`
	MapID       string      `json:"mapid" bson:"mapid"`
}

type Antenna struct {
	No         int   `json:"no" bson:"no"`
	IsSet      bool  `json:"isset" bson:"isset"`
	IsRunning  bool  `json:"isrunning" bson:"isrunning"`
	Location   Point `json:"location" bson:"location"`
	PowerLevel int   `json:"power_level" bson:"power_level"`
}

type AntennaData struct {
	Count    int       `json:"count,omitempty" bson:"count,omitempty"`
	Antennas []Antenna `json:"antennas,omitempty" bson:"antennas,omitempty"`
}

type Point struct {
	X float64 `json:"x,omitempty" bson:"x,omitempty"`
	Y float64 `json:"y,omitempty" bson:"y,omitempty"`
}
