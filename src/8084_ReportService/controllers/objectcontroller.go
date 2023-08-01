package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	http "github.com/helios/go-sdk/proxy-libs/helioshttp"

	"github.com/fatih/color"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"tracio.com/reportservice/db"
	middlewares "tracio.com/reportservice/handlers"
	"tracio.com/reportservice/models"
)

var client = db.Dbconnect()

func Test(c *gin.Context) {
	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}

var GetObjects = gin.HandlerFunc(func(c *gin.Context) {

	var objects []*models.Object
	collection := client.Database("reportservice").Collection("objects")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Objects Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var object *models.Object
		err := cursor.Decode(&object)
		if err != nil {
			color.Red("ObjectType Decode Failed in below api...")
		} else {
			objects = append(objects, object)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetObjects in below api...")
		return
	}

	middlewares.SuccessArrRespond(objects, "Object", c.Writer)
})

var AddNewObject = gin.HandlerFunc(func(c *gin.Context) {
	var object models.Object
	err := c.BindJSON(&object)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in below api...")
		return
	}

	object.ID = primitive.NewObjectID()
	collection := client.Database("reportservice").Collection("objects")

	result, err := collection.InsertOne(context.TODO(), object)
	if err != nil {
		color.Red("Device Create Failed in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessMessageResponse(`Inserted new device at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
})

var EditObject = gin.HandlerFunc(func(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var new_object models.Object
	err := c.BindJSON(&new_object)
	if err != nil {
		color.Red("Request Decoding Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	new_object.ID = id
	update := bson.D{{Key: "$set", Value: new_object}}
	collection := client.Database("reportservice").Collection("objects")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		color.Red("Object Update Failed with ID : %v...", id)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	if res.MatchedCount == 0 {
		color.Red("Object Not Exist with ID : %v...", id)
		middlewares.ErrorResponse("Object does not exist", c.Writer)
		return
	}
	middlewares.SuccessMessageResponse("Updated", c.Writer)
})

var DeleteObject = gin.HandlerFunc(func(c *gin.Context) {

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Invalid Object ID to Delete in below api...")
		middlewares.ErrorResponse("Invalid Object ID", c.Writer)
		return
	}

	collection := client.Database("reportservice").Collection("object")
	res, derr := collection.DeleteOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}})
	if derr != nil || res.DeletedCount == 0 {
		color.Red("Object Does Not Exist with ID:%v ...", id)
		middlewares.ErrorResponse("Object does not exist", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("Deleted", c.Writer)

})

var GetObjectByID = gin.HandlerFunc(func(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in below api...")
		middlewares.ErrorResponse("Invalid Device ID", c.Writer)
		return
	}

	var object models.Object
	collection := client.Database("reportservice").Collection("objects")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&object)
	if err != nil {
		color.Red("No Object Exist with ID: %v in below api...", id)
		middlewares.ErrorResponse("Object does not exist", c.Writer)
		return
	}

	middlewares.SuccessOneRespond(object, "Object", c.Writer)
})

func GetPrfidResponseFromAPI(url string, tokenString string) (models.PrfidResponse, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.PrfidResponse{}, errors.New("Invalid Request.")
	}
	req.Header.Set("Token", tokenString)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return models.PrfidResponse{}, errors.New("Server Error")
	}
	defer resp.Body.Close()

	var response models.PrfidResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return models.PrfidResponse{}, errors.New("Server Error")
	}

	return response, nil
}

var GetObjectReport = gin.HandlerFunc(func(c *gin.Context) {
	var tokenString string
	tokenString = c.Request.Header.Get("Token")
	var response models.PrfidResponse
	response = models.PrfidResponse{}
	// response = response.(models.PrfidResponse)
	response, err := GetPrfidResponseFromAPI("http://localhost/api/v1/devicemanagerservice/briefprfiddevices", tokenString)
	if err != nil {
		middlewares.ErrorResponse("Server Error", c.Writer)
		return
	}
	var ptrBriefPrfidDevices []*models.BriefPrfidDevice
	for i := range response.Data {
		ptrBriefPrfidDevices = append(ptrBriefPrfidDevices, &response.Data[i])
	}

	middlewares.SuccessArrRespond(ptrBriefPrfidDevices, "BriefPrfidDevice", c.Writer)
})

var GetTagReport = gin.HandlerFunc(func(c *gin.Context) {

})
