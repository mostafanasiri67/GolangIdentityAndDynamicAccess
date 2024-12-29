package main

import (
	"goLang/database"
	"goLang/entities"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var isLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("secret"),
})

func checkRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := claims["role"].(string)
		userId := claims["userId"].(float64)
		var permission []entities.Permission
		database.Instance.Where("user_id = ?", userId).Find(&permission)
		if role == "admin" {
			return next(c)
		}
		hasPermission := false
		for _, element := range permission {
			if c.Path() == element.Path {
				hasPermission = true
				break
			}
		}
		if hasPermission {
			return next(c)
		} else {
			return echo.ErrForbidden
		}

	}
}
