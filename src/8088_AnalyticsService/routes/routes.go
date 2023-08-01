package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"tracio.com/analyticsservice/controllers"
	middlewares "tracio.com/analyticsservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/analyticsservice/test", middlewares.IsAuthorized(controllers.Test))

}
