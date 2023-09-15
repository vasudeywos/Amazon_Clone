package controllers

import (
	"mymod/models"
	"mymod/utils"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type InputJSON struct {
    Address string `json:"address" binding:"required"`
}

func (s *Server) PlaceOrder(c *gin.Context) {
    var OrderItList []models.OrderItems
    var Order models.Order
    var prodlist []models.Product
	var InJSON InputJSON

    //JSON parse
    if err := c.ShouldBindJSON(&InJSON); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, _ := utils.CurrentUser(c)

    session := sessions.Default(c)

    cart, _ := GetCart(session, user)

    for _, cartit := range cart.Items {
        var orderit models.OrderItems
		orderit.Address=InJSON.Address
        orderit.Quantity = cartit.Quantity

        var prod models.Product
        if errs := s.db.Model(&models.Product{}).Where("name=?", cartit.ProductName).Find(&prod).Error; errs != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": errs.Error()})
            return
        }

        orderit.BoughtProduct = prod

        OrderItList = append(OrderItList, orderit)      //Add to Order
        q:=(prod.Quantity-cartit.Quantity)              //Sync Quantity
        prod.Quantity=q                                 
        cart.DeleteItem(cartit)                         //Sync Cart
        if errd :=SaveCart(session, cart); errd != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user's cart"})
            return
        }
    
        c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart"})
        prodlist=append(prodlist, prod)
    }

    Order.UserID = user.ID  //Foreign Key putting UserId
    Order.BoughtItems = OrderItList
    for _,prdcit:=range prodlist{
        if err := s.db.Model(&models.Product{}).Where("id = ?", prdcit.ID).Updates(&prdcit).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
            return
        }
    
        c.JSON(http.StatusOK,prdcit)
    }
    if err := s.db.Create(&Order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Order"})
        return
    }

    c.JSON(http.StatusOK, Order)
}
