package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"tracio.com/reportservice/controllers"
	middlewares "tracio.com/reportservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/reportservice/test", middlewares.IsAuthorized(controllers.Test))

	router.GET("/api/v1/reportservice/objects", middlewares.IsAuthorized(controllers.GetObjects))
	router.POST("/api/v1/reportservice/addobject", middlewares.IsAuthorized(controllers.AddNewObject))
	router.PUT("/api/v1/reportservice/editobject/:id", middlewares.IsAuthorized(controllers.EditObject))
	router.DELETE("/api/v1/reportservice/deleteobject/:id", middlewares.IsAuthorized(controllers.DeleteObject))
	router.GET("/api/v1/reportservice/objectbyid/:id", middlewares.IsAuthorized(controllers.GetObjectByID))

	router.GET("/api/v1/reportservice/objectreport", middlewares.IsAuthorized(controllers.GetObjectReport))
	router.GET("/api/v1/reportservice/tagreport", middlewares.IsAuthorized(controllers.GetTagReport))
}
