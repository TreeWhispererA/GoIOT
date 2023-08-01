package routes

import (
	gin "github.com/helios/go-sdk/proxy-libs/heliosgin"
	http "github.com/helios/go-sdk/proxy-libs/helioshttp"
	"tracio.com/userservice/controllers"
	middlewares "tracio.com/userservice/handlers"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/api/v1/userservice/test", middlewares.IsAuthorized(controllers.Test))

	router.POST("/api/v1/userservice/login", controllers.UserLogin)
	router.POST("/api/v1/userservice/register", controllers.UserRegister)
	router.GET("/api/v1/userservice/generateotp", middlewares.IsAuthorized(controllers.GenerateOTP))
	router.POST("/api/v1/userservice/verifyotp", middlewares.IsAuthorized(controllers.VerifyOTP))

	router.GET("/api/v1/userservice/roles", middlewares.IsAuthorized(controllers.GetRoles))
	router.POST("/api/v1/userservice/smartroles", middlewares.IsAuthorized(controllers.GetSmartRoles))
	router.POST("/api/v1/userservice/smartusers", middlewares.IsAuthorized(controllers.GetSmartUsers))

	router.GET("/api/v1/userservice/users", middlewares.IsAuthorized(controllers.GetUsers))

	router.POST("/api/v1/userservice/adduser", middlewares.IsAuthorized(controllers.AddNewUser))
	router.POST("/api/v1/userservice/addrole", middlewares.IsAuthorized(controllers.AddNewRole))

	router.DELETE("/api/v1/userservice/deleteuser/:id", middlewares.IsAuthorized(controllers.DeleteUser))
	router.DELETE("/api/v1/userservice/deleterole/:id", middlewares.IsAuthorized(controllers.DeleteRole))

	router.PUT("/api/v1/userservice/edituser/:id", middlewares.IsAuthorized(controllers.EditUser))
	router.PUT("/api/v1/userservice/editrole/:id", middlewares.IsAuthorized(controllers.EditRole))
	router.POST("/api/v1/userservice/changepassword", middlewares.IsAuthorized(controllers.ChangePassword))

	router.GET("/api/v1/userservice/usersbyrole/:id", middlewares.IsAuthorized(controllers.GetUsersByRole))
	router.GET("/api/v1/userservice/userbyid/:id", middlewares.IsAuthorized(controllers.GetUserByID))
	router.GET("/api/v1/userservice/getinfo", middlewares.IsAuthorized(controllers.GetInfo))
	router.GET("/api/v1/userservice/rolebyid/:id", middlewares.IsAuthorized(controllers.GetRoleByID))

	router.POST("api/v1/userservice/upload", middlewares.IsAuthorized(controllers.UploadFileEndpoint))

	router.StaticFS("api/v1/userservice/photo/", http.Dir("./uploaded/photo/"))
	router.StaticFS("api/v1/userservice/totp/", http.Dir("./uploaded/totp/"))
}
