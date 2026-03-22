package main

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Daddy-senpaii/Shorty_URL/internal/config"
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
    router.Run(":5000")
}
