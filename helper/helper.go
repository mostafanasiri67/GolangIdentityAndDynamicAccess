package helper

import (
	"encoding/json"
	"goLang/entities"
	"goLang/models"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func DateTime() string {
	currentTime := time.Now()
	result := currentTime.Format("2006-01-02 15:04:05") //yyyy-mm-dd HH:mm:ss
	return result
}

type Claims struct {
	UserId      uint   `json:"userId"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	DisplayName string `json:"displayName"`
	jwt.StandardClaims
}

func JwtGenerator(username, firstname, lastname, key string, role string, userId uint, expirationTime time.Time) string {
	claims := &Claims{
		Name:        username,
		Role:        role,
		DisplayName: firstname + " " + lastname,
		UserId:      userId,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string.
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return err.Error()
	}

	return tokenString
}
func PermissionJwtGenerator(permission []entities.Permission, key string, expirationTime time.Time) string {
	stringPermission, err := json.Marshal(permission)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}
	claims := &Claims{
		Name: string(stringPermission),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string.
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return err.Error()
	}

	return tokenString
}
func ErrorLog(rc int, detail, ext_ref string) models.Error {
	var error models.Error
	error.ResponseCode = rc
	error.Message = "Failed"
	error.Detail = detail
	error.ExternalReference = ext_ref

	return error
}
