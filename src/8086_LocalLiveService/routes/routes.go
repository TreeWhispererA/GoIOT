package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"tracio.com/localliveservice/controllers"
	middlewares "tracio.com/localliveservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/localliveservice/test", middlewares.IsAuthorized(controllers.Test))

}
