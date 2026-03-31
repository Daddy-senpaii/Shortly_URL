package middleware

import (
    "net/http"
    "strings"
    "github.com/Daddy-senpaii/Shorty_URL/internal/controller"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare() gin.HandlerFunc{
    return func(c *gin.Context){
        authHeader := c.GetHeader("Authorization")
        if authHeader == ""{
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) !=2 || parts[0] != "Bearer"{
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
            c.Abort()
            return
        }
        tokenString := parts[1]

        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error){
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(controller.JwtSecret) , nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok{
            c.Set("user_id", claims["user_id"])
            c.Set("email", claims["email"])
        }

        c.Next()
    }
}
