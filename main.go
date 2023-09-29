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

	store := gormsessions.NewStore(db, true, []byte("secret"))

	server:=controllers.NewServer(db,store)

	//API
	router.POST("/register",server.Register)
	router.POST("/login",server.Login)

	if err := godotenv.Load(); err != nil {
        log.Println("Error loading .env file")
    }

	router.MaxMultipartMemory = 8 << 20  // 8 MiB limit

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
	logged.POST("/uploadimage",controllers.FileUpload)
	logged.Static("/viewimage", "assets")


    //port := os.Getenv("PORT")
	//router.Run("localhost:8080")
	//router.Run(":80")
	
    certPath := "mycert.crt"
    keyPath := "mycert.key"

    // Run the server with HTTPS
    err := http.ListenAndServeTLS(":443", certPath, keyPath, router)
    if err != nil {
        panic(err)
    
	}
}