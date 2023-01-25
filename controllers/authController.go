package controllers

import (
	"encoding/json"
	"fmt"
	"goLang/database"
	"goLang/entities"
	"goLang/helper"
	"goLang/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

func Register(c echo.Context) error {
	user := new(entities.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	var result int
	database.Instance.Select("id").Where("username = ?").Find(&result)
	if result != 0 {
		resp := c.JSON(http.StatusInternalServerError, helper.ErrorLog(http.StatusInternalServerError, "Username existed", "EXT_REF"))
		return resp
	}
	//hashing password
	hash, error := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if error != nil {
		resp := c.JSON(http.StatusInternalServerError, helper.ErrorLog(http.StatusInternalServerError, " Error While Hashing Password", "EXT_REF"))
		return resp
	}
	user.Password = string(hash)
	user.DateCreated = helper.DateTime()
	user.Role = "member"
	database.Instance.Create(&user)
	return c.JSON(http.StatusCreated, user)
}
func Login(c echo.Context) error {
	loginUser := new(models.Login)
	if err := c.Bind(loginUser); err != nil {
		return err
	}
	var result models.User
	var user entities.User
	var permission []entities.Permission
	database.Instance.Where("username = ?", loginUser.Username).Find(&user)
	if user.ID == 0 {
		resp := c.JSON(http.StatusInternalServerError, helper.ErrorLog(http.StatusInternalServerError, "username not found", "EXT_REF"))
		return resp
	}

	error := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if error != nil {
		log.Error("Invalid Password :", error)
		resp := c.JSON(http.StatusInternalServerError, helper.ErrorLog(http.StatusInternalServerError, "Invalid Password", "EXT_REF"))
		return resp
	}
	//permission list
	database.Instance.Where("user_id = ?", user.ID).Find(&permission)

	//addUserToken
	AccessToken, RefreshToken, permissionToken := AddUserToken(user, permission)
	result.AccessToken = AccessToken
	result.RefreshToken = RefreshToken
	result.Permission = permissionToken
	//resp
	resp := c.JSON(http.StatusOK, result)
	log.Info()
	return resp
}
func AddUserToken(user entities.User, permission []entities.Permission) (string, string, string) {
	var result models.User
	expirationTime := time.Now().Add(1 * time.Minute)
	result.AccessToken = helper.JwtGenerator(user.Username, user.Firstname, user.Lastname, "secret", user.Role, user.ID, expirationTime)
	refreshTokenExpirationTime := time.Now().Add(5 * time.Minute)
	result.RefreshToken = helper.JwtGenerator(user.Username, user.Firstname, user.Lastname, "secret", user.Role, user.ID, refreshTokenExpirationTime)
	//addUserToken
	userToken := &entities.UserToken{UserId: &user.ID, Token: result.AccessToken, RefreshToken: result.RefreshToken}
	result.Permission = helper.PermissionJwtGenerator(permission, "secret", expirationTime)
	database.Instance.Create(&userToken)
	return result.AccessToken, result.RefreshToken, result.Permission
}
func GetDynamicPermission(c echo.Context) error {
	data, err := json.MarshalIndent(c.Echo().Routes(), "", "  ")
	if err != nil {
		return err
	}
	resp := c.JSON(http.StatusOK, string(data))
	return resp
}
func GetUserDynamicPermission(c echo.Context) error {
	userId := c.Param("id")
	var permission []entities.Permission
	database.Instance.Where("user_id = ?", userId).Find(&permission)
	data, err := json.MarshalIndent(c.Echo().Routes(), "", "  ")
	if err != nil {
		return err
	}
	var userPermission models.UserPermission
	userPermission.Permission = string(data)
	userPermission.UserPermission = permission

	resp := c.JSON(http.StatusOK, userPermission)
	return resp
}
func PostDynamicPermission(c echo.Context) error {
	var permission entities.Permission
	model := new(models.Permission)
	if err := c.Bind(model); err != nil {
		return err
	}

	permission.UserId = model.UserId
	permission.Name = model.Name
	permission.Path = model.Path
	database.Instance.Save(&permission)
	return c.JSON(http.StatusOK, "true")
}
func RemoveDynamicPermission(c echo.Context) error {
	var permission entities.Permission
	model := new(models.Permission)
	if err := c.Bind(model); err != nil {
		return err
	}
	permission.UserId = model.UserId
	permission.Name = model.Name
	permission.Path = model.Path
	permission.ID = model.ID
	database.Instance.Delete(&permission, model.ID)
	return c.JSON(http.StatusOK, "true")
}

func Logout(c echo.Context) error {
	refreshToken := c.QueryParam("refreshToken")
	fmt.Println(refreshToken)
	var userTokens []entities.UserToken
	database.Instance.Where("refresh_token = ?", refreshToken).Find(&userTokens)
	database.Instance.Delete(userTokens)
	return c.JSON(http.StatusOK, "true")
}
func Users(c echo.Context) error {
	var users []entities.User
	database.Instance.Find(&users)
	return c.JSON(http.StatusOK, users)
}
func ChangePassword(c echo.Context) error {
	var user entities.User
	model := new(models.ChangePassword)
	if err := c.Bind(model); err != nil {
		return err
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	username := claims["name"].(string)
	database.Instance.Where("username = ?", username).Find(&user)
	if user.ID == 0 {
		resp := c.JSON(http.StatusInternalServerError, helper.ErrorLog(http.StatusInternalServerError, "username not found", "EXT_REF"))
		return resp
	}
	error := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(model.OldPassword))
	if error != nil {
		log.Error("Invalid Password :", error)
		resp := c.JSON(http.StatusInternalServerError, helper.ErrorLog(http.StatusInternalServerError, "Invalid Password", "EXT_REF"))
		return resp
	}
	hash, error := bcrypt.GenerateFromPassword([]byte(model.NewPassword), 5)
	if error != nil {
		resp := c.JSON(http.StatusInternalServerError, helper.ErrorLog(http.StatusInternalServerError, " Error While Hashing Password", "EXT_REF"))
		return resp
	}
	user.Password = string(hash)
	database.Instance.Save(&user)
	return c.JSON(http.StatusOK, "true")
}
func RefreshToken(c echo.Context) error {
	model := new(models.RefreshToken)
	if err := c.Bind(model); err != nil {
		return err
	}

	var user entities.User
	var result models.User
	userToken := new(entities.UserToken)
	database.Instance.Where("refresh_token = ?", model.RefreshToken).Find(&userToken)
	if userToken.ID == 0 {
		resp := c.JSON(http.StatusInternalServerError, helper.ErrorLog(http.StatusInternalServerError, "not found", "EXT_REF"))
		return resp
	}
	database.Instance.First(&user, userToken.UserId)
	database.Instance.Delete(&userToken)
	var permission []entities.Permission
	database.Instance.Where("user_id = ?", user.ID).Find(&permission)
	AccessToken, RefreshToken, PermissionToken := AddUserToken(user, permission)
	result.AccessToken = AccessToken
	result.RefreshToken = RefreshToken
	result.Permission = PermissionToken
	//resp
	resp := c.JSON(http.StatusOK, result)
	return resp
}
