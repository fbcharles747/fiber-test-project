package main

import (
	"log"

	"github.com/fbcharles747/fiber-api/database"
	"github.com/fbcharles747/fiber-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type SomeStruct struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
	Msg  string `json:"msg`
}

func welcome(c *fiber.Ctx) error {
	data := SomeStruct{
		Msg:  "Welcom to my api",
		Name: "Charles",
		Age:  21,
	}
	return c.JSON(data)
}

func setupRoutes(app *fiber.App) {
	app.Get("/api", welcome)

	app.Post("/api/users", routes.CreateUser)

	app.Get("/api/users", routes.GetUsers)

	app.Get("api/user/:id", routes.GetUser)

	app.Put("/api/user/:id", routes.UpdateUser)

	app.Post("/api/devices/csv", routes.BulkUploadCSV)

	app.Post("/api/devices/json", routes.BulkUploadJSON)

}

func main() {
	database.ConnectDb()
	app := fiber.New()

	app.Use(cors.New())

	setupRoutes(app)

	log.Fatal(app.Listen(":3001"))
}
