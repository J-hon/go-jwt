package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/J-hon/go-jwt/database"
	"github.com/J-hon/go-jwt/helpers"
	"github.com/J-hon/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/bson"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(context.Background(), bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Email already exists"})
		}

		mobileCount, mobileError := userCollection.CountDocuments(context.Background(), bson.M{"mobile_number": user.Mobile_number})
		if mobileError != nil {
			log.Panic(mobileError)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking mobile number"})
			return
		}

		if mobileCount > 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Mobile number already exists"})
		}
	}
}

func SignIn()

func GetUsers()

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err := userCollection.FindOne(context.Background(), bson.M{"user_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
