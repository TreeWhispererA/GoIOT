package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Site struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id" `
	Name     string             `json:"name" bson:"name" `
	Level    int                `json:"level" bson:"level" `
	ParentID string             `json:"parentid" bson:"parentid" `
	ZoneData ZoneData           `json:"zonedata" bson:"zonedata"`
}

type ZoneData struct {
	Points    []Point `json:"points" bson:"points"`
	Color     string  `json:"color" bson:"color"`
	MapWidth  float64 `json:"mapwidth" bson:"mapwidth"`
	MapHeight float64 `json:"mapheight" bson:"mapheight"`
	Width     float64 `json:"width" bson:"width"`
	Height    float64 `json:"height" bson:"height"`
	Scale     float64 `json:"scale" bson:"scale"` // pixel/m
	IsFT      bool    `json:"isft" bson:"isft"`
}

type Point struct {
	XPOS int `json:"xpos" bson:"xpos"`
	YPOS int `json:"ypos" bson:"ypos"`
}

type RefPoint struct {
	XPOS      int     `json:"xpos" bson:"xpos"`
	YPOS      int     `json:"ypos" bson:"ypos"`
	Latitude  float64 `json:"lati" bson:"lati"`
	Longitude float64 `json:"long" bson:"long"`
}

type SiteInfo struct {
	ID            string `json:"_id" bson:"_id" `
	Level         int    `json:"level" bson:"level" `
	SiteGroupID   string `json:"sitegroupid" bson:"sitegroupid" `
	SiteGroupName string `json:"sitegroupname" bson:"sitegroupname" `
	SiteID        string `json:"siteid" bson:"siteid" `
	SiteName      string `json:"sitename" bson:"sitename" `
	MapID         string `json:"mapid" bson:"mapid" `
	MapName       string `json:"mapname" bson:"mapname" `
	ZoneGroupID   string `json:"zonegroupid" bson:"zonegroupid" `
	ZoneGroupName string `json:"zonegroupname" bson:"zonegroupname" `
	ZoneID        string `json:"zoneid" bson:"zoneid" `
	ZoneName      string `json:"zonename" bson:"zonename" `
}
