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
	Category string `json:"category"`
}
//STAFF
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
		Produc.UserID=user.ID //ForeignKey putting UserID
		if err := s.db.Create(&Produc).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.Products=append(user.Products, Produc)
		if err := s.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update userprodlist"})
			return
		}
		c.JSON(http.StatusCreated, Produc)
	}

//STAFF
func (s *Server)UpdateProd (c *gin.Context) {
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
	//Only Creator can edit
	if user.ID!=Produc.ID{
		c.JSON(http.StatusForbidden, gin.H{"error": "Cant Access"})
		return
	}
	//Update requires Primary Key
	if err := s.db.Model(&models.Product{}).Where("id = ?", Produc.ID).Updates(&Produc).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
        return
    }

    c.JSON(http.StatusOK, Produc)
}
//STAFF
//Get Product That You Sold
func(s *Server)SellerProd (c *gin.Context){
	var sellPrdLst []models.Product
	user,_:=utils.CurrentUser(c)
	if !user.IsStaff {
		//     // Only staff members can create products
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
		}
	
	list:=user.Products
	for _,prod:=range list{
		sellPrdLst = append(sellPrdLst, prod)
	}
	c.JSON(http.StatusOK,sellPrdLst)

}
//STAFF
//Delete product
func(s *Server)DeleteProd (c *gin.Context){
	var Produc models.Product
		
	if erre := c.ShouldBindJSON(&Produc); erre != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": erre.Error()})
		return
		}
		
	user, err := utils.CurrentUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Only Creator can edit
	if user.ID!=Produc.ID{
		c.JSON(http.StatusForbidden, gin.H{"error": "Cant Access"})
		return
	}
	delerr:=s.db.Delete(&models.Product{},Produc.ID).Error
	if delerr!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
        return
	}
	c.JSON(http.StatusOK,gin.H{"message":"Deleted Product"})
}

//Get All Producs
func(s *Server)GetAllProd (c *gin.Context){
	var prodls []models.Product

	err := s.db.Model(&models.Product{}).Preload("Users").Find(&prodls).Error
	if err!=nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
	}

	c.JSON(http.StatusOK,prodls)
}

//Search Products by NAme
func (s *Server)FindProd (c *gin.Context){
	var srchprodls []models.Product

	queryterm:=c.Query("name")
	if queryterm == "" {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
		return
	}

	err := s.db.Model(&models.Product{}).Where("name=?", queryterm).Find(&srchprodls).Error
	if err!=nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
	}

	c.JSON(http.StatusOK, srchprodls)
}

//Filter by category
func (s *Server)FilterCatg (c *gin.Context){
	var srchprodls []models.Product

	queryterm:=c.Query("category")
	if queryterm == "" {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
		return
	}

	err := s.db.Model(&models.Product{}).Where("category=?", queryterm).Find(&srchprodls).Error
	if err!=nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
	}

	c.JSON(http.StatusOK, srchprodls)
}
