package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/mongo/options"
	middlewares "tracio.com/userservice/handlers"
	"tracio.com/userservice/models"
	"tracio.com/userservice/validators"
)

var mySigningKey = []byte(middlewares.DotEnvVariable("JWT_SECRET"))

func Test(c *gin.Context) {
	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}

// AddNewUser -> create new user
func AddNewUser(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in AddNewUser()...")
		return
	}
	if ok, errors := validators.ValidateInputs(user); !ok {
		middlewares.ValidationResponse(errors, c.Writer)
		color.Red("User Validation Failed in AddNewUser()...")
		return
	}
	collection := client.Database("userservice").Collection("users")

	user.ID = primitive.NewObjectID()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Hash Password Failed in AddNewUser()...")
		return
	}
	user.Password = string(hashedPassword)

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		color.Red("User Register Failed in AddNewUser()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessMessageResponse(`Inserted new user at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
}

// AddNewRole -> create new role
func AddNewRole(c *gin.Context) {
	var role models.Role
	err := c.BindJSON(&role)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in AddNewRole()...")
		return
	}
	if ok, errors := validators.ValidateInputs(role); !ok {
		middlewares.ValidationResponse(errors, c.Writer)
		color.Red("Role Validation Error in AddNewRole()...")
		return
	}
	collection := client.Database("userservice").Collection("roles")

	err = collection.FindOne(context.TODO(), bson.D{{Key: "name", Value: role.Name}}).Decode(&role)
	if err == nil {
		middlewares.ErrorResponse("The role with that name already exists...", c.Writer)
		color.Red("The Role with This Name already exists in AddNewRole()...")
		return
	}

	role.ID = primitive.NewObjectID()
	result, err := collection.InsertOne(context.TODO(), role)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Insertion Role Failed in AddNewRole()...")
		return
	}

	res, _ := json.Marshal(result.InsertedID)
	res_str := strings.Replace(string(res), `"`, ``, 2)
	collection = client.Database("userservice").Collection("users")

	userIDs := make([]primitive.ObjectID, len(role.CurrentUsers))
	for i, userIDStr := range role.CurrentUsers {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			color.Red("Issue occured converting to primitive.ObjectID in AddNewRole()...")
		}
		userIDs[i] = userID
	}

	// Update the users with the new role ID and remove the role ID from other users
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	update_user := bson.M{
		"$addToSet": bson.M{"role": res_str},
	}
	_, err = collection.UpdateMany(context.Background(), filter, update_user)
	if err != nil {
		color.Red(err.Error())
	}

	filter = bson.M{"_id": bson.M{"$nin": userIDs}, "role": bson.M{"$in": []string{res_str}}}
	update_user = bson.M{"$pull": bson.M{"role": res_str}}
	_, err = collection.UpdateMany(context.Background(), filter, update_user)
	if err != nil {
		color.Red(err.Error())
	}

	middlewares.SuccessMessageResponse(`Inserted new role at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
	color.Green(`Inserted new role at ` + strings.Replace(string(res), `"`, ``, 2))
}

// GetRoles -> get all roles data
var GetRoles = gin.HandlerFunc(func(c *gin.Context) {
	var usercount map[string]int
	usercount = make(map[string]int)

	collection := client.Database("userservice").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Users Database Issue in GetRoles()...")
		middlewares.ServerErrResponse("Database Connection Issue...", c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			color.Red("One Role Decode Failed in GetRoles()...")
		} else {
			for _, roleid := range user.Role {
				usercount[roleid]++
			}
		}
	}

	var roles []*models.Role

	collection = client.Database("userservice").Collection("roles")
	cursor, err = collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Roles Database Issue in GetRoles()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	for cursor.Next(context.TODO()) {
		var role models.Role
		err := cursor.Decode(&role)
		if err != nil {
			color.Red("Decodig one Role Failed in GetRoles()...")
		} else {
			role.Users = usercount[role.ID.Hex()]
			roles = append(roles, &role)
		}
	}
	if err := cursor.Err(); err != nil {
		color.Red("Something Wrong in GetRoles in GetRoles()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	middlewares.SuccessArrRespond(roles, "Role", c.Writer)
})

// GetUsers -> get all users data
var GetUsers = gin.HandlerFunc(func(c *gin.Context) {

	var roles map[primitive.ObjectID]models.Role
	roles, error := GetRolesMap()

	if error != "" {
		color.Red("Roles Database Issue in GetUsers()...")
		middlewares.ServerErrResponse(error, c.Writer)
		return
	}

	var users []*models.User
	collection := client.Database("userservice").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.D{{}}, options.Find().SetProjection(bson.M{"password": 0}))
	if err != nil {
		color.Red("Users Database Issue in GetUsers()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var user *models.User
		err := cursor.Decode(&user)
		user.BriefRole = []models.BriefRole{}
		if err != nil {
			color.Red("One User Decode Failed in GetUsers()...")
		} else {
			for _, groupID := range user.Role {
				id, _ := primitive.ObjectIDFromHex(groupID)

				simpleRole := models.BriefRole{}

				role := roles[id]
				if _, ok := roles[id]; !ok {
					color.Cyan("Error in getting role : %v for user: %v", id, user.ID)
				} else {
					simpleRole = models.BriefRole{ID: role.ID.Hex(), Color: role.Color, Name: role.Name}
					user.BriefRole = append(user.BriefRole, simpleRole)
				}
			}
			photo, err := readFileAsBase64(fmt.Sprintf("uploaded/photo/%v", user.ID.Hex()))
			if err != nil {
				user.Photo = "none"
			} else {
				user.Photo = photo
			}

			users = append(users, user)
		}
	}
	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetUsers in GetUsers()...")
		return
	}

	middlewares.SuccessArrRespond(users, "User", c.Writer)
})

// GetUsersByRole -> get users for specific role id
var GetUsersByRole = gin.HandlerFunc(func(c *gin.Context) {

	var roles map[primitive.ObjectID]models.Role
	roles, error := GetRolesMap()

	if error != "" {
		color.Red("Roles Database Issue in GetUsers()...")
		middlewares.ServerErrResponse(error, c.Writer)
		return
	}

	id := c.Param("id")

	collection := client.Database("userservice").Collection("users")
	cur, err := collection.Find(context.TODO(), bson.M{"role": bson.M{"$in": []string{id}}})
	if err != nil {
		color.Red("Now User exist with Role ID: %v in GetUsers()...", id)
		middlewares.ErrorResponse("Users do not exist with role", c.Writer)
		return
	}

	var results []*models.User
	for cur.Next(context.Background()) {
		var result *models.User
		err := cur.Decode(&result)
		if err != nil {
			color.Red("Now User exist with Role ID: %v in GetUsers()...", id)
			middlewares.ErrorResponse("Error occured getting users.", c.Writer)
			return
		} else {
			for _, groupID := range result.Role {
				id, _ := primitive.ObjectIDFromHex(groupID)

				simpleRole := models.BriefRole{}

				role := roles[id]
				if _, ok := roles[id]; !ok {
					color.Cyan("Error in getting role : %v for user: %v", id, result.ID)
				} else {
					simpleRole = models.BriefRole{ID: role.ID.Hex(), Color: role.Color, Name: role.Name}
					result.BriefRole = append(result.BriefRole, simpleRole)
				}
			}
			photo, err := readFileAsBase64(fmt.Sprintf("uploaded/photo/%v", result.ID.Hex()))
			if err != nil {
				result.Photo = "none"
			} else {
				result.Photo = photo
			}

			results = append(results, result)
		}
	}

	middlewares.SuccessArrRespond(results, "User", c.Writer)
})

// GetRoleByID -> get role by id
var GetRoleByID = gin.HandlerFunc(func(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Converting Primitive ID issue in GetRoleByID()...", id)
		middlewares.ErrorResponse("Invalid Role ID", c.Writer)
		return
	}

	var role models.Role
	collection := client.Database("userservice").Collection("roles")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&role)
	if err != nil {
		color.Red("Role getting failed with id %v in GetRoleByID()...", id)
		middlewares.ErrorResponse("Role does not exist", c.Writer)
		return
	}

	middlewares.SuccessOneRespond(role, "Role", c.Writer)
})

// GetUserByID -> get user by id
var GetUserByID = gin.HandlerFunc(func(c *gin.Context) {

	var roles map[primitive.ObjectID]models.Role
	roles, error := GetRolesMap()

	if error != "" {
		color.Red("Roles Database Issue in GetUserByID()...")
		middlewares.ServerErrResponse(error, c.Writer)
		return
	}

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in GetUserByID()...")
		middlewares.ErrorResponse("Invalid User ID", c.Writer)
		return
	}

	var user models.User
	collection := client.Database("userservice").Collection("users")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&user)
	if err != nil {
		color.Red("No User Exist with ID: %v in GetUserByID()...", id)
		middlewares.ErrorResponse("User does not exist", c.Writer)
		return
	}
	user.Password = ""
	photo, err := readFileAsBase64(fmt.Sprintf("uploaded/photo/%v", user.ID.Hex()))
	if err != nil {
		user.Photo = "none"
	} else {
		user.Photo = photo
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

	middlewares.SuccessOneRespond(user, "User", c.Writer)
})

// GetUserByID -> get user by id
var GetInfo = gin.HandlerFunc(func(c *gin.Context) {

	tokenString := c.Request.Header.Get("Token")
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		// Handle error
		middlewares.ErrorResponse("Invalid JWT Token", c.Writer)
		fmt.Println("Error parsing JWT token in GetInfo():", err)
		return
	}

	var id string = ""
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		id = claims["ID"].(string)
	} else {
		middlewares.ErrorResponse("Invalid JWT Token", c.Writer)
		fmt.Println("Invalid JWT token")
	}

	var user models.User
	user, str_err := GetUserByIDFunc(id)
	user.Password = ""
	photo, err := readFileAsBase64(fmt.Sprintf("uploaded/photo/%v", user.ID.Hex()))
	if err != nil {
		user.Photo = "none"
	} else {
		user.Photo = photo
	}

	if str_err != "" {
		middlewares.ErrorResponse(str_err, c.Writer)
		return
	}

	middlewares.SuccessOneRespond(user, "User", c.Writer)
})

// DeleteRole -> delete role by id
var DeleteRole = gin.HandlerFunc(func(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Invalid Role ID to Delete DeleteRole()...")
		middlewares.ErrorResponse("Invalid Role ID", c.Writer)
		return
	}

	collection := client.Database("userservice").Collection("roles")
	res, derr := collection.DeleteOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}})
	if derr != nil || res.DeletedCount == 0 {
		color.Red("Role Does Not Exist with ID:%v ...", id)
		middlewares.ErrorResponse("Role does not exist", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("Deleted", c.Writer)
})

// DeleteUser -> delete user by id
var DeleteUser = gin.HandlerFunc(func(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var user models.User

	collection := client.Database("userservice").Collection("users")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		middlewares.ErrorResponse("User does not exist", c.Writer)
		color.Red("Invalid User ID to Delete DeleteUser()...")
		return
	}
	_, derr := collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}})
	if derr != nil {
		middlewares.ServerErrResponse(derr.Error(), c.Writer)
		color.Red("User Does Not Exist with ID:%v ...", id)
		return
	}
	middlewares.SuccessMessageResponse("Deleted", c.Writer)
})

// EditUser -> update user by id
var EditUser = gin.HandlerFunc(func(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var new_user models.User
	err := c.BindJSON(&new_user)
	if err != nil {
		color.Red("Request Decoding Issue EditUser()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(new_user.Password), bcrypt.DefaultCost)
	if err != nil {
		color.Red("Hash Password Generation Failed EditUser()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	exist_user, err_str := GetUserByIDFunc(c.Param("id"))

	if err_str != "" {
		color.Red("User Not Exist with ID : %v...", c.Param("id"))
		middlewares.ErrorResponse("User does not exist", c.Writer)
		return
	}

	new_user.ID = id
	new_user.OtpAuthUrl = exist_user.OtpAuthUrl
	new_user.OtpSecret = exist_user.OtpSecret

	if new_user.Password == "" {
		new_user.Password = exist_user.Password
	} else {
		new_user.Password = string(hashedPassword)
	}
	update := bson.D{{Key: "$set", Value: new_user}}
	collection := client.Database("userservice").Collection("users")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		color.Red("User Update Failed with ID : %v...", id)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	if res.MatchedCount == 0 {
		color.Red("User Not Exist with ID : %v...", id)
		middlewares.ErrorResponse("User does not exist", c.Writer)
		return
	}
	middlewares.SuccessMessageResponse("Updated", c.Writer)
})

var ChangePassword = gin.HandlerFunc(func(c *gin.Context) {

	var change_info models.PasswordChangeInfo
	err := c.BindJSON(&change_info)
	if err != nil {
		color.Red("Request Decoding Issue in ChangePassword()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	id, _ := primitive.ObjectIDFromHex(change_info.ID)
	new_pass := change_info.NewPassword

	new_pass_hash, err := bcrypt.GenerateFromPassword([]byte(new_pass), bcrypt.DefaultCost)
	if err != nil {
		color.Red("Hash Password Generation Failed in ChangePassword()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	new_pass_hash_string := string(new_pass_hash)

	var user models.User
	collection := client.Database("userservice").Collection("users")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&user)
	if err != nil {
		color.Red("No User Exist with ID: %v in ChangePassword()...", change_info.ID)
		middlewares.ErrorResponse("User does not exist", c.Writer)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(change_info.OldPassword)); err != nil {
		color.Red("Password Incorrect ID: %v in ChangePassword()...", change_info.ID)
		middlewares.ErrorResponse("Current password incorrect.", c.Writer)
		return
	}

	user.Password = new_pass_hash_string
	update := bson.D{{Key: "$set", Value: user}}
	collection = client.Database("userservice").Collection("users")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		color.Red("User Update Failed with ID : %v...", id)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	if res.MatchedCount == 0 {
		color.Red("User Not Exist with ID : %v...", id)
		middlewares.ErrorResponse("User does not exist", c.Writer)
		return
	}
	middlewares.SuccessMessageResponse("Updated", c.Writer)
})

// EditRole -> update role by id
var EditRole = gin.HandlerFunc(func(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Request Decoding Issue in EditRole()...")
		middlewares.ErrorResponse("Invalid Role ID", c.Writer)
		return
	}

	var new_role models.Role
	if err := c.ShouldBindJSON(&new_role); err != nil {
		color.Red("Hash Password Generation Failed in EditRole()...")
		middlewares.ErrorResponse("Invalid Request Body", c.Writer)
		return
	}
	new_role.ID = objID
	update := bson.D{{Key: "$set", Value: new_role}}

	collection := client.Database("userservice").Collection("roles")
	res, derr := collection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)
	if derr != nil || res.MatchedCount == 0 {
		color.Red("Role Not Exist with ID : %v...", id)
		middlewares.ErrorResponse("Role does not exist", c.Writer)
		return
	}

	collection = client.Database("userservice").Collection("users")

	userIDs := make([]primitive.ObjectID, len(new_role.CurrentUsers))
	for i, userIDStr := range new_role.CurrentUsers {
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			color.Red("Issue occured converting to primitive.ObjectID in EditRole()...")
			middlewares.ServerErrResponse(err.Error(), c.Writer)
		}
		userIDs[i] = userID
	}

	// Update the users with the new role ID and remove the role ID from other users
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	update_user := bson.M{
		"$addToSet": bson.M{"role": id},
	}
	_, err = collection.UpdateMany(context.Background(), filter, update_user)
	if err != nil {
		color.Red(err.Error())
		middlewares.ServerErrResponse(err.Error(), c.Writer)
	}

	filter = bson.M{"_id": bson.M{"$nin": userIDs}, "role": bson.M{"$in": []string{id}}}
	update_user = bson.M{"$pull": bson.M{"role": id}}
	_, err = collection.UpdateMany(context.Background(), filter, update_user)
	if err != nil {
		color.Red(err.Error())
		middlewares.ServerErrResponse(err.Error(), c.Writer)
	}

	middlewares.SuccessMessageResponse("Successfully Updated", c.Writer)
})

var GetSmartRoles = gin.HandlerFunc(func(c *gin.Context) {

	var usercount map[string]int
	usercount = make(map[string]int)

	collection := client.Database("userservice").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Users Database Issue in GetSmartRoles()...")
		middlewares.ServerErrResponse("Database Connection Issue...", c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			color.Red("One User Decode Failed in GetSmartRoles()...")
		} else {
			for _, roleid := range user.Role {
				usercount[roleid]++
			}
		}
	}

	var smartrule models.SmartRule
	err = c.BindJSON(&smartrule)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in GetSmartRoles()...")
		return
	}

	var roles []*models.Role
	var results []*models.Role

	collection = client.Database("userservice").Collection("roles")
	cursor, err = collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Roles Database Issue in GetSmartRoles()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var role models.Role
		err := cursor.Decode(&role)
		if err != nil {
			color.Red("Decodig one Role Failed in GetSmartRoles()...")
		} else {
			role.Users = usercount[role.ID.Hex()]
			roles = append(roles, &role)
		}
	}

	if err := cursor.Err(); err != nil {
		color.Red("Something Wrong in GetRoles in GetSmartRoles()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	if smartrule.Filtering {
		for _, role := range roles {
			if smartrule.Filter.Reports.Access == -1 && role.PagePermission.Reports.Access {
				continue
			}
			if smartrule.Filter.LocalLive.Access == -1 && role.PagePermission.LocalLive.Access {
				continue
			}
			if smartrule.Filter.History.Access == -1 && role.PagePermission.History.Access {
				continue
			}
			if smartrule.Filter.Alert.Access == -1 && role.PagePermission.Alert.Access {
				continue
			}
			if smartrule.Filter.SiteManager.Access == -1 && role.PagePermission.SiteManager.Access {
				continue
			}
			if smartrule.Filter.DeviceManager.Access == -1 && role.PagePermission.DeviceManager.Access {
				continue
			}
			if smartrule.Filter.RulesManager.Access == -1 && role.PagePermission.RulesManager.Access {
				continue
			}
			if smartrule.Filter.Analytics.Access == -1 && role.PagePermission.Analytics.Access {
				continue
			}
			if smartrule.Filter.User.Access == -1 && role.PagePermission.User.Access {
				continue
			}
			if smartrule.Filter.System.Access == -1 && role.PagePermission.System.Access {
				continue
			}

			if smartrule.Filter.Reports.Access == 1 {
				if !role.PagePermission.Reports.Access {
					continue
				}
				if smartrule.Filter.Reports.Group1 == 1 && !role.PagePermission.Reports.ObjectReport.Access {
					continue
				}
				if smartrule.Filter.Reports.Group2 == 1 && !role.PagePermission.Reports.TagReport.Access {
					continue
				}
				if smartrule.Filter.Reports.Group1 == -1 && role.PagePermission.Reports.ObjectReport.Access {
					continue
				}
				if smartrule.Filter.Reports.Group2 == -1 && role.PagePermission.Reports.TagReport.Access {
					continue
				}
			}
			if smartrule.Filter.LocalLive.Access == 1 {
				if !role.PagePermission.LocalLive.Access {
					continue
				}
			}
			if smartrule.Filter.History.Access == 1 {
				if !role.PagePermission.History.Access {
					continue
				}
				if smartrule.Filter.History.Group1 == 1 && !role.PagePermission.History.ObjectEvent.Access {
					continue
				}
				if smartrule.Filter.History.Group2 == 1 && !role.PagePermission.History.TagBlink.Access {
					continue
				}
				if smartrule.Filter.History.Group1 == -1 && role.PagePermission.History.ObjectEvent.Access {
					continue
				}
				if smartrule.Filter.History.Group2 == -1 && role.PagePermission.History.TagBlink.Access {
					continue
				}
			}
			if smartrule.Filter.Alert.Access == 1 {
				if !role.PagePermission.Alert.Access {
					continue
				}
				if smartrule.Filter.Alert.Group1 == 1 && !role.PagePermission.Alert.ObjectAlert.Access {
					continue
				}
				if smartrule.Filter.Alert.Group2 == 1 && !role.PagePermission.Alert.SystemAlert.Access {
					continue
				}
				if smartrule.Filter.Alert.Group1 == -1 && role.PagePermission.Alert.ObjectAlert.Access {
					continue
				}
				if smartrule.Filter.Alert.Group2 == -1 && role.PagePermission.Alert.SystemAlert.Access {
					continue
				}
			}
			if smartrule.Filter.SiteManager.Access == 1 {
				if !role.PagePermission.SiteManager.Access {
					continue
				}
			}
			if smartrule.Filter.DeviceManager.Access == 1 {
				if !role.PagePermission.DeviceManager.Access {
					continue
				}
			}
			if smartrule.Filter.RulesManager.Access == 1 {
				if !role.PagePermission.RulesManager.Access {
					continue
				}
			}
			if smartrule.Filter.Analytics.Access == 1 {
				if !role.PagePermission.Analytics.Access {
					continue
				}
			}
			if smartrule.Filter.User.Access == 1 {
				if !role.PagePermission.User.Access {
					continue
				}
			}
			if smartrule.Filter.System.Access == 1 {
				if !role.PagePermission.System.Access {
					continue
				}
			}
			results = append(results, role)
		}
	} else {
		results = append(results, roles...)
	}

	if smartrule.Sorting {
		switch smartrule.SortIndex {
		case "roles":
			sort.Slice(results, func(i, j int) bool {
				return results[i].Name > results[j].Name
			})
		case "description":
			sort.Slice(results, func(i, j int) bool {
				return results[i].Description > results[j].Description
			})
		case "users":
			sort.Slice(results, func(i, j int) bool {
				return results[i].Users > results[j].Users
			})
		case "report":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.Reports.Access && !results[j].PagePermission.Reports.Access
			})
		case "locallive":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.LocalLive.Access && !results[j].PagePermission.LocalLive.Access
			})
		case "history":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.History.Access && !results[j].PagePermission.History.Access
			})
		case "alerts":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.Alert.Access && !results[j].PagePermission.Alert.Access
			})
		case "sitemanager":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.SiteManager.Access && !results[j].PagePermission.SiteManager.Access
			})
		case "devicemanager":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.DeviceManager.Access && !results[j].PagePermission.DeviceManager.Access
			})
		case "rulesmanager":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.RulesManager.Access && !results[j].PagePermission.RulesManager.Access
			})
		case "analytics":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.Analytics.Access && !results[j].PagePermission.Analytics.Access
			})
		case "user":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.User.Access && !results[j].PagePermission.User.Access
			})
		case "system":
			sort.Slice(results, func(i, j int) bool {
				return results[i].PagePermission.System.Access && !results[j].PagePermission.System.Access
			})
		}

		if smartrule.Ascending == false {
			for i := 0; i < len(results)/2; i++ {
				j := len(results) - i - 1
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	totalpages := len(results)
	color.Red("%v %v %v", totalpages, (smartrule.PageNum-1)*smartrule.PerPage+1, smartrule.PageNum*smartrule.PerPage)
	if (smartrule.PageNum-1)*smartrule.PerPage+1 <= totalpages {
		max_page := smartrule.PageNum * smartrule.PerPage
		if max_page > totalpages {
			max_page = totalpages
		}
		results = results[(smartrule.PageNum-1)*smartrule.PerPage : max_page]
	} else {
		// pagenum = 1
		max_page := smartrule.PerPage
		if max_page > totalpages {
			max_page = totalpages
		}
		results = results[0:max_page]
	}

	middlewares.SuccessSmartRuleRespond(results, "Role", totalpages, c.Writer)
})

var GetSmartUsers = gin.HandlerFunc(func(c *gin.Context) {

	var roles map[primitive.ObjectID]models.Role
	roles, error := GetRolesMap()

	if error != "" {
		color.Red("Roles Database Issue in GetSmartUsers...")
		middlewares.ServerErrResponse(error, c.Writer)
		return
	}

	var smartrule models.SmartRule
	err := c.BindJSON(&smartrule)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in GetSmartUsers...")
		return
	}

	var users []*models.User
	var results []*models.User

	collection := client.Database("userservice").Collection("users")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Users Database Issue in GetSmartUsers...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	for cursor.Next(context.TODO()) {
		var user *models.User
		err = cursor.Decode(&user)
		user.BriefRole = []models.BriefRole{}
		for _, groupID := range user.Role {
			id, _ := primitive.ObjectIDFromHex(groupID)

			simpleRole := models.BriefRole{}

			role := roles[id]
			if _, ok := roles[id]; !ok {
				color.Cyan("Error in getting role : %v for user: %v", id, user.ID)
			} else {
				simpleRole = models.BriefRole{ID: role.ID.Hex(), Color: role.Color, Name: role.Name}
				user.BriefRole = append(user.BriefRole, simpleRole)
			}
		}

		photo, err := readFileAsBase64(fmt.Sprintf("uploaded/photo/%v", user.ID.Hex()))
		if err != nil {
			user.Photo = "none"
		} else {
			user.Photo = photo
		}
		user.Password = ""
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		color.Red("Something Wrong in GetRoles in GetSmartUsers...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	if smartrule.Filtering {
		// 	for _, role := range roles {
		// 		if smartrule.Filter.Reports.Access == -1 && role.PagePermission.Reports.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.LocalLive.Access == -1 && role.PagePermission.LocalLive.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.History.Access == -1 && role.PagePermission.History.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.Alert.Access == -1 && role.PagePermission.Alert.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.SiteManager.Access == -1 && role.PagePermission.SiteManager.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.DeviceManager.Access == -1 && role.PagePermission.DeviceManager.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.RulesManager.Access == -1 && role.PagePermission.RulesManager.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.Analytics.Access == -1 && role.PagePermission.Analytics.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.User.Access == -1 && role.PagePermission.User.Access {
		// 			continue
		// 		}
		// 		if smartrule.Filter.System.Access == -1 && role.PagePermission.System.Access {
		// 			continue
		// 		}

		// 		if smartrule.Filter.Reports.Access == 1 {
		// 			if !role.PagePermission.Reports.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.Reports.Group1 == 1 && !role.PagePermission.Reports.ObjectReport.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.Reports.Group2 == 1 && !role.PagePermission.Reports.TagReport.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.Reports.Group1 == -1 && role.PagePermission.Reports.ObjectReport.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.Reports.Group2 == -1 && role.PagePermission.Reports.TagReport.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.LocalLive.Access == 1 {
		// 			if !role.PagePermission.LocalLive.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.History.Access == 1 {
		// 			if !role.PagePermission.History.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.History.Group1 == 1 && !role.PagePermission.History.ObjectEvent.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.History.Group2 == 1 && !role.PagePermission.History.TagBlink.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.History.Group1 == -1 && role.PagePermission.History.ObjectEvent.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.History.Group2 == -1 && role.PagePermission.History.TagBlink.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.Alert.Access == 1 {
		// 			if !role.PagePermission.Alert.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.Alert.Group1 == 1 && !role.PagePermission.Alert.ObjectAlert.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.Alert.Group2 == 1 && !role.PagePermission.Alert.SystemAlert.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.Alert.Group1 == -1 && role.PagePermission.Alert.ObjectAlert.Access {
		// 				continue
		// 			}
		// 			if smartrule.Filter.Alert.Group2 == -1 && role.PagePermission.Alert.SystemAlert.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.SiteManager.Access == 1 {
		// 			if !role.PagePermission.SiteManager.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.DeviceManager.Access == 1 {
		// 			if !role.PagePermission.DeviceManager.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.RulesManager.Access == 1 {
		// 			if !role.PagePermission.RulesManager.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.Analytics.Access == 1 {
		// 			if !role.PagePermission.Analytics.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.User.Access == 1 {
		// 			if !role.PagePermission.User.Access {
		// 				continue
		// 			}
		// 		}
		// 		if smartrule.Filter.System.Access == 1 {
		// 			if !role.PagePermission.System.Access {
		// 				continue
		// 			}
		// 		}
		// 		results = append(results, role)
		// }
	} else {
		results = append(results, users...)
	}

	if smartrule.Sorting {
		switch smartrule.SortIndex {
		case "username":
			sort.Slice(results, func(i, j int) bool {
				return results[i].Username > results[j].Username
			})
		case "firstname":
			sort.Slice(results, func(i, j int) bool {
				return results[i].Firstname > results[j].Firstname
			})
		case "lastname":
			sort.Slice(results, func(i, j int) bool {
				return results[i].Lastname > results[j].Lastname
			})
		case "description":
			sort.Slice(results, func(i, j int) bool {
				return results[i].Description > results[j].Description
			})
		case "groups":
			sort.Slice(results, func(i, j int) bool {
				if len(results[i].BriefRole) == 0 {
					return false
				} else if len(results[j].BriefRole) == 0 {
					return true
				}
				return results[i].BriefRole[0].Name > results[j].BriefRole[0].Name
			})
		case "inactive":
			sort.Slice(results, func(i, j int) bool {
				return results[i].LastConnect.Before(results[j].LastConnect)
			})
		}

		if smartrule.Ascending == false {
			for i := 0; i < len(results)/2; i++ {
				j := len(results) - i - 1
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	totalpages := len(results)
	color.Red("%v %v %v", totalpages, (smartrule.PageNum-1)*smartrule.PerPage+1, smartrule.PageNum*smartrule.PerPage)
	if (smartrule.PageNum-1)*smartrule.PerPage+1 <= totalpages {
		max_page := smartrule.PageNum * smartrule.PerPage
		if max_page > totalpages {
			max_page = totalpages
		}
		results = results[(smartrule.PageNum-1)*smartrule.PerPage : max_page]
	} else {
		// pagenum = 1
		max_page := smartrule.PerPage
		if max_page > totalpages {
			max_page = totalpages
		}
		results = results[0:max_page]
	}

	middlewares.SuccessSmartRuleRespond(results, "User", totalpages, c.Writer)
})
