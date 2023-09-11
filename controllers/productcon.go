package controllers

import (
	"mymod/models"
	"mymod/utils"
	"net/http"
	"github.com/gin-gonic/gin"
)

type NewProdIn struct {
    Name string `json:"name" binding:"required"`
	Price string `json:"price" binding:"required"`
	Quantity int `json:"quantity" binding:"required"`
}

func (s *Server)SellProd (c *gin.Context){
		var Produc models.Product
	
		if err := c.ShouldBindJSON(&Produc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		user, err := utils.CurrentUser(c)
		if err != nil {
	
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		if !user.IsStaff {
			//     // Only staff members can create products
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			return
			}
	
		if err := s.db.Create(&Produc).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	
		c.JSON(http.StatusCreated, Produc)
	}
	


func(s *Server)GetProd (c *gin.Context){
	var prodls []models.Product

	if err := s.db.Find(&prodls).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

	c.JSON(http.StatusOK,prodls)
}