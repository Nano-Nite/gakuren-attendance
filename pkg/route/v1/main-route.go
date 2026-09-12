package v1

import (
	"log"
	"os"
	"strings"

	"gakuren-system.com/pkg/helper"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes() {
	app := fiber.New()

	app.Use(func(c fiber.Ctx) error {
		log.Printf("API hit : %s %s <> IP Address : %s <> User Agent : %s\n", c.Method(), c.OriginalURL(), c.IP(), c.UserAgent())

		return c.Next()
	})

	app.Get("/"+strings.Trim(helper.API_VERSION, "/")+"/attendance/health", func(c fiber.Ctx) error {
		return c.SendString("ready to go !!!!!!!!!!")
	})

	// child route setup
	SetupLocationRoutes(app, helper.API_VERSION)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
