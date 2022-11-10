package handler

import (
	"blog/database"
	"blog/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func PostComment(ctx *fiber.Ctx) error {

	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	commentOwner := claims["name"]

	collection := database.Mg.Db.Collection("comments")
	idParam := ctx.Params("articleID")
	articleID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	/*
		articleFilter := bson.D{{Key: "_id", Value: articleID}}
		article := new(models.Article)
		if err := database.Mg.Db.Collection("articles").FindOne(ctx.Context(), articleFilter).Decode(&article); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		fmt.Println(article)
	*/

	comment := new(models.CreateComment)
	if err := ctx.BodyParser(comment); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	comment.Owner = commentOwner.(string)
	comment.Date = time.Now()
	insertResult, _ := collection.InsertOne(ctx.Context(), comment)
	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	createdComment := collection.FindOne(ctx.Context(), filter)

	returnComment := &models.Comment{}
	if err := createdComment.Decode(returnComment); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	articleFilter := bson.D{{Key: "_id", Value: articleID}}
	article := new(models.Article)
	update := bson.D{
		{
			Key: "$push",
			Value: bson.D{
				{Key: "comments", Value: returnComment},
			},
		}}

	if err := database.Mg.Db.Collection("articles").FindOneAndUpdate(ctx.Context(), articleFilter, update).Decode(&article); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.Status(fiber.StatusOK).JSON(&article)
}
