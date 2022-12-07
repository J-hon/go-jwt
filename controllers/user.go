package controllers

import (
	"github.com/J-hon/go-jwt/database"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection("users")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func SignUp()

func SignIn()

func GetUsers()

func GetUser()
