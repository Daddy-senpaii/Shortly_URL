package main
import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Daddy-senpaii/Shorty_URL/internal/config"
    "github.com/Daddy-senpaii/Shorty_URL/internal/utils"
    "github.com/Daddy-senpaii/Shorty_URL/internal/controller"
)

func getPing(context *gin.Context){
    context.JSON(http.StatusOK, gin.H{
        "message": "pong connected",
    })
}


func main(){
    config.MakeConnection()
    defer config.Client.Disconnect(context.Background())

    router := gin.Default()
    admin := router.Group("/api")

    admin.GET("/ping", getPing)

    admin.GET("/test-code", func(c *gin.Context){
        code, err := utils.GenerateShortCode(8)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"short-code":code})
    })

    admin.POST("/post-url", controller.PostURL)
    router.GET("/:code", controller.GetURL)
    router.Run(":5000")
}
