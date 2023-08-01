package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"tracio.com/alertservice/controllers"
	middlewares "tracio.com/alertservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/alertservice/test", middlewares.IsAuthorized(controllers.Test))

}
