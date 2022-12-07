package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/J-hon/go-jwt/database"
	"github.com/J-hon/go-jwt/helpers"
	"github.com/J-hon/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

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

		user.Password = helpers.HashPassword(user.Password)
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Id = primitive.NewObjectID()
		token, refresh_token, _ := helpers.GenerateAllTokens(&user.Email, &user.First_name, &user.Last_name, &user.User_type, &user.Id)

		user.Token = token
		user.Refresh_token = refresh_token

		res, err := userCollection.InsertOne(context.Background(), user)
		if err != nil {
			msg := fmt.Sprintf("Error occurred")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user, foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User with such email doesn't exist"})
			return
		}

		isValid, msg := helpers.VerifyPassword(user.Password, foundUser.Password)
	}
}

func GetUsers()

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err := userCollection.FindOne(context.Background(), bson.M{"_id": userId}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
