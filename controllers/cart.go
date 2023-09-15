package controllers

import (
	"mymod/models"
	"mymod/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (s *Server) AddToCart(c *gin.Context) {
	var input CartItem
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, erra := utils.CurrentUser(c)
	if erra != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": erra.Error()})
	return
	}
	session := sessions.Default(c)
	cart, errb :=GetCartForUser(session,user)
	if errb != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errb.Error()})
		return
	}
	cart.UserID=user.ID
	//Check if Quantity is There
	var producs models.Product
	if errs := s.db.Model(&models.Product{}).Where("name=?", input.ProductName).Find(&producs).Error; errs != nil {
	c.JSON(http.StatusNotFound, gin.H{"error": errs.Error()})
	return
	}
	if producs.Quantity < input.Quantity{
		c.JSON(http.StatusNotFound, gin.H{"error": "Amount of Product added is more than that available for sale"})
	}
	cart.AddItem(input)
	if errd :=SaveCart(session, cart); errd != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user's cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})
}

func (s *Server) GetCart(c *gin.Context) {
	var cart *Cart
	session := sessions.Default(c)
	user, err := utils.CurrentUser(c)
	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
	}
	
	cart,_=GetCart(session,user)
	// if err!= nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found"})
	// 	return
	// }

	c.JSON(http.StatusOK, cart)
}