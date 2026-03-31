package controller

import (
    "fmt"
    "context"
    "time"
    "unicode"
    "strings"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/v2/bson"
    "github.com/golang-jwt/jwt/v5"
    "github.com/Daddy-senpaii/Shorty_URL/internal/models"
    "github.com/Daddy-senpaii/Shorty_URL/internal/config"
    "net/http"
    "net/mail"
)

var JwtSecret = []byte("secret-key")    

func IsValidEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil 
}

func IsValidPassword(password string) bool {
    if len(password) < 8 {
        return false
    }
    var hasUpper, hasLower, hasNumber, hasSpecial bool

    for _ , ch := range password {
        switch {
            case unicode.IsUpper(ch):
                hasUpper = true
            case unicode.IsLower(ch):
                hasLower = true
            case unicode.IsDigit(ch):
                hasNumber = true
            case strings.ContainsRune("@$!%*?&", ch):
                hasSpecial = true
        }
    }
    return hasUpper && hasLower && hasNumber && hasSpecial
}

func HashPassword(password string )(string, error){
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func Register(c *gin.Context){
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H {"error": "Invalid Format"})
        return
    }

    fmt.Printf("Password raw: %#v\n", user.Password)

    if !IsValidEmail(user.Email){
        c.JSON(http.StatusBadRequest, gin.H{"error": "Write correct email man"})
        return
    }
    if !IsValidPassword(user.Password){
        c.JSON(http.StatusBadRequest, gin.H{"error": "Follow Proper Password Format"})
        return
    }

    hashPassword, err := HashPassword(user.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Hashing failed"})
        return
    }
    user.Password = hashPassword
    user.ID = bson.NewObjectID()

    //// write into the database make collection in db and put the details in there
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := config.GetCollection("users")
    _, err = collection.InsertOne(ctx, user)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H {"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func GenerateJWT(userID, email string)(string, error){
    claims := jwt.MapClaims{
        "user_id": userID,
        "email": email,
        "exp": time.Now().Add(24*time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JwtSecret)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

func LogIn(c *gin.Context){
    var user models.User
    var dbUser models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Details"})
        return 
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel();

    collection := config.GetCollection("users")
    err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User Not Found"})
        return
    }
    /// Password Validation
    err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password"})
        return
    }

    // now generate JWT
    token, err := GenerateJWT(dbUser.ID.Hex(), dbUser.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Login Successful",
        "token": token,
        "user": gin.H{
            "id": dbUser.ID,
            "email": dbUser.Email,
        },
    })

}
