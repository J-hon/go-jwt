package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectId `bson:"_id"`
	First_name    string             `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     string             `json:"last_name" validate:"required,min=2,max=100"`
	Password      string             `json:"password" validate:"required,min=6"`
	Email         string             `json:"email" validate:"email,required"`
	Mobile_number string             `json:"mobile_number" validate:"required"`
	User_type     string             `json:"mobile_number" validate:"required, eq=ADMIN|eq=USER"`
	Token         string             `json:"token"`
	Refresh_token string             `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
}
