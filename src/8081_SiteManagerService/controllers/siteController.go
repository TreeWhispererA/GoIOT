package controllers

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"tracio.com/sitemanagerservice/db"
	middlewares "tracio.com/sitemanagerservice/handlers"
	"tracio.com/sitemanagerservice/models"
	"tracio.com/sitemanagerservice/validators"

	"tracio.com/sitemanagerservice/font"
	"tracio.com/sitemanagerservice/wmf"

	"image"

	"golang.org/x/image/tiff"
)

var client = db.Dbconnect()

// GetSites -> get all sites data

func Test(c *gin.Context) {
	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}

var GetSites = gin.HandlerFunc(func(c *gin.Context) {

	var sites []*models.Site
	collection := client.Database("sitemanagerservice").Collection("sites")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		color.Red("Sites Database Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	count := 0
	for cursor.Next(context.TODO()) {
		count++
		var site *models.Site
		err := cursor.Decode(&site)
		if err != nil {
			color.Red("%v", count)
			color.Red("One Site Data Decode Failed in below api...")
		} else {
			sites = append(sites, site)
		}
	}
	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Something Wrong in GetSites in below api...")
		return
	}

	middlewares.SuccessArrRespond(sites, "Site", c.Writer)
})

// AddNewSite -> create new site
func AddNewSite(c *gin.Context) {
	var site models.Site
	err := c.BindJSON(&site)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Bad Request in below api...")
		return
	}
	if ok, errors := validators.ValidateInputs(site); !ok {
		middlewares.ValidationResponse(errors, c.Writer)
		color.Red("Site Validation Failed in below api...")
		return
	}
	collection := client.Database("sitemanagerservice").Collection("sites")

	site.ID = primitive.NewObjectID()
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		color.Red("Hash Password Failed in below api...")
		return
	}

	result, err := collection.InsertOne(context.TODO(), site)
	if err != nil {
		color.Red("Site Insertion Failed in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessMessageResponse(`Inserted new site at `+strings.Replace(string(res), `"`, ``, 2), c.Writer)
}

// DeleteSite -> delete site by id
var DeleteSite = gin.HandlerFunc(func(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var site models.Site

	collection := client.Database("sitemanagerservice").Collection("sites")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&site)
	if err != nil {
		middlewares.ErrorResponse("Site does not exist", c.Writer)
		color.Red("Invalid Site ID to Delete in below api...")
		return
	}
	_, derr := collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}})
	if derr != nil {
		middlewares.ServerErrResponse(derr.Error(), c.Writer)
		color.Red("Site Does Not Exist with ID:%v ...", id)
		return
	}
	middlewares.SuccessMessageResponse("Deleted", c.Writer)
})

// EditSiteByID -> Edit site by id
var EditSiteByID = gin.HandlerFunc(func(c *gin.Context) {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var new_site models.Site
	err := c.BindJSON(&new_site)
	if err != nil {
		color.Red("Request Decoding Issue in below api...")
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}

	new_site.ID = id
	update := bson.D{{Key: "$set", Value: new_site}}
	collection := client.Database("sitemanagerservice").Collection("sites")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, update)
	if err != nil {
		color.Red("Site Update Failed with ID : %v...", id)
		middlewares.ServerErrResponse(err.Error(), c.Writer)
		return
	}
	if res.MatchedCount == 0 {
		color.Red("No Site Exist with ID : %v...", id)
		middlewares.ErrorResponse("Site does not exist with ID", c.Writer)
		return
	}
	middlewares.SuccessMessageResponse("Updated", c.Writer)
})

// GetSiteByID -> get site by id
var GetSiteByID = gin.HandlerFunc(func(c *gin.Context) {

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in below api...")
		middlewares.ErrorResponse("Invalid Site ID", c.Writer)
		return
	}

	var site models.Site
	collection := client.Database("sitemanagerservice").Collection("sites")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&site)
	if err != nil {
		color.Red("No Site Exist with ID: %v in below api...", id)
		middlewares.ErrorResponse("Site does not exist", c.Writer)
		return
	}

	middlewares.SuccessOneRespond(site, "Site", c.Writer)
})

// GetSiteByID -> get site by id
var ServeMapTile = gin.HandlerFunc(func(c *gin.Context) {

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in below api...")
		middlewares.ErrorResponse("Invalid Site ID", c.Writer)
		return
	}

	var site models.Site
	collection := client.Database("sitemanagerservice").Collection("sites")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&site)
	if err != nil {
		color.Red("No Site Exist with ID: %v in below api...", id)
		middlewares.ErrorResponse("Site does not exist", c.Writer)
		return
	}

	// Construct file path for requested tile image
	if site.Level != 2 {
		color.Red("Invalid Map ID...", id)
		middlewares.ErrorResponse("Invalid Map ID", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("MapTile Server Started...", c.Writer)
})

// UploadFileEndpoint -> upload file
func UploadFileEndpoint(c *gin.Context) {
	file, err := c.FormFile("file")

	id := c.PostForm("id")

	// fileName := c.PostForm("file_name")
	if err != nil {
		color.Red("File Upload Failed in below api...")
		middlewares.ErrorResponse("File Upload Failed.", c.Writer)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in below api...")
		middlewares.ErrorResponse("Invalid Site ID", c.Writer)
		return
	}

	var site models.Site
	collection := client.Database("sitemanagerservice").Collection("sites")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&site)
	if err != nil {
		color.Red("No Site Exist with ID: %v in below api...", id)
		middlewares.ErrorResponse("Site does not exist", c.Writer)
		return
	}

	if site.Level != 2 {
		color.Red("Invalid Map ID...", id)
		middlewares.ErrorResponse("Invalid Map...", c.Writer)
		return
	}

	origin := strings.ToLower(filepath.Ext(file.Filename))

	color.Green("Filename : %v", origin)

	filename := fmt.Sprintf("uploaded/%v%v", id, origin)
	output := fmt.Sprintf("uploaded/%v", id)

	fmt.Println(filename)
	fmt.Println(output)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		color.Red("File Save Error In Server in below api...")
		middlewares.ErrorResponse("File Save Error In Server...", c.Writer)
		return
	}
	defer f.Close()

	if err := c.SaveUploadedFile(file, f.Name()); err != nil {
		middlewares.ErrorResponse("File Save Error.", c.Writer)
		return
	}

	// Call the MapTiler CLI using the exec.Command function
	//maptiler-engine -raster -o output -zoom 1 5 -watermark_opacity 100 images/test.png
	// cmd := exec.Command("maptiler-engine ", "-raster ", "-o ", "uploaded/newone ", "-zoom 1 5 ", "uploaded/upload.png")
	cmd := exec.Command("rm", "-r", output)
	cmd.Start()

	if origin == ".wmf" {
		if _, err = os.Stat(filename); err == nil {
			os.Remove(filename)
		}

		font.UpdateDraw2dFontSettings()

		var img_data *image.RGBA
		img_data, err = wmf.LoadAsImage(filename, true)
		if err != nil {
			fmt.Println("Failed to load file")
			middlewares.ErrorResponse("Failed to load file.", c.Writer)
			return
		}

		dstFilename := fmt.Sprintf("%v.tiff", id)
		err = saveImageAsTiff(dstFilename, img_data)
		if err != nil {
			fmt.Println("Failed to save image")
			middlewares.ErrorResponse("Failed to save image.", c.Writer)
			return
		}

		middlewares.SuccessMessageResponse("WMF success", c.Writer)
	}

	filename = fmt.Sprintf("uploaded/%v.tiff", id)

	cmd = exec.Command("C:\\Program Files\\MapTiler Engine\\maptiler-engine.exe", "-raster", "-o", output, "-zoom", "1", "5", filename)
	// Capture the command's output and error streams
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error capturing output: %v\n", err)
		middlewares.ErrorResponse("File Save Error In Server.", c.Writer)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error capturing error: %v\n", err)
		middlewares.ErrorResponse("File Save Error In Server.", c.Writer)
		return
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		middlewares.ErrorResponse("File Save Error In Server.", c.Writer)
		return
	}

	// Read the output stream
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Read the error stream
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Wait for the command to finish
	err = cmd.Wait()

	// Check for errors
	if err != nil {
		fmt.Printf("Error generating tiles: %v\n", err)
		middlewares.ErrorResponse("Error generating tiles:...", c.Writer)
		os.Exit(1)
	}

	middlewares.SuccessMessageResponse("Uploaded Successfully", c.Writer)
}

// UploadFileEndpoint -> upload file
func TestUploadFileEndpoint(c *gin.Context) {
	file, err := c.FormFile("file")

	id := c.PostForm("id")
	// fileName := c.PostForm("file_name")
	if err != nil {
		color.Red("File Upload Failed in below api...")
		middlewares.ErrorResponse("File Upload Failed.", c.Writer)
		return
	}
	origin := strings.ToLower(filepath.Ext(file.Filename))

	color.Green("Filename : %v", origin)

	filename := fmt.Sprintf("uploaded/%v.%v", id, origin)
	output := fmt.Sprintf("uploaded/%v", id)

	fmt.Println(filename)
	fmt.Println(output)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		color.Red("File Save Error In Server in below api...")
		middlewares.ErrorResponse("File Save Error In Server...", c.Writer)
		return
	}
	defer f.Close()

	if err := c.SaveUploadedFile(file, f.Name()); err != nil {
		middlewares.ErrorResponse("File Save Error.", c.Writer)
		return
	}

	middlewares.SuccessMessageResponse("Uploaded Successfully", c.Writer)
}

var GetSiteInfoByID = gin.HandlerFunc(func(c *gin.Context) {

	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		color.Red("Getting Param Issue in below api...")
		middlewares.ErrorResponse("Invalid Site ID", c.Writer)
		return
	}

	var result models.SiteInfo

	var site models.Site
	collection := client.Database("sitemanagerservice").Collection("sites")
	err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&site)
	if err != nil {
		color.Red("No Site Exist with ID: %v in below api...", id)
		middlewares.ErrorResponse("Site does not exist", c.Writer)
		return
	}

	result.ID = site.ID.Hex()
	result.Level = site.Level

	for {

		color.Red("%v", site.Level)
		color.Red("%v", site.ID.Hex())
		color.Red("%v", site.Name)
		color.Red("%v\n", site.ParentID)
		switch site.Level {
		case 0:
			result.SiteGroupID = site.ID.Hex()
			result.SiteGroupName = site.Name
		case 1:
			result.SiteID = site.ID.Hex()
			result.SiteName = site.Name
		case 2:
			result.MapID = site.ID.Hex()
			result.MapName = site.Name
		case 3:
			result.ZoneGroupID = site.ID.Hex()
			result.ZoneGroupName = site.Name
		case 4:
			result.ZoneID = site.ID.Hex()
			result.ZoneName = site.Name
		default:
		}
		if site.Level == 0 {
			break
		}

		collection := client.Database("sitemanagerservice").Collection("sites")

		objID, err = primitive.ObjectIDFromHex(site.ParentID)
		if err != nil {
			color.Red("Getting Param Issue in below api...")
			middlewares.ErrorResponse("Invalid Site ID", c.Writer)
			return
		}
		err = collection.FindOne(context.Background(), bson.D{primitive.E{Key: "_id", Value: objID}}).Decode(&site)
		if err != nil {
			color.Red("No Site Exist with ID: %v in below api...", site.ParentID)
			middlewares.ErrorResponse("Site does not exist", c.Writer)
			return
		}
	}

	middlewares.SuccessOneRespond(result, "SiteInfo", c.Writer)
})

func saveImageAsTiff(fileName string, img *image.RGBA) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

	if err == nil {
		// var buf bytes.Buffer
		options := tiff.Options{Compression: tiff.Deflate}
		err = tiff.Encode(f, img, &options)
		err = f.Close()
	}

	return err
}
