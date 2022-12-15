package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/J-hon/go-jwt/helpers"
	"github.com/J-hon/go-jwt/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

		count, err = userCollection.CountDocuments(context.Background(), bson.M{"mobile_number": user.Mobile_number})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking mobile number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Mobile number already exists"})
			return
		}

		password := helpers.HashPassword(*user.Password)
		user.Password = &password
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Id = primitive.NewObjectID()
		user.User_id = user.Id.Hex()
		token, refresh_token, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)

		user.Token = &token
		user.Refresh_token = &refresh_token

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

		err := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "User with such email doesn't exist"})
			return
		}

		isValid, msg := helpers.VerifyPassword(*user.Password, *foundUser.Password)
		if !isValid {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		helpers.UpdateTokens(token, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(context.Background(), bson.M{"_id": foundUser.Id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)
	}
}
