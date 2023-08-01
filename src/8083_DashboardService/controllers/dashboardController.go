package controllers

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"

	"tracio.com/dashboardservice/db"
	middlewares "tracio.com/dashboardservice/handlers"
)

var client = db.Dbconnect()

func Test(c *gin.Context) {
	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}
