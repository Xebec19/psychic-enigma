package main

import (
	"github.com/Xebec19/psychic-enigma/internal"
	"github.com/gofiber/fiber/v2"
)

func main() {

	internal.NewAWS()

	app := fiber.New()

	app.Post("/upload", func(c *fiber.Ctx) error {

		form, err := c.MultipartForm()

		if err != nil {
			return err
		}

		ch := internal.UploadImage(form.File["image"][0])

		response := <-ch
		if response.Err != nil {
			return response.Err
		}

		c.SendString(response.Value)
		return nil
	})

	app.Listen(":3000")
}
