package handler

import (
	"blog/database"
	"blog/models"
	"blog/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// Login
func Login(ctx *fiber.Ctx) error {
	collection := database.Mg.Db.Collection("users")
	user := new(models.LoginUser)
	dbUser := new(models.DbUser)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	filter := bson.D{{Key: "email", Value: user.Email}}

	if err := collection.FindOne(ctx.Context(), filter).Decode(dbUser); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)
	if passErr != nil {
		return ctx.Status(http.StatusBadRequest).SendString("Password and email does not match!")
	}

	jwtToken, err := utils.GenerateJWT(dbUser)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	return ctx.Status(http.StatusOK).JSON(jwtToken)
}

// Register
func Register(ctx *fiber.Ctx) error {
	collection := database.Mg.Db.Collection("users")
	user := new(models.CreateUser)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	user.Password = utils.HashPassword([]byte(user.Password))
	insertedUser, err := collection.InsertOne(ctx.Context(), user)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	filter := bson.D{{Key: "_id", Value: insertedUser.InsertedID}}
	createdUser := collection.FindOne(ctx.Context(), filter)
	returnUser := &models.User{}
	if err := createdUser.Decode(returnUser); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	return ctx.Status(http.StatusOK).JSON(returnUser)
}

// Create user
func CreateUser(ctx *fiber.Ctx) error {

	authErr := utils.VerifyAuthentication(ctx)
	if authErr != nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"err_code": "authError", "status_code": fiber.StatusUnauthorized, "err_message": authErr.Error()})
	}

	collection := database.Mg.Db.Collection("users")

	user := new(models.CreateUser)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	insertResult, err := collection.InsertOne(ctx.Context(), user)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	createdUser := collection.FindOne(ctx.Context(), filter)
	fmt.Println(createdUser)
	returnUser := &models.User{}
	//returnUser.ID = insertResult.InsertedID

	if err := createdUser.Decode(returnUser); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	return ctx.Status(http.StatusCreated).JSON(returnUser)
}

// Delete user
func DeleteUser(ctx *fiber.Ctx) error {

	authErr := utils.VerifyAuthentication(ctx)
	if authErr != nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"err_code": "authError", "status_code": fiber.StatusUnauthorized, "err_message": authErr.Error()})
	}

	params := ctx.Params("id")

	userID, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	filter := bson.D{{"_id", userID}}

	var user models.User
	if err := database.Mg.Db.Collection("users").FindOneAndDelete(ctx.Context(), filter).Decode(&user); err != nil {
		return ctx.SendStatus(http.StatusNotFound)
	}
	return ctx.Status(http.StatusNoContent).SendString("User deleted.")
}

// Update user
func UpdateUser(ctx *fiber.Ctx) error {

	authErr := utils.VerifyAuthentication(ctx)
	if authErr != nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"err_code": "authError", "status_code": fiber.StatusUnauthorized, "err_message": authErr.Error()})
	}

	params := ctx.Params("id")
	userID, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	user := new(models.CreateUser)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	filter := bson.D{{"_id", userID}}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "name", Value: user.Name},
				{Key: "email", Value: user.Email},
				{Key: "password", Value: user.Password},
			},
		}}
	var returnUser models.User
	if err := database.Mg.Db.Collection("users").FindOneAndUpdate(ctx.Context(), filter, update).Decode(&returnUser); err != nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.Status(http.StatusOK).JSON(returnUser)

}

// Get user by ID
func GetUserById(ctx *fiber.Ctx) error {

	authErr := utils.VerifyAuthentication(ctx)
	if authErr != nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"err_code": "authError", "status_code": fiber.StatusUnauthorized, "err_message": authErr.Error()})
	}

	params := ctx.Params("id")
	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	filter := bson.D{{"_id", _id}}

	var returnResult models.User
	if err := database.Mg.Db.Collection("users").FindOne(ctx.Context(), filter).Decode(&returnResult); err != nil {
		return ctx.SendStatus(http.StatusNotFound)
	}

	return ctx.Status(http.StatusOK).JSON(returnResult)
}

// Get all users
func GetAllUsers(ctx *fiber.Ctx) error {
	authErr := utils.VerifyAuthentication(ctx)
	if authErr != nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"err_code": "authError", "status_code": fiber.StatusUnauthorized, "err_message": authErr.Error()})
	}

	query := bson.D{{}}
	cursor, err := database.Mg.Db.Collection("users").Find(ctx.Context(), query)
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	var users []models.User = make([]models.User, 0)
	if err := cursor.All(ctx.Context(), &users); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	return ctx.Status(http.StatusOK).JSON(users)

}
