package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"tracio.com/ruleservice/controllers"
	middlewares "tracio.com/ruleservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/ruleservice/test", middlewares.IsAuthorized(controllers.Test))

}
