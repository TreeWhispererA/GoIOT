package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"tracio.com/staticservice/db"
	middlewares "tracio.com/staticservice/handlers"
	"tracio.com/staticservice/models"
)

var client = db.Dbconnect()

func Test(c *gin.Context) {

	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}

var GetObjectTypes = gin.HandlerFunc(func(c *gin.Context) {

	var objecttypes []*models.ObjectType
	collection := client.Database("staticservice").Collection("objecttypes")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("ObjectTypes Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var objecttype *models.ObjectType
		err := cursor.Decode(&objecttype)
		if err != nil {
			color.Red("ObjectType Decode Failed in below api...")
		} else {
			objecttypes = append(objecttypes, objecttype)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetObjectTypes in below api...")
		return
	}

	middlewares.SuccessArrRespond(objecttypes, "ObjectType", c.Writer)
})

var GetDeviceTypes = gin.HandlerFunc(func(c *gin.Context) {

	var devicetypes []*models.DeviceType
	collection := client.Database("staticservice").Collection("devicetypes")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("DeviceTypes Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var devicetype *models.DeviceType
		err := cursor.Decode(&devicetype)
		if err != nil {
			color.Red("ObjectType Decode Failed in below api...")
		} else {
			devicetypes = append(devicetypes, devicetype)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetDeviceTypes in below api...")
		return
	}

	middlewares.SuccessArrRespond(devicetypes, "DeviceType", c.Writer)
})

func UploadObjectIcon(c *gin.Context) {
	file, err := c.FormFile("file")
	name := c.PostForm("name")
	description := c.PostForm("description")

	var objectIcon models.ObjectIcon

	objectIcon.Name = name
	objectIcon.Description = description
	// fileName := c.PostForm("file_name")
	if err != nil {
		color.Red("File Upload Failed in below api...")
		return
	}

	collection := client.Database("staticservice").Collection("objecticons")

	objectIcon.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(context.TODO(), objectIcon)
	if err != nil {
		color.Red("ObjectIcon creation Failed in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)

	id := strings.Replace(string(res), `"`, ``, 2)

	filename := fmt.Sprintf("uploaded/objecticons/%v", id)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		color.Red("File Save Error In Server in below api...")
		return
	}
	defer f.Close()

	if err := c.SaveUploadedFile(file, f.Name()); err != nil {
		middlewares.ErrorResponse("File Save Error", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse(`New Icon Uploaded at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
}
