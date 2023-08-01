package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
	"tracio.com/sitemanagerservice/controllers"
	middlewares "tracio.com/sitemanagerservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/sitemanagerservice/test", middlewares.IsAuthorized(controllers.Test))

	router.GET("/api/v1/sitemanagerservice/sites", middlewares.IsAuthorized(controllers.GetSites))

	router.POST("/api/v1/sitemanagerservice/addsite", middlewares.IsAuthorized(controllers.AddNewSite))

	router.DELETE("/api/v1/sitemanagerservice/deletesite/:id", middlewares.IsAuthorized(controllers.DeleteSite))

	router.PUT("/api/v1/sitemanagerservice/setscale/:id", middlewares.IsAuthorized(controllers.EditSiteByID))
	router.PUT("/api/v1/sitemanagerservice/setunit/:id", middlewares.IsAuthorized(controllers.EditSiteByID))

	router.GET("/api/v1/sitemanagerservice/sitebyid/:id", middlewares.IsAuthorized(controllers.GetSiteByID))

	router.POST("/api/v1/sitemanagerservice/upload", controllers.UploadFileEndpoint)
	router.POST("/api/v1/sitemanagerservice/testupload", controllers.TestUploadFileEndpoint)

	router.GET("/api/v1/sitemanagerservice/maptile/:id", middlewares.IsAuthorized(controllers.ServeMapTile))
	router.StaticFS("/api/v1/sitemanagerservice/static/", http.Dir("./uploaded/"))

	router.GET("/api/v1/sitemanagerservice/siteinfo/:id", middlewares.IsAuthorized(controllers.GetSiteInfoByID))
}
