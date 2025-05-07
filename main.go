package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// JWT Secret Key
	secretKey := os.Getenv("SECRET_KEY")
	// Login route
	app.Post("/login", login(secretKey))
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(secretKey),
	}))

	// CORS middleware
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*", // Adjust this to be more restrictive if needed
	// 	AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// }))

	// Setup routes
	app.Get("/book", getBooks)
	app.Get("/book/:id", getBook)
	app.Post("/book", createBook)
	app.Put("/book/:id", updateBook)
	app.Delete("/book/:id", deleteBook)

	// Setup route
	app.Post("/upload", uploadImage)

	if secretKey == "" {
		log.Fatal("SECRET_KEY not set in environment")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}
	app.Listen(":" + port)

}

func uploadImage(c *fiber.Ctx) error {
	// Read file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Save the file to the server
	os.MkdirAll("./uploads", os.ModePerm)
	err = c.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File uploaded successfully: " + file.Filename)
}

// Dummy user for example
var user = struct {
	Email    string
	Password string
}{
	Email:    "admin@example.com",
	Password: "admin1234",
}

func login(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var request LoginRequest
		if err := c.BodyParser(&request); err != nil {
			return err
		}

		// Check credentials - In real world, you should check against a database
		if request.Email != user.Email || request.Password != user.Password {
			return fiber.ErrUnauthorized
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "John Doe"
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token
		t, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
}
