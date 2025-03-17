# RapidGo - Fast, Minimalist Web Framework for Go
RapidGo is a lightweight and high-performance web framework for building modern web applications and APIs in Go. Designed with simplicity and speed in mind, RapidGo provides an intuitive API for routing, middleware integration, and request handling while maintaining minimal overhead.

## Features
- **Fast Routing**: Built on a radix tree-based router for efficient route matching.
- **Middleware Integration** : Supports both global and group-level middleware for flexible request handling.
- **Lightweight Design** : A clean and minimalistic API that prioritizes simplicity and rapid development.
- **Environment Management** : Built-in support for loading and managing environment variables.
- **Customizable Error Responses** : Define custom "Not Found" messages or fallback behavior for unmatched routes.
- **Debugging Capabilities** : Includes a debugging mode to log registered routes and server details during development.
- **Route Grouping** : Organize routes with shared base paths and middleware for better structure and reusability.

---

## Installation
To install RapidGo, use `go get`:
```bash
go get github.com/rwiteshbera/rapidgo
```

---

## Quick Start
Here’s how you can quickly set up a basic web server using RapidGo:
```go
package main
import (
	"github.com/rwiteshbera/rapidgo"
)
func main() {
	router := rapidgo.New()
	// Define a simple route
	router.Get("/", func(c *rapidgo.Context) {
		c.String(200, "Hello, World!")
	})
	// Define a route with parameters
	router.Get("/user/:id", func(c *rapidgo.Context) {
		id := c.Param("id")
		c.JSON(200, map[string]string{"id": id})
	})
	// Start the server
	if err := router.Listen(); err != nil {
		panic(err)
	}
}
```
Run the application:
```bash
go run main.go
```
By default, the server will start on port `8080`. You can customize the port by setting the `PORT` environment variable in terminal or in `.env` file:
```bash
PORT=3000 go run main.go
```
---
## Contributing
RapidGo is currently in its early stages of development. We appreciate your patience and encourage you to contribute or provide feedback as we continue to improve the framework. If you'd like to contribute, please follow these steps:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Open a pull request.
   
---
## Discussion
If you have any questions, need support, or want to share feedback, feel free to reach out by creating a new issue on our GitHub repository. Whether it’s a bug report, feature request, or general question, we welcome your input.

👉 **Create a new issue here**: [RapidGo Issues](https://github.com/rwiteshbera/rapidgo/issues/new)

---
Start building fast, scalable web applications with RapidGo today! 🚀