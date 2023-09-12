package utils

import (
	"errors"
	"fmt"
	"mymod/models"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(user models.User) (string,error){
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"exp": time.Now().Add(time.Hour*24*30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	// tokenString, err := token.SignedString(hmacSampleSecret)

	// fmt.Println(tokenString, err)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}

func getTokenFromRequest(c *gin.Context) string {
    bearerToken := c.Request.Header.Get("Authorization")

    splitToken := strings.Split(bearerToken, " ")
    if len(splitToken) == 2 {
        return splitToken[1]
    }
    return ""
}


func GetToken(c *gin.Context) (*jwt.Token, error) {
    tokenString := getTokenFromRequest(c)
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("SECRET")), nil   //Have to get token in []byte()format
    })
    return token, err
}


func ValidateToken (c *gin.Context) error {
    token, err := GetToken(c)

    if err != nil {
        return err
    }

    _, ok := token.Claims.(jwt.MapClaims)
    if ok && token.Valid {
        return nil
    }

    return errors.New("invalid token provided")
}

func GetUserIDFromContext(c *gin.Context) (uint, error) {
    // Get the JWT token
    token, err := GetToken(c)
    if err != nil {
        return 0, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return 0, errors.New("invalid token")
    }

    // Extract the user ID from claims
    userID := uint(claims["id"].(float64))

    return userID, nil
}


func CurrentUser(c *gin.Context) (models.User, error) {
    err := ValidateToken(c)
    if err != nil {
        return models.User{}, err
    }
    token, _ := GetToken(c)
    claims, _ := token.Claims.(jwt.MapClaims)
    userId := uint(claims["id"].(float64))

    user, err := models.GetUserById(userId)
    if err != nil {
        return models.User{}, err
    }
    return user, nil
}
