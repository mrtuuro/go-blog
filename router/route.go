package router

import (
	"blog/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	// Burada middleware ekleyebiliriz.
	v1 := app.Group("/v1")

	// Route'larımız burada
	v1.Get("/article", handler.GetAllArticles)
	v1.Get("/article/:id", handler.GetSingleArticle)
	v1.Post("/article", handler.CreateArticle)
	v1.Put("/article/:id", handler.UpdateArticle)
	v1.Get("/search", handler.Search)
	v1.Delete("/article/:id", handler.DeleteArticle)

}
