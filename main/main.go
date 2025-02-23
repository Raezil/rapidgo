package main

import (
	"fmt"
	"os"

	"github.com/rwiteshbera/rapidgo"
)

func main() {
	rapidgo.LoadEnv()

	router := rapidgo.New()

	fmt.Println(os.Getenv("PORT"))
	fmt.Println(os.Getenv("PASS"))
	fmt.Println(os.Getenv("USERNAME"))
	router.SetDebug(false)

	// Middleware examples
	router.Use(rapidgo.SampleMiddleware1())
	router.Use(rapidgo.SampleMiddleware2())
	router.Use(rapidgo.SampleMiddleware3())

	// 🔹 Static Routes
	router.Get("/", func(c *rapidgo.Context) {
		c.JSON(200, map[string]string{"message": "Home Page"})
	})
	router.Get("/about", func(c *rapidgo.Context) {
		c.JSON(200, map[string]string{"message": "About Page"})
	})

	// 🔹 Basic Dynamic Routes
	router.Get("/users/:id", func(c *rapidgo.Context) {
		id := c.Param("id")
		c.JSON(200, map[string]string{"message": "User Profile", "id": id})
	})
	router.Get("/products/:sku", func(c *rapidgo.Context) {
		sku := c.Param("sku")
		c.JSON(200, map[string]string{"message": "Product Details", "sku": sku})
	})

	// 🔹 Nested Dynamic Routes
	router.Get("/users/:id/profile", func(c *rapidgo.Context) {
		id := c.Param("id")
		c.JSON(200, map[string]string{"message": "User Profile Page", "id": id})
	})
	router.Get("/users/:id/orders/:orderID", func(c *rapidgo.Context) {
		id := c.Param("id")
		orderID := c.Param("orderID")
		c.JSON(200, map[string]string{"message": "User Order Details", "id": id, "orderID": orderID})
	})

	// 🔹 Mixed Static & Dynamic
	router.Get("/blog/:category/:postID", func(c *rapidgo.Context) {
		category := c.Param("category")
		postID := c.Param("postID")
		c.JSON(200, map[string]string{"message": "Blog Post", "category": category, "postID": postID})
	})

	// 🔹 Wildcard Route (Catch-All)
	router.Get("/static/*filepath", func(c *rapidgo.Context) {
		filepath := c.Param("filepath")
		c.JSON(200, map[string]string{"message": "Serving Static File", "filepath": filepath})
	})

	// 🔹 Testing Edge Cases
	router.Get("/users/wow", func(c *rapidgo.Context) {
		c.JSON(200, map[string]string{"message": "Special Route"})
	})
	router.Get("/users/:id/settings", func(c *rapidgo.Context) {
		id := c.Param("id")
		c.JSON(200, map[string]string{"message": "User Settings", "id": id})
	})

	// Start server
	router.Listen()
}
