package models

import (
    "go.mongodb.org/mongo-driver/v2/bson"
)

type ShortURL struct {
    ID          bson.ObjectID       `bson:"_id,omitempty"  json:"id,omitempty"`
    ShortCode   string              `bson:"short_code"     json:"short-code"`
    OriginalURL string              `bson:"original_url"   json:"original_url"`
    UserID      bson.ObjectID       `bson:"user_id"        json:"user_id,omitempty"`
    ClickCount  int64               `bson:"click_count"    json:"click_count"`
    CreatedAt   time.Time           `bson: created_at      json:"create_at"`
}
