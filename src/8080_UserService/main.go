package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"github.com/helios/go-sdk/sdk"
	middlewares "tracio.com/userservice/handlers"
	"tracio.com/userservice/routes"
)

func main() {

	port := middlewares.DotEnvVariable("PORT")

	sdk.Initialize("UserService", "eb243b5a8bd81d5e7fa4", sdk.WithEnvironment("local_env"))

	fmt.Println("Port number is: " + port)
	color.Cyan("🌏 Server running on localhost:" + port)

	router := gin.Default()

	// router.Use(customLogger())
	routes.Routes(router)

	err := router.RunTLS(":"+port, "../certificate/cert.pem", "../certificate/key.pem")
	if err != nil {
		log.Fatal(err)
	}
}
