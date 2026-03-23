package models

import (
    "go.mongodb.org/mongo-drvier/v2/bson"
)

type User struct {
    ID bson.ObjectID    `bson:"_id, omitempty"  json:"id,omitempty"`
    Email string        `bson:"email"           json:"email" validate: "required, email"`
    Password string     `bson:"password"        json:"-"`
}
