package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/Daddy-senpaii/Shorty_URL/internal/models"
    "github.com/Daddy-senpaii/Shorty_URL/internal/utils"
    "github.com/Daddy-senpaii/Shorty_URL/internal/config"
    "go.mongodb.org/mongo-driver/v2/bson"
    //"go.mongodb.org/mongo-driver/v2/mongo"
    "context"
    "errors"
    "log"
    "net/url"
    "time"
    "net/http"
    "fmt"
)

func ValidateURL(original_url string) bool {
    parsed, err := url.ParseRequestURI(original_url)
    return err == nil && parsed.Scheme != "" && parsed.Scheme == ""
}

func GenerateShortURL(original_url, shortCode string) (string, error) {
    parser, err := url.Parse(original_url)
    if err != nil {
        return "" , errors.New("We have error ")
    }
    
    return parser.Scheme + "://" + parser.Host + "/" + shortCode, nil
}



func PostURL(c *gin.Context){
    fmt.Println("Handler Hit man")
    var url_struct models.ShortURL

    if err := c.ShouldBindJSON(&url_struct); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Having some internal problems",
        })
        return 
    }

    
    original_url := url_struct.OriginalURL

    if original_url == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "url required",
        })
        return
    }

    if ValidateURL(original_url){
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Provide correct url",
        })

        return
    }
    url_struct.ID = bson.NewObjectID()
    url_struct.OriginalURL = original_url

    short_code, err := utils.GenerateShortCode(8)
    if err != nil {
        log.Fatal(err)
    }

    url_struct.ShortCode = short_code

    shortURL, err := GenerateShortURL(original_url, short_code)
    if err != nil {
        log.Fatal(err)
    }

    url_struct.ShortURL = shortURL
    url_struct.CreatedAt = time.Now()

   // fmt.Println(url_struct)

   collection := config.GetCollection("url")
   ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
   defer cancel()

   _, err = collection.InsertOne(ctx, url_struct)
   if err != nil {
       c.JSON(http.StatusInternalServerError, gin.H {
           "error": "Could not connect to database",
       })
       return
   }

   c.JSON(http.StatusCreated, url_struct)
}
