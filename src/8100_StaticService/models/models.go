package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeviceType struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" `
	Name        string             `bson:"name" json:"name"`
	Type        int                `bson:"type" json:"type"`
	Description string             `bson:"description" json:"description"`
}

type ObjectType struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" `
	Name           string             `bson:"name" json:"name"`
	DisplayName    string             `bson:"display_name" json:"display_name"`
	Description    string             `bson:"description" json:"description"`
	ExpirationTime int                `bson:"expiration_time" json:"expiration_time"`
	Color          string             `bson:"color" json:"color"`
	Icon           string             `bson:"icon" json:"icon"`
}


type ObjectIcon struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" `
	Name           string             `bson:"name" json:"name"`
    Description    string             `bson:"description" json:"description"`
}
