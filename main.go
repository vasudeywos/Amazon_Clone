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
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
)
func init(){        //Runs before main
	initializers.LoadEnvVariables()
}

func DbInit() (*gorm.DB, error) {
    db, err := models.Setup()
    if err != nil {
        log.Printf("Error setting up database: %v\n", err)
        return nil, err
    }
    return db, nil
}

func main(){

	router:=gin.Default()

	db, err := DbInit()
    if err != nil {
        log.Fatalf("Error initializing database: %v\n", err)
    }
    log.Printf("Value of db: %v\n", db)

    store := gormsessions.NewStore(db, true, []byte("secret"))

	server:=controllers.NewServer(db,store)

	//API
	router.POST("/register",server.Register)
	router.POST("/login",server.Login)

	if err := godotenv.Load(); err != nil {
        log.Println("Error loading .env file")
    }

	logged:=router.Group("/loggedin")

	logged.Use(middlewares.JwtAuthMiddleware(),sessions.Sessions("moisessions",store))
    logged.GET("/getprod", server.GetAllProd)
	logged.GET("/searchprodname", server.FindProd)
	logged.GET("/filtcateg", server.FilterCatg)
    logged.POST("/sellprod", server.SellProd)
	logged.GET("/sellerprodlist", server.SellerProd)
	logged.POST("/updateprod",server.UpdateProd)
	logged.POST("/deleteprod",server.DeleteProd)
	logged.POST("/addtocart", server.AddToCart)
	logged.GET("/getcart", server.GetCart)
	logged.POST("/orderplace",server.PlaceOrder)

    //port := os.Getenv("PORT")
	router.Run("localhost:8080")

}