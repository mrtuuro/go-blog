package main

import (
	"blog/database"
	"blog/router"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal(err) // TODO------------->>>>>  Buraya CUSTOM EXCEPTION uygulaması ekle !!
	}

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":" + port))
}
