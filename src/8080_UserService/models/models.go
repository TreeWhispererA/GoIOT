package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Username     string             `json:"username" bson:"username" `
	Password     string             `json:"password" bson:"password"`
	Email        string             `json:"email" bson:"email"`
	DisplayMode  string             `json:"display_mode" bson:"display_mode"`
	Is2FAEnabled bool               `json:"is_2fa_enabled" bson:"is_2fa_enabled"`
	Language     string             `json:"language" bson:"language"`
	Phone        string             `json:"phone" bson:"phone"`
	Firstname    string             `json:"firstname" bson:"firstname" validate:"required,alpha"`
	Lastname     string             `json:"lastname" bson:"lastname" validate:"required,alpha"`
	Description  string             `json:"description" bson:"description"`
	LastConnect  time.Time          `json:"lastconnect" bson:"lastconnect"`
	Role         []string           `json:"role" bson:"role"`
	BriefRole    []BriefRole        `json:"briefrole" bson:"briefrole"`
	Photo        string             `json:"photo" bson:"photo"`
	OtpSecret    string             `json:"otpsecret" bson:"otpsecret"`
	OtpAuthUrl   string             `json:"otpauthurl" bson:"otpauthurl"`
}

type PasswordChangeInfo struct {
	ID          string `json:"id" bson:"id"`
	OldPassword string `json:"old_password" bson:"old_password"`
	NewPassword string `json:"new_password" bson:"new_password"`
}

type BriefRole struct {
	ID    string `json:"_id" bson:"_id"`
	Color string `json:"color" bson:"color"`
	Name  string `json:"name" bson:"name"`
}

type Role struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
	Name             string             `json:"name" bson:"name" `
	Color            string             `json:"color" bson:"color" `
	Description      string             `json:"description" bson:"description"`
	SiteGroupID      string             `json:"sitegroupid" bson:"sitegroupid"`
	IsAdmin          int                `json:"isadmin" bson:"isadmin"`
	SitePermission   []Permission       `json:"sitepermission" bson:"sitepermission"`
	PagePermission   PagePermission     `json:"pagepermission" bson:"pagepermission"`
	DevicePermission []Permission       `json:"devicepermission" bson:"devicepermission"`
	Users            int                `json:"users" bson:"users" `
	CurrentUsers     []string           `json:"currentusers" bson:"currentusers"`
}

type PagePermission struct {
	Reports       ReportsPermission       `json:"reports" bson:"reports" `
	LocalLive     LocalLivePermission     `json:"locallive" bson:"locallive" `
	History       HistoryPermission       `json:"history" bson:"history" `
	Alert         AlertPermission         `json:"alert" bson:"alert" `
	SiteManager   SiteManagerPermission   `json:"sitemanager" bson:"sitemanager" `
	DeviceManager DeviceManagerPermission `json:"devicemanager" bson:"devicemanager" `
	RulesManager  RulesManagerPermission  `json:"rulesmanager" bson:"rulesmanager" `
	Analytics     AnalyticsPermission     `json:"analytics" bson:"analytics" `
	User          UserPermission          `json:"user" bson:"user" `
	System        SystemPermission        `json:"system" bson:"system" `
}

type ArrayPermission struct {
	Access          bool         `json:"access" bson:"access" `
	PermissionGroup []Permission `json:"data" bson:"data" `
}

type HistoryArrayPermission struct {
	Access  bool         `json:"access" bson:"access" `
	Sites   []Permission `json:"sites" bson:"sites" `
	Objects []Permission `json:"objects" bson:"objects" `
}

type AccessOnlyPermission struct {
	Access bool `json:"access" bson:"access" `
}

type ReportArrayPermission struct {
	Access     bool         `json:"access" bson:"access" `
	Object     []Permission `json:"object" bson:"object" `
	Technology []Permission `json:"technology" bson:"technology" `
	Sites      []Permission `json:"sites" bson:"sites" `
}

type ReportsPermission struct {
	Access       bool                  `json:"access" bson:"access" `
	ObjectReport ReportArrayPermission `json:"objectreport" bson:"objectreport" `
	TagReport    ReportArrayPermission `json:"tagreport" bson:"tagreport" `
}

type LocalLivePermission struct {
	Access bool `json:"access" bson:"access" `
	Zones  int  `json:"zones" bson:"zones" `
}

type HistoryPermission struct {
	Access      bool                   `json:"access" bson:"access" `
	ObjectEvent HistoryArrayPermission `json:"objectevent" bson:"objectevent" `
	TagBlink    HistoryArrayPermission `json:"tagblink" bson:"tagblink" `
}

type AlertPermission struct {
	Access      bool            `json:"access" bson:"access" `
	ObjectAlert ArrayPermission `json:"objectalert" bson:"objectalert" `
	SystemAlert ArrayPermission `json:"systemalert" bson:"systemalert" `
}

type SiteManagerPermission struct {
	Access    bool         `json:"access" bson:"access" `
	Sitegroup []Permission `json:"sitegroup" bson:"sitegroup" `
}

type DeviceManagerPermission struct {
	Access  bool         `json:"access" bson:"access" `
	Devices []Permission `json:"devices" bson:"devices" `
}

type RulesManagerPermission struct {
	Access bool `json:"access" bson:"access" `
	Rules  int  `json:"rules" bson:"rules" `
}

type AnalyticsPermission struct {
	Access    bool `json:"access" bson:"access" `
	Analytics int  `json:"analytics" bson:"analytics" `
}

type UserPermission struct {
	Access bool `json:"access" bson:"access" `
	Add    bool `json:"add" bson:"add" `
	Edit   bool `json:"edit" bson:"edit" `
	Delete bool `json:"delete" bson:"delete" `
	View   bool `json:"view" bson:"view" `
}

type SystemPermission struct {
	Access           bool `json:"access" bson:"access" `
	SystemPermission int  `json:"systempermission" bson:"systempermission" `
}

type Permission struct {
	PermissionID    string `json:"pname" bson:"pname"`
	PermissionValue int    `json:"pvalue" bson:"pvalue"`
}

type ObjectType struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type SmartRule struct {
	PerPage   int        `json:"perpage" bson:"perpage"`
	PageNum   int        `json:"pagenum" bson:"pagenum"`
	Sorting   bool       `json:"sorting" bson:"sorting"`
	Filtering bool       `json:"filtering" bson:"filtering"`
	SortIndex string     `json:"sortfield" bson:"sortfield"`
	Ascending bool       `json:"ascending" bson:"ascending"`
	Filter    FilterRule `json:"filter" bson:"filter"`
}

type FilterRule struct {
	Reports       GroupFilter  `json:"reports" bson:"reports"`
	LocalLive     AccessFilter `json:"locallive" bson:"locallive"`
	History       GroupFilter  `json:"history" bson:"history"`
	Alert         GroupFilter  `json:"alert" bson:"alert"`
	SiteManager   AccessFilter `json:"sitemanager" bson:"sitemanager"`
	DeviceManager AccessFilter `json:"devicemanager" bson:"devicemanager"`
	RulesManager  AccessFilter `json:"rulesmanager" bson:"rulesmanager"`
	Analytics     AccessFilter `json:"analytics" bson:"analytics"`
	User          AccessFilter `json:"user" bson:"user"`
	System        AccessFilter `json:"system" bson:"system"`
}

type AccessFilter struct {
	Access int `json:"access" bson:"access" `
}

type GroupFilter struct {
	Access int `json:"access" bson:"access" `
	Group1 int `json:"group1" bson:"group1" `
	Group2 int `json:"group2" bson:"group2" `
}
