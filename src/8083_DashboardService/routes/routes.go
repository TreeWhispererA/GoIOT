package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"tracio.com/dashboardservice/controllers"
	middlewares "tracio.com/dashboardservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/dashboardservice/test", middlewares.IsAuthorized(controllers.Test))

}
