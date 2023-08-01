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

	"tracio.com/devicemanagerservice/db"
	middlewares "tracio.com/devicemanagerservice/handlers"
	"tracio.com/devicemanagerservice/models"
	"tracio.com/devicemanagerservice/validators"
)

var client = db.Dbconnect()

func Test(c *gin.Context) {

	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}

var GetDevices = gin.HandlerFunc(func(c *gin.Context) {

	var model_s map[string]models.Model
	model_s, errs := GetModelsMap()

	if errs != "" {
		color.Red("Devices Database Issue in below api...")
	}

	var manufacturers map[string]models.Manufacturer
	manufacturers, errs = GetManufacturersMap()

	if errs != "" {
		color.Red("Manufacturers Database Issue in below api...")
	}

	var devices []*models.Devices
	collection := client.Database("devicemanagerservice").Collection("devices")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Devices Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var device *models.Devices
		err := cursor.Decode(&device)
		if err != nil {
			color.Red("ObjectType Decode Failed in below api...")
		} else {
			if device.Type == 0 {
				device.PRFID.Manufacturer_Name = manufacturers[device.PRFID.Manufacturer].Name
				device.PRFID.Model_Name = model_s[device.PRFID.Model].Name
			} else if device.Type == 1 {
				device.BLE.Manufacturer_Name = manufacturers[device.BLE.Manufacturer].Name
				device.BLE.Model_Name = model_s[device.BLE.Model].Name
			}
			devices = append(devices, device)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetDevices in below api...")
		return
	}

	middlewares.SuccessArrRespond(devices, "Devices", c.Writer)
})

var GetBriefDevices = gin.HandlerFunc(func(c *gin.Context) {

	var briefdevices []*models.BriefPrfidDevice
	collection := client.Database("devicemanagerservice").Collection("devices")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Devices Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var device *models.Devices
		var briefdevice *models.BriefPrfidDevice
		err := cursor.Decode(&device)
		if err != nil {
			color.Red("ObjectType Decode Failed in below api...")
		} else {
			briefdevice = &models.BriefPrfidDevice{}
			if device.Type == 0 {
				briefdevice.AntennaData = device.PRFID.AntennaData
				briefdevice.MAC = device.PRFID.MAC
				briefdevice.MapID = device.PRFID.MapID
			}
			briefdevices = append(briefdevices, briefdevice)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetDevices in below api...")
		return
	}

	middlewares.SuccessArrRespond(briefdevices, "BriefDevice", c.Writer)
})

var GetTemplates = gin.HandlerFunc(func(c *gin.Context) {

	var model_s map[string]models.Model
	model_s, errs := GetModelsMap()

	if errs != "" {
		color.Red("Models Database Issue in below api...")
	}

	var manufacturers map[string]models.Manufacturer
	manufacturers, errs = GetManufacturersMap()

	if errs != "" {
		color.Red("Manufacturers Database Issue in below api...")
	}

	var templates []*models.Template
	collection := client.Database("devicemanagerservice").Collection("templates")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Templates Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var template *models.Template
		err := cursor.Decode(&template)
		if err != nil {
			color.Red("ObjectType Decode Failed in below api...")
		} else {
			if template.Type == 0 {
				template.PRFID.Manufacturer_Name = manufacturers[template.PRFID.Manufacturer].Name
				template.PRFID.Model_Name = model_s[template.PRFID.Model].Name
			} else if template.Type == 1 {
				template.BLE.Manufacturer_Name = manufacturers[template.BLE.Manufacturer].Name
				template.BLE.Model_Name = model_s[template.BLE.Model].Name
			}
			templates = append(templates, template)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetTemplates in below api...")
		return
	}

	middlewares.SuccessArrRespond(templates, "Template", c.Writer)
})

var GetDeviceTemplates = gin.HandlerFunc(func(c *gin.Context) {

	var device_templates []*models.DeviceTemplate
	collection := client.Database("devicemanagerservice").Collection("devicetemplates")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Device Templates Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var device_template *models.DeviceTemplate
		err := cursor.Decode(&device_template)
		if err != nil {
			color.Red("ObjectType Decode Failed in below api...")
		} else {
			device_templates = append(device_templates, device_template)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetDeviceTemplates in below api...")
		return
	}

	middlewares.SuccessArrRespond(device_templates, "DeviceTemplate", c.Writer)
})

var GetManufacturers = gin.HandlerFunc(func(c *gin.Context) {

	var manufacturers []*models.Manufacturer
	collection := client.Database("devicemanagerservice").Collection("manufacturers")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Manufacturers Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var manufacturer *models.Manufacturer
		err := cursor.Decode(&manufacturer)
		if err != nil {
			color.Red("GetManufacturers Decode Failed in below api...")
		} else {
			manufacturers = append(manufacturers, manufacturer)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetManufacturers in below api...")
		return
	}

	middlewares.SuccessArrRespond(manufacturers, "Manufacturer", c.Writer)
})

var GetModels = gin.HandlerFunc(func(c *gin.Context) {

	var model_s []*models.Model
	collection := client.Database("devicemanagerservice").Collection("models")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Models Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var model *models.Model
		err := cursor.Decode(&model)
		if err != nil {
			color.Red("GetModels Decode Failed in below api...")
		} else {
			model_s = append(model_s, model)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetModels in below api...")
		return
	}

	middlewares.SuccessArrRespond(model_s, "Model", c.Writer)
})

var GetModelsByManufacturers = gin.HandlerFunc(func(c *gin.Context) {

	id := c.Param("id")

	var model_s []*models.Model
	collection := client.Database("devicemanagerservice").Collection("models")
	cursor, err := collection.Find(context.TODO(), bson.M{"manufacturer": id})
	if err != nil {
		color.Red("Models Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var model *models.Model
		err := cursor.Decode(&model)
		if err != nil {
			color.Red("GetModels Decode Failed in below api...")
		} else {
			model_s = append(model_s, model)
		}
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetModels in below api...")
		return
	}

	middlewares.SuccessArrRespond(model_s, "Model", c.Writer)
})

var AddNewDevice = gin.HandlerFunc(func(c *gin.Context) {
	var device models.Devices
	err := c.BindJSON(&device)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in below api...")
		return
	}

	device.ID = primitive.NewObjectID()
	collection := client.Database("devicemanagerservice").Collection("devices")

	result, err := collection.InsertOne(context.TODO(), device)
	if err != nil {
		color.Red("Device Create Failed in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessMessageResponse(`Inserted new device at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
})

var EditDevice = gin.HandlerFunc(func(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var new_device models.Devices
	err := c.BindJSON(&new_device)
	if err != nil {
		color.Red("Request Decoding Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	new_device.ID = id
	update := bson.D{{Key: "$set", Value: new_device}}
	collection := client.Database("devicemanagerservice").Collection("devices")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		color.Red("Device Update Failed with ID : %v...", id)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	if res.MatchedCount == 0 {
		color.Red("Device Not Exist with ID : %v...", id)
		middlewares.ErrorResponse("Device does not exist", c.Writer)
		return
	}
	middlewares.SuccessMessageResponse("Updated", c.Writer)
})

var DeleteDevice = gin.HandlerFunc(func(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Invalid Device ID to Delete in below api...")
		middlewares.ErrorResponse("Invalid Device ID", c.Writer)
		return
	}

	collection := client.Database("devicemanagerservice").Collection("devices")
	res, derr := collection.DeleteOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}})
	if derr != nil || res.DeletedCount == 0 {
		color.Red("Device Does Not Exist with ID:%v ...", id)
		middlewares.ErrorResponse("Device does not exist", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("Deleted", c.Writer)
})

var GetDeviceByID = gin.HandlerFunc(func(c *gin.Context) {

	var model_s map[string]models.Model
	model_s, errs := GetModelsMap()

	if errs != "" {
		color.Red("Devices Database Issue in below api...")
	}

	var manufacturers map[string]models.Manufacturer
	manufacturers, errs = GetManufacturersMap()

	if errs != "" {
		color.Red("Manufacturers Database Issue in below api...")
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in below api...")
		middlewares.ErrorResponse("Invalid Device ID", c.Writer)
		return
	}

	var device models.Devices
	collection := client.Database("devicemanagerservice").Collection("devices")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&device)
	if err != nil {
		color.Red("No Device Exist with ID: %v in below api...", id)
		middlewares.ErrorResponse("Device does not exist", c.Writer)
		return
	}

	if device.Type == 0 {
		device.PRFID.Manufacturer_Name = manufacturers[device.PRFID.Manufacturer].Name
		device.PRFID.Model_Name = model_s[device.PRFID.Model].Name
	} else if device.Type == 1 {
		device.BLE.Manufacturer_Name = manufacturers[device.BLE.Manufacturer].Name
		device.BLE.Model_Name = model_s[device.BLE.Model].Name
	}

	middlewares.SuccessOneRespond(device, "Devices", c.Writer)
})

var AddNewDeviceTemplate = gin.HandlerFunc(func(c *gin.Context) {
	var devicetemplate models.DeviceTemplate
	err := c.BindJSON(&devicetemplate)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in AddNewDeviceTemplate()...")
		return
	}

	devicetemplate.ID = primitive.NewObjectID()
	collection := client.Database("devicemanagerservice").Collection("devicetemplates")

	result, err := collection.InsertOne(context.TODO(), devicetemplate)
	if err != nil {
		color.Red("Device Template Create Failed in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessMessageResponse(`Inserted new device-template at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
})

var EditDeviceTemplate = gin.HandlerFunc(func(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var new_device_template models.DeviceTemplate
	err := c.BindJSON(&new_device_template)
	if err != nil {
		color.Red("Request Decoding Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	new_device_template.ID = id
	update := bson.D{{Key: "$set", Value: new_device_template}}
	collection := client.Database("devicemanagerservice").Collection("devicetemplates")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		color.Red("Device Template Update Failed with ID : %v...", id)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	if res.MatchedCount == 0 {
		color.Red("Device Template Not Exist with ID : %v...", id)
		middlewares.ErrorResponse("Device Template does not exist", c.Writer)
		return
	}
	middlewares.SuccessMessageResponse("Updated", c.Writer)
})

var DeleteDeviceTemplate = gin.HandlerFunc(func(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Invalid Device Template ID to Delete in below api...")
		middlewares.ErrorResponse("Invalid Device ID", c.Writer)
		return
	}

	collection := client.Database("devicemanagerservice").Collection("devicetemplates")
	res, derr := collection.DeleteOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}})
	if derr != nil || res.DeletedCount == 0 {
		color.Red("Device Template Does Not Exist with ID:%v ...", id)
		middlewares.ErrorResponse("Device Template does not exist", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("Deleted", c.Writer)
})

var GetDeviceTemplateByID = gin.HandlerFunc(func(c *gin.Context) {

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in below api...")
		middlewares.ErrorResponse("Invalid Device Template ID", c.Writer)
		return
	}

	var device_template models.DeviceTemplate
	collection := client.Database("devicemanagerservice").Collection("devicetemplates")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&device_template)
	if err != nil {
		color.Red("No Device Template Exist with ID: %v in below api...", id)
		middlewares.ErrorResponse("Device Template does not exist", c.Writer)
		return
	}

	middlewares.SuccessOneRespond(device_template, "DeviceTemplate", c.Writer)
})

func ImportJsonDeviceData(c *gin.Context) {
	file, err := c.FormFile("file")
	devicetype := c.PostForm("type")
	// fileName := c.PostForm("file_name")
	if err != nil {
		color.Red("File Upload Failed in below api...")
		return
	}
	// timenow := time.Now().Format("2006.1.2_15:04:05")
	var filename string
	filename = fmt.Sprintf("uploaded/devicedata/json/%v", file.Filename)
	fmt.Println(filename)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		color.Red("File Save Error In Server in below api...")
		middlewares.ErrorResponse("Import File Failed...", c.Writer)
		return
	}
	defer f.Close()

	if err := c.SaveUploadedFile(file, f.Name()); err != nil {
		middlewares.ErrorResponse("File Save Error", c.Writer)
		return
	}

	file_open, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		middlewares.ErrorResponse("Import File Failed...", c.Writer)
		return
	}
	defer file_open.Close()

	decoder := json.NewDecoder(file_open)
	deviceList := []models.Devices{}
	err = decoder.Decode(&deviceList)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		middlewares.ErrorResponse("Error decoding JSON...", c.Writer)
		return
	}

	collection := client.Database("devicemanagerservice").Collection("devices")

	var result_str string
	var failed_str string

	for _, device := range deviceList {
		// Save device to database here
		if device.Type == 0 {
			if ok, errors := validators.ValidateInputs(device.PRFID); !ok {
				middlewares.ValidationResponse(errors, c.Writer)
				color.Red("Validation Error...")
				return
			}
			device.BLE = models.Device_ble{}
		}

		if device.Type == 1 {
			if ok, errors := validators.ValidateInputs(device.BLE); !ok {
				middlewares.ValidationResponse(errors, c.Writer)
				color.Red("Validation Error...")
				return
			}
			device.PRFID = models.Device_prfid{}
		}
		if devicetype == "prfid" && device.Type != 0 {
			middlewares.ErrorResponse("Please Input validate PRFID Device Data...", c.Writer)
			return
		}
		if devicetype == "ble" && device.Type != 1 {
			middlewares.ErrorResponse("Please Input validate PRFID Device Data...", c.Writer)
			return
		}

	}
	// Loop through the devices and save them to our database
	for index, device := range deviceList {
		// Save device to database here
		device.ID = primitive.NewObjectID()

		result, err := collection.InsertOne(context.TODO(), device)
		if err != nil {
			color.Red("Device Create Failed in below api...")
			failed_str = fmt.Sprintf("%v%v,\n", result_str, index)
			// middlewares.ServerErrResponse(err.Error(), c.Writer)
			continue
		} else {
			res, _ := json.Marshal(result.InsertedID)
			color.Green("Inserted new device at %v\n", strings.Replace(string(res), `"`, ``, 2))
			result_str = fmt.Sprintf("%v%v,\n", result_str, strings.Replace(string(res), `"`, ``, 2))
		}
	}
	result_str = fmt.Sprintf("Inserted at %v", result_str)
	if len(failed_str) > 0 {
		result_str = fmt.Sprintf("%v Failed for Index %v", result_str, failed_str)
	}
	middlewares.SuccessMessageResponse(result_str, c.Writer)

}

var AddNewTemplate = gin.HandlerFunc(func(c *gin.Context) {
	var template models.Template
	err := c.BindJSON(&template)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in below api...")
		return
	}

	template.ID = primitive.NewObjectID()
	collection := client.Database("devicemanagerservice").Collection("templates")

	result, err := collection.InsertOne(context.TODO(), template)
	if err != nil {
		color.Red("Template Create Failed in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessMessageResponse(`Inserted new template at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
})

var EditTemplate = gin.HandlerFunc(func(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var new_template models.Template
	err := c.BindJSON(&new_template)
	if err != nil {
		color.Red("Request Decoding Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	new_template.ID = id
	update := bson.D{{Key: "$set", Value: new_template}}
	collection := client.Database("devicemanagerservice").Collection("templates")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		color.Red("Template Update Failed with ID : %v...", id)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	if res.MatchedCount == 0 {
		color.Red("Template Not Exist with ID : %v...", id)
		middlewares.ErrorResponse("Template does not exist", c.Writer)
		return
	}
	middlewares.SuccessMessageResponse("Updated", c.Writer)
})

var DeleteTemplate = gin.HandlerFunc(func(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Invalid Template ID to Delete in below api...")
		middlewares.ErrorResponse("Invalid Template ID", c.Writer)
		return
	}

	collection := client.Database("devicemanagerservice").Collection("templates")
	res, derr := collection.DeleteOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}})
	if derr != nil || res.DeletedCount == 0 {
		color.Red("Template Does Not Exist with ID:%v ...", id)
		middlewares.ErrorResponse("Template does not exist", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("Deleted", c.Writer)
})

var GetTemplateByID = gin.HandlerFunc(func(c *gin.Context) {

	var model_s map[string]models.Model
	model_s, errs := GetModelsMap()

	if errs != "" {
		color.Red("Models Database Issue in below api...")
	}

	var manufacturers map[string]models.Manufacturer
	manufacturers, errs = GetManufacturersMap()

	if errs != "" {
		color.Red("Manufacturers Database Issue in below api...")
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in below api...")
		middlewares.ErrorResponse("Invalid Device ID", c.Writer)
		return
	}

	var template models.Template
	collection := client.Database("devicemanagerservice").Collection("templates")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&template)
	if err != nil {
		color.Red("No Template Exist with ID: %v in below api...", id)
		middlewares.ErrorResponse("Template does not exist", c.Writer)
		return
	}

	if template.Type == 0 {
		template.PRFID.Manufacturer_Name = ""
		if template.PRFID.Manufacturer != "" {
			template.PRFID.Manufacturer_Name = manufacturers[template.PRFID.Manufacturer].Name
		}
		template.PRFID.Model_Name = ""
		if template.PRFID.Model != "" {
			template.PRFID.Model_Name = model_s[template.PRFID.Model].Name
		}
	} else if template.Type == 1 {
		template.BLE.Manufacturer_Name = ""
		if template.BLE.Manufacturer != "" {
			template.BLE.Manufacturer_Name = manufacturers[template.BLE.Manufacturer].Name
		}
		template.BLE.Model_Name = ""
		if template.BLE.Model != "" {
			template.BLE.Model_Name = model_s[template.BLE.Model].Name
		}
	}

	middlewares.SuccessOneRespond(template, "Template", c.Writer)
})

var BulkAction = gin.HandlerFunc(func(c *gin.Context) {

	var selected models.SelectData
	err := c.BindJSON(&selected)
	if err != nil {
		color.Red("Request Decoding Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	var objectids []primitive.ObjectID
	for _, stringID := range selected.SelectedID {
		objectID, err := primitive.ObjectIDFromHex(stringID)
		if err != nil {
			panic(err)
		}
		objectids = append(objectids, objectID)
	}
	collection := client.Database("devicemanagerservice").Collection("devices")

	if selected.Operation == "delete" {
		filter := bson.M{"_id": bson.M{"$in": objectids}}
		result, err := collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			panic(err)
		}
		middlewares.SuccessMessageResponse(fmt.Sprintf("%v devices deleted successfully", result.DeletedCount), c.Writer)
	}

})
