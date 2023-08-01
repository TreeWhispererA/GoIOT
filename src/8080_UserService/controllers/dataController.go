package controllers

import (
	"context"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// "go.mongodb.org/mongo-driver/mongo/options"

	"tracio.com/userservice/models"
)

func GetRolesMap() (map[primitive.ObjectID]models.Role, string) {
	var roles map[primitive.ObjectID]models.Role
	roles = make(map[primitive.ObjectID]models.Role)

	collection := client.Database("userservice").Collection("roles")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Roles Database Issue in GetRolesMap()...")
		return nil, "Database Connection Issue in GetRolesMap()..."
	}
	for cursor.Next(context.TODO()) {
		var role models.Role
		err := cursor.Decode(&role)
		if err != nil {
			color.Red("One Role Decode Failed in GetRolesMap()...")
		} else {
			roles[role.ID] = role
		}
	}
	return roles, ""
}

func GetUserByIDFunc(id string) (models.User, string) {

	var roles map[primitive.ObjectID]models.Role
	roles, errs := GetRolesMap()

	if errs != "" {
		color.Red("Roles Database Issue GetUserByIDFunc...")
		return models.User{}, "Getting Roles Information Issue"
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue GetUserByIDFunc...")
		return models.User{}, "Invalid User ID"
	}

	var user models.User
	collection := client.Database("userservice").Collection("users")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&user)
	if err != nil {
		color.Red("No User Exist with ID: %v GetUserByIDFunc...", id)
		return models.User{}, "User Does not Exist"
	}

	for _, groupID := range user.Role {
		id, _ := primitive.ObjectIDFromHex(groupID)

		simpleRole := models.BriefRole{}

		role := roles[id]
		if _, ok := roles[id]; !ok {
			color.Cyan("Error in getting role data with ID : %v", id)
		} else {
			simpleRole = models.BriefRole{ID: role.ID.Hex(), Color: role.Color, Name: role.Name}
			user.BriefRole = append(user.BriefRole, simpleRole)
		}
	}

	return user, ""
}
