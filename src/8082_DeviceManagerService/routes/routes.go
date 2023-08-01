package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"tracio.com/devicemanagerservice/controllers"
	middlewares "tracio.com/devicemanagerservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/devicemanagerservice/test", middlewares.IsAuthorized(controllers.Test))

	router.GET("/api/v1/devicemanagerservice/manufacturers", middlewares.IsAuthorized(controllers.GetManufacturers))
	router.GET("/api/v1/devicemanagerservice/models", middlewares.IsAuthorized(controllers.GetModels))
	router.GET("/api/v1/devicemanagerservice/models/:id", middlewares.IsAuthorized(controllers.GetModelsByManufacturers))

	router.GET("/api/v1/devicemanagerservice/devices", middlewares.IsAuthorized(controllers.GetDevices))
	router.GET("/api/v1/devicemanagerservice/briefprfiddevices", middlewares.IsAuthorized(controllers.GetBriefDevices))
	router.POST("/api/v1/devicemanagerservice/adddevice", middlewares.IsAuthorized(controllers.AddNewDevice))
	router.PUT("/api/v1/devicemanagerservice/editdevice/:id", middlewares.IsAuthorized(controllers.EditDevice))
	router.DELETE("/api/v1/devicemanagerservice/deletedevice/:id", middlewares.IsAuthorized(controllers.DeleteDevice))
	router.GET("/api/v1/devicemanagerservice/devicebyid/:id", middlewares.IsAuthorized(controllers.GetDeviceByID))

	router.GET("/api/v1/devicemanagerservice/templates", middlewares.IsAuthorized(controllers.GetTemplates))
	router.POST("/api/v1/devicemanagerservice/addtemplate", middlewares.IsAuthorized(controllers.AddNewTemplate))
	router.PUT("/api/v1/devicemanagerservice/edittemplate/:id", middlewares.IsAuthorized(controllers.EditTemplate))
	router.DELETE("/api/v1/devicemanagerservice/deletetemplate/:id", middlewares.IsAuthorized(controllers.DeleteTemplate))
	router.GET("/api/v1/devicemanagerservice/templatebyid/:id", middlewares.IsAuthorized(controllers.GetTemplateByID))

	router.GET("/api/v1/devicemanagerservice/devicetemplates", middlewares.IsAuthorized(controllers.GetDeviceTemplates))
	router.POST("/api/v1/devicemanagerservice/adddevicetemplate", middlewares.IsAuthorized(controllers.AddNewDeviceTemplate))
	router.PUT("/api/v1/devicemanagerservice/editdevicetemplate/:id", middlewares.IsAuthorized(controllers.EditDeviceTemplate))
	router.DELETE("/api/v1/devicemanagerservice/deletedevicetemplate/:id", middlewares.IsAuthorized(controllers.DeleteDeviceTemplate))
	router.GET("/api/v1/devicemanagerservice/devicetemplatebyid/:id", middlewares.IsAuthorized(controllers.GetDeviceTemplateByID))

	router.POST("api/v1/devicemanagerservice/importjson", middlewares.IsAuthorized(controllers.ImportJsonDeviceData))

	router.POST("api/v1/devicemanagerservice/bulkaction", middlewares.IsAuthorized(controllers.BulkAction))

}
