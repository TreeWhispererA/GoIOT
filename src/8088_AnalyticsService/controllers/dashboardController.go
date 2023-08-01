package controllers

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"

	"tracio.com/analyticsservice/db"
	middlewares "tracio.com/analyticsservice/handlers"
)

var client = db.Dbconnect()

func Test(c *gin.Context) {
	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}
