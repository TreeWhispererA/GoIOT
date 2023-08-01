package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Devices struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Type  int                `json:"type" bson:"type" `
	PRFID Device_prfid       `json:"prfid,omitempty" bson:"prfid,omitempty"`
	BLE   Device_ble         `json:"ble,omitempty" bson:"ble,omitempty"`
	// Appliance Device_appliance   `json:"appliance,omitempty" bson:"appliance,omitempty"`
}

type Template struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Type  int                `json:"type" bson:"type"`
	PRFID Device_prfid       `json:"prfid,omitempty" bson:"prfid,omitempty"`
	BLE   Device_ble         `json:"ble,omitempty" bson:"ble,omitempty"`
	// Appliance Device_appliance   `json:"appliance,omitempty" bson:"appliance,omitempty"`
}

type Device_prfid struct {
	IsActive          bool        `json:"isactive,omitempty" bson:"isactive,omitempty"`
	Name              string      `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Description       string      `json:"description,omitempty" bson:"description,omitempty"`
	MapID             string      `json:"mapid,omitempty" bson:"mapid,omitempty"`
	MapName           string      `json:"mapname,omitempty" bson:"mapname,omitempty"`
	IsLocated         bool        `json:"islocated" bson:"islocated"`
	Hostname          string      `json:"hostname,omitempty" bson:"hostname,omitempty" validate:"required"`
	IP                string      `json:"ip,omitempty" bson:"ip,omitempty"`
	Port              int         `json:"port,omitempty" bson:"port,omitempty"`
	TLS_ENABLED       bool        `json:"tls_enabled,omitempty" bson:"tls_enabled,omitempty"`
	MAC               string      `json:"mac,omitempty" bson:"mac,omitempty" validate:"required"`
	Manufacturer      string      `json:"manufacturer,omitempty" bson:"manufacturer,omitempty"`
	Model             string      `json:"model,omitempty" bson:"model,omitempty"`
	ReportInterval    int         `json:"report_interval,omitempty" bson:"report_interval,omitempty" validate:"required"`
	Appliance_UUID    string      `json:"appliance_uuid,omitempty" bson:"appliance_uuid,omitempty"`
	AntennaData       AntennaData `json:"antennadata,omitempty" bson:"antennadata,omitempty"`
	Template          string      `json:"template,omitempty" bson:"template,omitempty"`
	Model_Name        string      `json:"model_name,omitempty" bson:"model_name,omitempty"`
	Manufacturer_Name string      `json:"manufacturer_name,omitempty" bson:"manufacturer_name,omitempty"`
	OperationalMode   OpMode      `json:"operational_model,omitempty" bson:"operational_model,omitempty"`
}

type BriefPrfidDevice struct {
	MAC         string      `json:"mac" bson:"mac"`
	AntennaData AntennaData `json:"antennadata" bson:"antennadata"`
	MapID       string      `json:"mapid" bson:"mapid"`
}

type OpMode struct {
	Mode                         string `json:"mode,omitempty" bson:"mode,omitempty"`
	TagIDFilter                  string `json:"tag_id_filter,omitempty" bson:"tag_id_filter,omitempty"`
	FilterMatch                  string `json:"filter_match,omitempty" bson:"filter_match,omitempty"`
	FilterOperation              string `json:"filter_operation,omitempty" bson:"filter_operation,omitempty"`
	ReportingInterval            int    `json:"reporting_interval,omitempty" bson:"reporting_interval,omitempty"`
	Units                        string `json:"units,omitempty" bson:"units,omitempty"`
	MinimumRSSI                  int    `json:"minimum_rssi,omitempty" bson:"minimum_rssi,omitempty"`
	GPIPort                      int    `json:"gpi_port,omitempty" bson:"gpi_port,omitempty"`
	Signal                       string `json:"signal,omitempty" bson:"signal,omitempty"`
	Interval                     int    `json:"interval,omitempty" bson:"interval,omitempty"`
	UserDefinedJson              string `json:"user_defined_json,omitempty" bson:"user_defined_json,omitempty"`
	IsEnabledReaderConfiguration bool   `json:"is_enabled_reader_configuration,omitempty" bson:"is_enabled_reader_configuration,omitempty"`
	ReaderConfiguration          string `json:"reader_configuration,omitempty" bson:"reader_configuration,omitempty"`
}

type AntennaData struct {
	Count    int       `json:"count,omitempty" bson:"count,omitempty"`
	Antennas []Antenna `json:"antennas,omitempty" bson:"antennas,omitempty"`
}

type Device_ble struct {
	IsActive          bool   `json:"isactive,omitempty" bson:"isactive,omitempty" `
	Name              string `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Description       string `json:"description,omitempty" bson:"description,omitempty"`
	IsLocated         bool   `json:"islocated" bson:"islocated"`
	MapID             string `json:"mapid,omitempty" bson:"mapid,omitempty"`
	Hostname          string `json:"hostname,omitempty" bson:"hostname,omitempty" validate:"required"`
	MAC               string `json:"mac,omitempty" bson:"mac,omitempty" validate:"required"`
	Manufacturer      string `json:"manufacturer,omitempty" bson:"manufacturer,omitempty" validate:"required"`
	Model             string `json:"model,omitempty" bson:"model,omitempty"`
	Appliance_UUID    string `json:"appliance_uuid,omitempty" bson:"appliance_uuid,omitempty"`
	Location          Point  `json:"location,omitempty" bson:"location,omitempty"`
	Template          string `json:"template,omitempty" bson:"template,omitempty"`
	Model_Name        string `json:"model_name,omitempty" bson:"model_name,omitempty"`
	Manufacturer_Name string `json:"manufacturer_name,omitempty" bson:"manufacturer_name,omitempty"`
}

type Device_appliance struct {
	Name          string   `json:"name" bson:"name"`
	UUID          string   `json:"uuid" bson:"uuid"`
	Hostname      string   `json:"hostname" bson:"hostname"`
	IP            string   `json:"ip" bson:"ip"`
	Services      []string `json:"services" bson:"services"`
	Configuration []Config `json:"configuration" bson:"configuration"`
}

type DeviceTemplate struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	Type          int                `json:"type" bson:"type"`
	Configuration []Config           `json:"configuration" bson:"configuration"`
}

type Config struct {
	Name   string `json:"name" bson:"name"`
	Config string `json:"config" bson:"config"`
}

type Point struct {
	X float64 `json:"x,omitempty" bson:"x,omitempty"`
	Y float64 `json:"y,omitempty" bson:"y,omitempty"`
}

type Antenna struct {
	No         int   `json:"no" bson:"no"`
	IsSet      bool  `json:"isset" bson:"isset"`
	IsRunning  bool  `json:"isrunning" bson:"isrunning"`
	Location   Point `json:"location" bson:"location"`
	PowerLevel int   `json:"power_level" bson:"power_level"`
}

type DeviceType struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id" `
	Name        string             `bson:"name" json:"name"`
	Type        int                `bson:"type" json:"type"`
	Description string             `bson:"description" json:"description"`
}

type ObjectType struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id" `
	Name           string             `bson:"name" json:"name"`
	DisplayName    string             `bson:"display_name" json:"display_name"`
	Description    string             `bson:"description" json:"description"`
	ExpirationTime int                `bson:"expiration_time" json:"expiration_time"`
	Color          string             `bson:"color" json:"color"`
	Icon           string             `bson:"icon" json:"icon"`
}

type Manufacturer struct {
	ID   primitive.ObjectID `bson:"_id" json:"_id" `
	Type int                `bson:"type" json:"type"`
	Name string             `bson:"name" json:"name"`
}

type Model struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id" `
	Name         string             `bson:"name" json:"name"`
	Manufacturer string             `bson:"manufacturer" json:"manufacturer"`
}

type SelectData struct {
	Operation  string   `bson:"operation" json:"operation"`
	SelectedID []string `bson:"selectedid" json:"selectedid"`
}
