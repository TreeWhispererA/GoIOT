package middlewares

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tracio.com/userservice/db"
	"tracio.com/userservice/models"
)

var client = db.Dbconnect()

var mySigningKey = []byte(DotEnvVariable("JWT_SECRET"))

// IsAuthorized -> verify jwt header
func IsAuthorized(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		if tokenString = c.Request.Header.Get("Token"); tokenString != "" {
		} else {
			color.Red("Not Authorized...")
			AuthorizationResponse("Not Authorized", c.Writer)
			return
		}

		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})

		if err != nil {
			// Handle error
			color.Red("Invalid JWT token...")
			AuthorizationResponse("Invalid JWT token", c.Writer)
			return
		}

		var id string = ""
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			id = claims["ID"].(string)
		} else {
			color.Red("Invalid JWT token...")
			AuthorizationResponse("Invalid JWT token", c.Writer)
			return
		}

		user, str_err := GetUserByIDFunc(id)
		user.Password = ""

		if str_err == "" {
			color.Green("Username : %v", user.Username)
			color.Green("Email : %v", user.Email)
			next(c)
			return
		}
		color.Red("Invalid JWT token...")
		AuthorizationResponse("Invalid JWT Token", c.Writer)
	}

}

// GenerateJWT -> generate jwt
func GenerateJWT(email string, username string, ID string) (string, error) {

	if email == "" {
		email = "json037923@gmail.com"
	}
	if username == "" {
		username = "aidenlee"
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["username"] = username
	claims["ID"] = ID
	claims["exp"] = time.Now().Add(time.Minute * 14400).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		color.Red("Token Generate Failed...")
		return "", err
	}

	return tokenString, nil
}

func GetUserByIDFunc(id string) (models.User, string) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Primitive ID generate failed for %v in GetUserByIDFunc()...", id)
		return models.User{}, "Invalid User ID"
	}

	var user models.User
	collection := client.Database("userservice").Collection("users")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&user)
	if err != nil {
		color.Red("No User Exist with ID: %v in GetUserByIDFunc()...", id)
		return models.User{}, "User Does not Exist"
	}
	return user, ""
}
