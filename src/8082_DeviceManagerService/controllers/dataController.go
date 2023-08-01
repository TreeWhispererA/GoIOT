package controllers

import (
	"context"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson"

	"tracio.com/devicemanagerservice/models"
)

func GetModelsMap() (map[string]models.Model, string) {
	var model_s map[string]models.Model
	model_s = make(map[string]models.Model)

	collection := client.Database("devicemanagerservice").Collection("models")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Models Database Issue in below api...")
		return nil, "Database Connection Issue in below api..."
	}
	for cursor.Next(context.TODO()) {
		var model models.Model
		err := cursor.Decode(&model)
		if err != nil {
			color.Red("One Model Decode Failed in below api...")
		} else {
			model_s[model.ID.Hex()] = model
		}
	}
	return model_s, ""
}

func GetManufacturersMap() (map[string]models.Manufacturer, string) {
	var manufacturers map[string]models.Manufacturer
	manufacturers = make(map[string]models.Manufacturer)

	collection := client.Database("devicemanagerservice").Collection("manufacturers")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Manufacturers Database Issue in below api...")
		return nil, "Database Connection Issue in below api..."
	}
	for cursor.Next(context.TODO()) {
		var manufacturer models.Manufacturer
		err := cursor.Decode(&manufacturer)
		if err != nil {
			color.Red("One Device Decode Failed in below api...")
		} else {
			manufacturers[manufacturer.ID.Hex()] = manufacturer
		}
	}
	return manufacturers, ""
}

// func GetUserByIDFunc(id string) (models.User, string) {

// 	var roles map[primitive.ObjectID]models.Role
// 	roles, errs := GetRolesMap()

// 	if errs != "" {
// 		color.Red("Roles Database Issue in below api...")
// 		return models.User{}, "Getting Roles Information Issue"
// 	}

// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		color.Red("Getting Param Issue in below api...")
// 		return models.User{}, "Invalid User ID"
// 	}

// 	var user models.User
// 	collection := client.Database("userservice").Collection("users")
// 	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&user)
// 	if err != nil {
// 		color.Red("No User Exist with ID: %v in below api...", id)
// 		return models.User{}, "User Does not Exist"
// 	}
// 	user.Password = ""

// 	for _, groupID := range user.Role {
// 		id, _ := primitive.ObjectIDFromHex(groupID)

// 		simpleRole := models.BriefRole{}

// 		role := roles[id]
// 		if _, ok := roles[id]; !ok {
// 			color.Cyan("Error in getting role data with ID : %v", id)
// 		} else {
// 			simpleRole = models.BriefRole{ID: role.ID.Hex(), Color: role.Color, Name: role.Name}
// 			user.BriefRole = append(user.BriefRole, simpleRole)
// 		}
// 	}

// 	return user, ""
// }
