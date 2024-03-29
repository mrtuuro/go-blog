package handler

import (
	"blog/database"
	"blog/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// GetAllArticles Tüm kayıtları getir.
func GetAllArticles(ctx *fiber.Ctx) error {
	// tüm kayıtları cursor olarak alıyoruz
	query := bson.D{{}}
	cursor, err := database.Mg.Db.Collection("articles").Find(ctx.Context(), query)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error()) //TODO restful api'ye göre düzenle
	}

	var articles []models.Article = make([]models.Article, 0)
	// All iterates the cursor and decodes each document into results.
	if err := cursor.All(ctx.Context(), &articles); err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.JSON(articles)
}

// GetSingleArticle id'ye göre kayıt getir
func GetSingleArticle(ctx *fiber.Ctx) error {
	// get id by params
	params := ctx.Params("id")

	// ObjectIDFromHex creates a new ObjectID from a hex string.
	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// D is an ordered representation of a BSON document. This type should be used when the order of the elements matters,
	// such as MongoDB command documents. If the order of the elements does not matter, an M should be used instead.
	filter := bson.D{{"_id", _id}}

	var result models.Article

	// Decode will unmarshal the document represented by this SingleResult into v. If there was an error from the operation
	// that created this SingleResult, that error will be returned. If the operation returned no documents, Decode will
	// return ErrNoDocuments.
	if err := database.Mg.Db.Collection("articles").FindOne(ctx.Context(), filter).Decode(&result); err != nil {
		return ctx.Status(http.StatusNotFound).SendString("Not found") // TODO err'ün içine bakıp status code'a karar ver
	}

	return ctx.Status(http.StatusOK).JSON(result)
}

// CreateArticle Inserting a new article to DB
func CreateArticle(ctx *fiber.Ctx) error {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	authorID := claims["id"].(string)
	collection := database.Mg.Db.Collection("articles")

	article := new(models.Article)
	if err := ctx.BodyParser(article); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// ID'yi mongodb kendisi oluştursun.
	article.ID = ""
	article.Author = authorID

	// Insert the record to DB
	insertResult, err := collection.InsertOne(ctx.Context(), article)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Insert yapılan document'i ID'sine gire filtreleyip geri döndürüyoruz.
	// The _id of the inserted document. A value generated by the driver will be of type primitive.ObjectID.
	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	createdRecord := collection.FindOne(ctx.Context(), filter) // Az önce yarattığımız article'ı değişkende tuttuk.

	createdArticle := &models.Article{}                          //
	if err := createdRecord.Decode(createdArticle); err != nil { // Az önce yarattığımız article'ı Decode ile unmarshal'layıp createdArticle'a atadık.
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.Status(http.StatusCreated).JSON(createdArticle)
}

// UpdateArticle Update an article
func UpdateArticle(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	// ObjectIDFromHex creates a new ObjectID from a hex string.
	articleID, err := primitive.ObjectIDFromHex(idParam)

	// Girilen ID'nin valid olup olmadığını kontrol ediyoruz.
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Boş article oluşturduk. Request'in body'sinden gelen data'yı boş article'e yazdırdık.
	article := new(models.Article)
	if err := ctx.BodyParser(article); err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Parametreden gelen ID ile query oluşturup ilgili data'yı buluyoruz.
	query := bson.D{{Key: "_id", Value: articleID}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "name", Value: article.Name},
				{Key: "author", Value: article.Author},
				{Key: "rating", Value: article.Rating},
				{Key: "content", Value: article.Content},
			},
		},
	}
	err = database.Mg.Db.Collection("articles").FindOneAndUpdate(ctx.Context(), query, update).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments { // ErrNoDocuments means that the filter did not match any documents.
			return ctx.SendStatus(http.StatusBadRequest)
		}
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	// Return updated article
	article.ID = idParam
	return ctx.Status(http.StatusOK).JSON(article)
}

// Search contents endpoint
func Search(ctx *fiber.Ctx) error {

	filter := bson.M{}

	if q := ctx.Query("q"); q != "" {
		filter = bson.M{
			"content": bson.M{
				"$regex": primitive.Regex{
					Pattern: q,
				},
			},
		}
	}

	cursor, err := database.Mg.Db.Collection("articles").Find(ctx.Context(), filter)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(http.StatusOK).JSON("")
	}

	var articles []models.Article = make([]models.Article, 0)
	if err := cursor.All(ctx.Context(), &articles); err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	//return ctx.SendString(ctx.Query("q"))
	return ctx.Status(http.StatusOK).JSON(articles)
}

// DeleteArticle Delete an article
func DeleteArticle(ctx *fiber.Ctx) error {
	articleID, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Verilen ID ile db'den kayıtı bul ve sil.
	query := bson.D{{Key: "_id", Value: articleID}}
	result, err := database.Mg.Db.Collection("articles").DeleteOne(ctx.Context(), &query)
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if result.DeletedCount < 1 { // Article bulunamamıştır.
		return ctx.Status(http.StatusNotFound).SendString("Article not found.")
	}
	return ctx.Status(http.StatusNoContent).SendString("Article deleted.")
}
