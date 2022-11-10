package router

import (
	"blog/handler"
	"blog/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// Burada middleware ekleyebiliriz.
	v1 := app.Group("/v1", middleware.CheckForAuth())
	auth := app.Group("/auth")

	auth.Post("/register", handler.Register) // rate limit her zaman olmalı
	auth.Post("/login", handler.Login)

	// Route'larımız burada
	v1.Get("/article", handler.GetAllArticles)
	v1.Get("/article/:id", handler.GetSingleArticle)
	v1.Post("/article", handler.CreateArticle)
	v1.Put("/article/:id", handler.UpdateArticle)
	v1.Get("/search", handler.Search)
	v1.Delete("/article/:id", handler.DeleteArticle)

	v1.Post("/article/:articleID/comment", handler.PostComment)

	v1.Post("/user", handler.CreateUser)
	v1.Get("/users", handler.GetAllUsers)
	v1.Get("/user/:id", handler.GetUserById)
	v1.Delete("/user/:id", handler.DeleteUser)
	v1.Put("/user/:id", handler.UpdateUser)
}
