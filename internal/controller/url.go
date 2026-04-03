package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/Daddy-senpaii/Shorty_URL/internal/models"
    "github.com/Daddy-senpaii/Shorty_URL/internal/utils"
    "github.com/Daddy-senpaii/Shorty_URL/internal/config"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
    //"github.com/golang-jwt/jwt/v5"
    "context"
    "errors"
    "log"
    "net/url"
    "time"
    "net/http"
    "fmt"
)


func GetURL(c *gin.Context){
    fmt.Println("Handling Get handler")
    code := c.Param("code")
    fmt.Println(code)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var url models.ShortURL

    collection := config.GetCollection("url")
    
    if err := collection.FindOne(ctx, bson.M{"short_code": code}).Decode(&url); err != nil {

        if errors.Is(err, mongo.ErrNoDocuments){
            c.JSON(http.StatusNotFound, gin.H{"error": "urlnot found"})
            return 
        } else if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
            return 
        }
    }

    go func(){
        ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel2()

        collection.UpdateOne(
                            ctx2,
                            bson.M{"short_code": code},
                            bson.M{"$inc": bson.M{"click_count": 1}},
                            )
    }()
    c.Redirect(http.StatusMovedPermanently, url.OriginalURL)


}

// get my url

func GetMyURL(c *gin.Context){
    id,exists := c.Get("user_id")
    if !exists{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user doesn't exist"})
        return 
    }
    userID, ok := id.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid format"})
        return
    }

    objectID, err := bson.ObjectIDFromHex(userID)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid format"})
        return 
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := config.GetCollection("url")

    cursor, err := collection.Find(ctx, bson.M{"user_id": objectID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "error in fetching urls"})
        return 
    }

    defer cursor.Close(ctx)

    var urls []models.ShortURL

    for cursor.Next(ctx){
        var url models.ShortURL
        if err := cursor.Decode(&url); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error in decoding"})
            return
        }

        urls = append(urls, url)
    }

    if err := cursor.Err(); err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": "no urls found"})
        return 
    }
    if len(urls) == 0 {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "no urls found"})
        return 
    }

    c.JSON(http.StatusOK, urls)
}

// get single links

func GetLinkById(c *gin.Context){
//    fmt.Println("hit getlinkbyid")
    id := c.Param("id")
  //  fmt.Println("id is: ",id)
    linkID, err := bson.ObjectIDFromHex(id)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
        return
    }

    // Get user_id from JWT middleware
    userIDRaw, exists := c.Get("user_id")
    if !exists{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User Id is not found"})
        return
    }

    userIDStr, ok := userIDRaw.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid id format"})
        return
    }
    ObjectID, err := bson.ObjectIDFromHex(userIDStr)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid id format"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var link models.ShortURL

    err = config.GetCollection("url").FindOne(ctx, bson.M{"_id": linkID, "user_id": ObjectID}).Decode(&link)

    if err == mongo.ErrNoDocuments{
        c.JSON(http.StatusNotFound, gin.H{"error": "Link not found man"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch"})
        return
    }
    c.JSON(http.StatusOK, link)
}


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

    //validation of url

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
    // Generating Short Code 
    short_code, err := utils.GenerateShortCode(8)
    if err != nil {
        log.Fatal(err)
    }

    url_struct.ShortCode = short_code

    shortURL, err := GenerateShortURL(original_url, short_code)
    if err != nil {
        log.Fatal(err)
    }
    // extracting url from jwt context
    userID, exists := c.Get("user_id")
    if !exists{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not found"})
        return
    }
    uidStr, ok := userID.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id format"})
        return
    }

    objectID, err := bson.ObjectIDFromHex(uidStr)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id format"})
        return
    }

    url_struct.ShortURL = shortURL
    url_struct.CreatedAt = time.Now()
    url_struct.UserID = objectID

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
