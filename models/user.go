package models

import (

    "gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
    "html"
	"strings"
    "log"
    "errors"

)

type User struct {
    gorm.Model
    Username string `gorm:"size:255" json:"username" binding:"required"`
    Password string `gorm:"size:255" json:"password" binding:"required"`
    IsStaff  bool   `json:"is_staff"`
    Products []Product //Many to One
    CartID uint
    Orders []Order //Has Many;Many to One
}

func (user *User) HashPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)
	
    user.Username = html.EscapeString(strings.TrimSpace(user.Username))

    return nil
}

func VerifyPassword(password, hashedPassword string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetUserById(uid uint) (User, error) {
    var user User

    db, err := Setup()

    if err != nil {
        log.Println(err)
        return User{}, err
    }
    if err := db.Preload("Products").Where("id=?", uid).Find(&user).Error; err != nil {
        return user, errors.New("user not found")

    }

    return user, nil
}