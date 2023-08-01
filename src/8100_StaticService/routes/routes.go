package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
	"tracio.com/staticservice/controllers"
	middlewares "tracio.com/staticservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/staticservice/test", middlewares.IsAuthorized(controllers.Test))

	router.GET("/api/v1/staticservice/objecttypes", middlewares.IsAuthorized(controllers.GetObjectTypes))
	router.GET("/api/v1/staticservice/devicetypes", middlewares.IsAuthorized(controllers.GetDeviceTypes))

	router.POST("/api/v1/staticservice/objecticonupload", middlewares.IsAuthorized(controllers.UploadObjectIcon))
	router.StaticFS("api/v1/staticservice/objecticon/", http.Dir("./uploaded/objecticons/"))
}
