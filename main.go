package main

import (
	"blog/database"
	"blog/router"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal(err) // TODO------------->>>>>  Buraya CUSTOM EXCEPTION uygulaması ekle !!
	}

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
