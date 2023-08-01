package controllers

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"

	"tracio.com/localliveservice/db"
	middlewares "tracio.com/localliveservice/handlers"
)

var client = db.Dbconnect()

func Test(c *gin.Context) {
	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}
