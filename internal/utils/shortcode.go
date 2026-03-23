package utils

import  (
    "crypto/rand"
    "encoding/base64"
)


func GenerateShortCode(length int)(string, error){

    if length < 6 {
        length = 8
    }

    b := make([]byte, length)
    _, err := rand.Read(b)

    if err != nil {
        return "", err
    }

    code := base64.RawURLEncoding.EncodeToString(b)

    if len(code) > length{
        code = code[:length]
    }

    return code , nil
}
