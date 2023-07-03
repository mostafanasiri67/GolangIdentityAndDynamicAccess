package main

import (
	"encoding/json"
	"fmt"
	"goLang/controllers"
	"goLang/database"
	"goLang/models"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	// Load Configurations from config.json using Viper
	LoadAppConfig()

	// Initialize Database
	database.Connect(AppConfig.ConnectionString)
	database.Migrate()

	// Register Routes
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"authorization", "Content-Type"},
		AllowCredentials: true,
		AllowMethods:     []string{echo.OPTIONS, echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	RegisterProductRoutes(e)

	e.Logger.Fatal(e.Start(AppConfig.Port))
}
func SkipperFn(skipURLs []string) func(echo.Context) bool {
	return func(context echo.Context) bool {
		for _, url := range skipURLs {
			if url == context.Request().URL.String() {
				return true
			}
		}
		return false
	}
}
func RegisterProductRoutes(e *echo.Echo) {
	e.POST("/api/auth/signup", controllers.Register).Name = "signup"
	e.POST("/api/auth/signin", controllers.Login).Name = "signin"
	e.POST("/api/auth/refreshToken", controllers.RefreshToken, isLoggedIn).Name = "refreshToken"
	e.GET("/api/auth/logout", controllers.Logout).Name = "logout"
	e.POST("/api/ChangePassword", controllers.ChangePassword, isLoggedIn).Name = "ChangePassword"
	e.GET("/api/ApiSettings", getConfig).Name = "ApiSettings"
	e.GET("/api/Users", controllers.Users, isLoggedIn, checkRole).Name = "Users"
	e.GET("/api/DynamicPermission", controllers.GetDynamicPermission, isLoggedIn, checkRole).Name = "DynamicPermission"
	e.GET("/api/UserDynamicPermission/:id", controllers.GetUserDynamicPermission, isLoggedIn, checkRole).Name = "UserDynamicPermission"
	e.POST("/api/DynamicPermission", controllers.PostDynamicPermission, isLoggedIn, checkRole).Name = "DynamicPermissionPost"
	e.POST("/api/RemoveDynamicPermission", controllers.RemoveDynamicPermission, isLoggedIn, checkRole).Name = "RemoveDynamicPermission"

	// adminGroup := e.Group("api/admin", isLoggedIn, checkRole)
	// adminGroup.GET("/Users", controllers.Users).Name = "adminUsers"
	data, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(data))
}
func getConfig(c echo.Context) error {
	apiConfig := &models.ApiConfig{AdminRoleName: "admin", LoginPath: "auth/signin", RegisterPath: "auth/signup", LogoutPath: "auth/logout", AccessTokenObjectKey: "AccessToken", RefreshTokenObjectKey: "RefreshToken", RefreshTokenPath: "auth/refreshToken"}
	return c.JSON(http.StatusOK, apiConfig)
}
