package main

import (
	"log"
	"mymod/controllers"
	"mymod/initializers"
	"mymod/models"
	"mymod/middlewares"
	"net/http"
	_"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)
func init(){        //Runs before main
	initializers.LoadEnvVariables()
}

func DbInit() *gorm.DB {
    db, err := models.Setup()
    if err != nil {
        log.Println("Problem setting up database")
    }
    return db
}
func main(){
	router:=gin.Default()

	db:=DbInit()

	router.GET("/ping",func(c *gin.Context){
		c.JSON(http.StatusOK,gin.H{
			"body":"This is it",
		})
	})

	server:=controllers.NewServer(db)

	//API
	router.POST("/register",server.Register)
	router.POST("/login",server.Login)

	if err := godotenv.Load(); err != nil {
        log.Println("Error loading .env file")
    }

	logged:=router.Group("/loggedin")

	logged.Use(middlewares.JwtAuthMiddleware())
    logged.GET("/getprod", server.GetProd)
    logged.POST("/sellprod", server.SellProd)

    //port := os.Getenv("PORT")
	router.Run("localhost:8080")

}