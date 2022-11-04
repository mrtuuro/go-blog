package utils

import (
	"blog/database"
	"blog/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

var JwtKey = []byte("login")

type Claims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

func HashPassword(password []byte) string {

	hashedPassword, err := bcrypt.GenerateFromPassword(password, 8)
	fmt.Println(string(hashedPassword))
	if err != nil {
		log.Fatal(err)
	}

	return string(hashedPassword)
}

func GenerateJWT(user *models.DbUser) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(72 * time.Hour).Unix()
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["jti"] = uuid.New()
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyAuthentication(ctx *fiber.Ctx) error {

	// TokenClaim struct'ında ID olmasının sebebi: Register olduktan sonra oluşturulan
	// token'ın payload'ında ID, email, name bulunuyor. Bu dataları kullanıp db'de sorgu yapıyoruz.
	type TokenClaim struct {
		ID string `json:"id"`
		jwt.RegisteredClaims
	}
	// Header'dan Authorization'ı alıp "Bearer " kısmını atıyoruz. Elimizde sadece token bulunuyor.
	reqToken := ctx.Get("Authorization")
	reqToken = strings.Split(reqToken, "Bearer ")[1]

	// Token'da bulunan datayı almak için ParseWithClaims() fonksiyonunu çağırıyoruz.
	token, err := jwt.ParseWithClaims(reqToken, &TokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("login"), nil
	})
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	// token'ın claim'lerini(payload'da bulunan data) aldık.
	claims, ok := token.Claims.(*TokenClaim)
	if !ok && !token.Valid {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	userID, _ := primitive.ObjectIDFromHex(claims.ID) // Payload'daki ID'yi ObjectID'ye çevirdik.
	filter := bson.D{{Key: "_id", Value: userID}}
	var returnResult models.User

	if err := database.Mg.Db.Collection("users").FindOne(ctx.Context(), filter).Decode(&returnResult); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"err_code": "authError", "status_code": fiber.StatusUnauthorized, "err_message": err.Error()})
	}
	fmt.Println(returnResult)

	return nil
}
