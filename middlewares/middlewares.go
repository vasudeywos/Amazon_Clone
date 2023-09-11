package middlewares

import (
    "fmt"
    "mymod/utils"
    "github.com/gin-gonic/gin"
    "net/http"
)

func JwtAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {


        err := utils.ValidateToken(c)

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"Unauthorized":"Authentication required"})
            fmt.Println(err)
            c.Abort()
            return
        }
        c.Next()
    }
}
