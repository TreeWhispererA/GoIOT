package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	"github.com/helios/go-sdk/sdk"
	middlewares "tracio.com/devicemanagerservice/handlers"
	"tracio.com/devicemanagerservice/routes"
)

func main() {

	port := middlewares.DotEnvVariable("PORT")

	sdk.Initialize("DeviceManagerService", "eb243b5a8bd81d5e7fa4", sdk.WithEnvironment("local_env"))

	fmt.Println("Port number is: " + port)
	color.Cyan("🌏 Server running on localhost:" + port)

	router := gin.Default()
	router.Use(gin.Logger())
	routes.Routes(router)

	err := router.RunTLS(":"+port, "../certificate/cert.pem", "../certificate/key.pem")
	if err != nil {
		log.Fatal(err)
	}
}
