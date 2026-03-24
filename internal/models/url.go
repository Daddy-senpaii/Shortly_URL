package models

import (
    "go.mongodb.org/mongo-driver/v2/bson"
    "time"
)

type ShortURL struct {
    ID          bson.ObjectID       `bson:"_id,omitempty"  json:"id,omitempty"`
    ShortCode   string              `bson:"short_code"     json:"short_code"`
    ShortURL    string              `bson:"short_url"      json:"short_url"`      
    OriginalURL string              `bson:"original_url"   json:"url"`
    UserID      bson.ObjectID       `bson:"user_id"        json:"user_id,omitempty"`
    ClickCount  int64               `bson:"click_count"    json:"click_count"`
    CreatedAt   time.Time           `bson: "created_at"      json:"create_at"`
}
