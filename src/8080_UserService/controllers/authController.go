package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"tracio.com/userservice/db"
	middlewares "tracio.com/userservice/handlers"
	"tracio.com/userservice/models"
	"tracio.com/userservice/validators"

	"github.com/pquerna/otp/totp"
	qrcode "github.com/skip2/go-qrcode"
)

var client = db.Dbconnect()

func UserRegister(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		color.Red("Bad Request in UserRegister()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	if ok, errors := validators.ValidateInputs(user); !ok {
		middlewares.ValidationResponse(errors, c.Writer)
		color.Red("Validation Error in UserRegister()...")
		return
	}

	collection := client.Database("userservice").Collection("users")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Password Hash Error in UserRegister()...")
		return
	}
	user.Password = string(hashedPassword)
	user.LastConnect = time.Now()
	var ex_user models.User

	err = collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).Decode(&ex_user)
	if err == nil {
		color.Red("Email %v Already Exist in UserRegister()...", user.Email)
		middlewares.ErrorResponse("The user with that email already exists...", c.Writer)
		return
	}
	err = collection.FindOne(context.TODO(), bson.D{{Key: "username", Value: user.Username}}).Decode(&ex_user)
	if err == nil {
		color.Red("Username %v Already Exist in UserRegister()...", user.Username)
		middlewares.ErrorResponse("The user with that username already exists...", c.Writer)
		return
	}

	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		color.Red("Insert User in Database failed in UserRegister()...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessMessageResponse(`Inserted new user at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
}

func UserLogin(c *gin.Context) {

	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in UserLogin()...")
		return
	}

	collection := client.Database("userservice").Collection("users")

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// Test the email against the regular expression
	var ex_user models.User
	if emailRegex.MatchString(user.Email) {
		err = collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: user.Email}}).Decode(&ex_user)
		if err != nil {
			middlewares.ErrorResponse("No User exist with that email...", c.Writer)
			color.Red("No User Exist with email %v in UserLogin()...", user.Email)
			return
		}
	} else {
		err = collection.FindOne(context.TODO(), bson.D{{Key: "username", Value: user.Email}}).Decode(&ex_user)
		if err != nil {
			middlewares.ErrorResponse("No User exist with that username...", c.Writer)
			color.Red("No User Exist with username %v in UserLogin()...", user.Username)
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(ex_user.Password), []byte(user.Password)); err != nil {
		middlewares.ErrorResponse("Invalid email or password...", c.Writer)
		color.Red("Invalid Password with email %v in UserLogin()...")
		return
	}

	filter := bson.M{"email": ex_user.Email} // Use a filter to select the user with the desired email

	update := bson.M{"$set": bson.M{"lastconnect": time.Now()}} // Use the $set modifier to update the desired field

	_, err = collection.UpdateOne(context.Background(), filter, update) // Use the UpdateOne method to apply the update

	if err != nil {
		color.Red("LastConnect Update failed in UserLogin() for ID : %v", ex_user.ID)
	}

	token, _ := middlewares.GenerateJWT(ex_user.Email, ex_user.Username, ex_user.ID.Hex())

	middlewares.SuccessMessageResponse(string(token), c.Writer)

}

func GenerateOTP(c *gin.Context) {

	var tokenString string
	if tokenString = c.Request.Header.Get("Token"); tokenString == "" {
		color.Red("Getting Token Header Failed in GenerateOTP()...")
		middlewares.AuthorizationResponse("Not Authorized", c.Writer)
		return
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		middlewares.AuthorizationResponse("Invalid JWT token", c.Writer)
		color.Red("Invalid JWT Token in GenerateOTP()...")
		return
	}

	var id string = ""
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		id = claims["ID"].(string)
	} else {
		color.Red("Claming Token Failed in GenerateOTP()...")
		middlewares.AuthorizationResponse("Invalid JWT token", c.Writer)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Primitive ID generate failed for %v in GenerateOTP()...", id)
		middlewares.ErrorResponse("Invalid User ID", c.Writer)
		return
	}

	user, err_str := GetUserByIDFunc(id)
	if err_str != "" {
		middlewares.AuthorizationResponse("Invalid JWT Token", c.Writer)
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "tracio.com",
		AccountName: user.Email,
		SecretSize:  15,
	})

	if err != nil {
		color.Red("TOTP Generation failed for email %v in GenerateOTP()...", user.Email)
		middlewares.ErrorResponse("Key Generation Issue.", c.Writer)
	}

	user.OtpSecret = key.Secret()
	user.OtpAuthUrl = key.URL()

	update := bson.D{{Key: "$set", Value: user}}
	collection := client.Database("userservice").Collection("users")

	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: objID}}, update)
	if err != nil {
		color.Red("User Update Failed with ID : %v in GenerateOTP()...", id)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	if res.MatchedCount == 0 {
		color.Red("User Not Exist with ID : %v in GenerateOTP()...", id)
		middlewares.ErrorResponse("User does not exist", c.Writer)
		return
	}

	// Generate the QR code
	qr, err := qrcode.New(user.OtpAuthUrl, qrcode.High)
	if err != nil {
		color.Red("Error generating QR code:", err)
		middlewares.ErrorResponse("QRCode Generation Failed...", c.Writer)
		return
	}

	filename := fmt.Sprintf("uploaded/totp/%v", user.ID.Hex())
	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		color.Red("Error creating output file:", err)
		middlewares.ErrorResponse("QRCode Generation Failed...", c.Writer)
		return
	}
	defer file.Close()

	// Save the QR code as a PNG file
	err = qr.Write(256, file)
	if err != nil {
		color.Red("Error saving QR code:", err)
		middlewares.ErrorResponse("QRCode Generation Failed...", c.Writer)
		return
	}

	color.Cyan("Generated Successfully")
	content, _ := readFileAsBase64(filename)
	middlewares.SuccessMessageResponse(content, c.Writer)
}

func readFileAsBase64(filepath string) (string, error) {
	// Read the file data
	fileData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	// Encode the file data as base64
	encodedData := base64.StdEncoding.EncodeToString(fileData)
	return encodedData, nil
}

// UploadFileEndpoint -> upload file
func UploadFileEndpoint(c *gin.Context) {
	file, err := c.FormFile("file")
	id := c.PostForm("id")
	// fileName := c.PostForm("file_name")
	if err != nil {
		color.Red("User Image Upload Failed in UploadFileEndpoint()...")
		return
	}

	filename := fmt.Sprintf("uploaded/photo/%v", id)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		color.Red("File Save Error Upload User Image in UploadFileEndpoint()...")
		middlewares.ErrorResponse("File Save Error", c.Writer)
		return
	}
	defer f.Close()

	if err := c.SaveUploadedFile(file, f.Name()); err != nil {
		color.Red("File Save Error Upload User Image in UploadFileEndpoint()...")
		middlewares.ErrorResponse("File Save Error", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("Uploaded Successfully", c.Writer)
}

func VerifyOTP(c *gin.Context) {

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		color.Red("Bad Request in VerifyOTP()...", c.Writer)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	token, ok := req["token"]
	if !ok {
		middlewares.ServerErrResponse("Bad Request...", c.Writer)
		color.Red("Bad Request in VerifyOTP()...")
		return
	}

	var tokenString string
	if tokenString = c.Request.Header.Get("Token"); tokenString != "" {

	} else {
		color.Red("Bad Request in VerifyOTP()...")
		middlewares.AuthorizationResponse("Not Authorized", c.Writer)
		return
	}

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		// Handle error
		color.Red("Claming Token Failed in VerifyOTP()...")
		middlewares.AuthorizationResponse("Invalid JWT token", c.Writer)
		return
	}

	var id string = ""
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		id = claims["ID"].(string)
	} else {
		color.Red("Claming Token Failed in VerifyOTP()...")
		middlewares.AuthorizationResponse("Invalid JWT token", c.Writer)
	}

	user, err_str := GetUserByIDFunc(id)

	if err_str != "" {
		middlewares.AuthorizationResponse("Invalid JWT Token", c.Writer)
		return
	}
	color.Red("%v %v %v", token, user.Email, user.OtpSecret)

	valid := totp.Validate(token, user.OtpSecret)
	if !valid {
		color.Red("OTP Validation Failed in VerifyOTP()...")
		middlewares.ErrorResponse("OTP Validation Failed.", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("OTP Validation Success.", c.Writer)

}
