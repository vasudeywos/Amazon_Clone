package controllers

import (
	"mymod/models"
	"mymod/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    IsStaff  bool   `json:"is_staff"`
}

type LoginInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type Server struct {
    db *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
    return &Server{db: db}
}

func (s *Server) Register(c *gin.Context) {
    var input RegisterInput

    if err := c.ShouldBind(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := models.User{Username: input.Username, Password: input.Password,IsStaff: input.IsStaff}
    user.HashPassword()

    if err := s.db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    }

    c.JSON(http.StatusCreated , gin.H{"message": "User created"})
}

//THis is a mtehod of server
func (s *Server)LoginCheck(username, password string) (string, error) {
    var err error

    user := models.User{}

    if err = s.db.Model(models.User{}).Where("username=?", username).Take(&user).Error; err != nil {
        return "", err
    }

    err = models.VerifyPassword(password, user.Password)

    if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
        return "", err
    }

    token, err := utils.GenerateToken(user)

    if err != nil {
        return "", err
    }

    return token, nil

}

func (s *Server) Login(c *gin.Context) {
    //Get email and password from user
    var LogInt LoginInput

    if err := c.ShouldBind(&LogInt); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }    
    //Find requested user
    user:=models.User{Username:LogInt.Username,Password: LogInt.Password}

    token,err:=s.LoginCheck(user.Username,user.Password)//THis is a mtehod of server
    //compare the credentials
    if err!=nil{
        c.JSON(http.StatusBadRequest,gin.H{
            "error":"Wrong credentials",
        })
    }

    //return it
    c.JSON(http.StatusOK, gin.H{"token": token})
}
