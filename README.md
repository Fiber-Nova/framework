# FiberNova

**FiberNova** is a full-stack, Laravel-inspired web framework built on [Go Fiber](https://gofiber.io/), designed to combine blazing-fast performance with elegant, structured development for modern cloud-native applications.

## Features

✨ **Laravel-Inspired Design** - Familiar patterns for developers coming from Laravel  
⚡ **Blazing Fast** - Built on Go Fiber for exceptional performance  
🏗️ **Structured** - Clean architecture with separation of concerns  
🔒 **Secure** - Built-in middleware for authentication and CORS  
📝 **Validation** - Fluent validation system  
💾 **Database** - Query builder with ORM-like interface  
🛠️ **CLI Tool** - Artisan-inspired command-line interface  
☁️ **Cloud-Native** - Ready for containerization and microservices  

## Installation

```bash
go get github.com/CasperHK/FiberNova
```

## Quick Start

```go
package main

import (
    "github.com/CasperHK/FiberNova"
    "github.com/CasperHK/FiberNova/config"
    "github.com/CasperHK/FiberNova/middleware"
    "github.com/gofiber/fiber/v2"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Create application
    app := fibernova.New(cfg)
    
    // Add middleware
    app.Use(middleware.Logger())
    app.Use(middleware.CORS())
    
    // Define routes
    router := app.Router()
    router.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Welcome to FiberNova!",
        })
    })
    
    // Start server
    app.Listen()
}
```

## Routing

FiberNova provides a Laravel-inspired routing system with method chaining and route groups:

```go
router := app.Router()

// Basic routes
router.Get("/users", listUsers)
router.Post("/users", createUser)
router.Put("/users/:id", updateUser)
router.Delete("/users/:id", deleteUser)

// Route groups
router.Group("/api", func(group *routing.RouterGroup) {
    group.Get("/health", healthCheck)
    group.Get("/version", versionInfo)
})

// Protected routes with middleware
router.Group("/admin", func(group *routing.RouterGroup) {
    group.Middleware(middleware.Auth("your-secret-key"))
    group.Get("/dashboard", adminDashboard)
})
```

## Validation

Built-in fluent validation system:

```go
import "github.com/CasperHK/FiberNova/validation"

validator := validation.New()
validator.Required("email", email).Email("email", email)
validator.Required("password", password).Min("password", password, 8)

if validator.Fails() {
    return c.Status(422).JSON(fiber.Map{
        "errors": validator.Errors(),
    })
}
```

## Middleware

FiberNova includes several built-in middleware:

```go
import "github.com/CasperHK/FiberNova/middleware"

// Logger middleware
app.Use(middleware.Logger())

// CORS middleware
app.Use(middleware.CORS())

// Authentication middleware
app.Use(middleware.Auth("your-secret-key"))

// Rate limiting
app.Use(middleware.RateLimiter(100, time.Minute))
```

## Database

Query builder with ORM-like interface:

```go
import "github.com/CasperHK/FiberNova/database"

// Connect to database
db, err := database.New("postgres", "postgresql://user:pass@localhost/db")

// Query builder
users := db.Table("users").
    Where("status", "=", "active").
    OrderBy("created_at", "DESC").
    Limit(10).
    Get()

// Get first result
user := db.Table("users").
    Where("id", "=", 1).
    First()
```

## Configuration

FiberNova uses environment variables for configuration:

```env
APP_NAME=FiberNova
APP_ENV=development
APP_DEBUG=true
APP_KEY=your-secret-key

DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fibernova
DB_USER=root
DB_PASS=secret

SERVER_HOST=0.0.0.0
SERVER_PORT=3000
```

## CLI Tool

FiberNova includes an Artisan-inspired CLI tool:

```bash
# Start development server
go run artisan.go serve --port=3000

# Create a controller
go run artisan.go make:controller UserController

# Create a model
go run artisan.go make:model User

# Run migrations
go run artisan.go migrate
```

## Project Structure

```
.
├── app/
│   ├── controllers/    # HTTP controllers
│   ├── models/         # Database models
│   └── application.go  # Core application
├── config/             # Configuration
├── database/           # Database layer
├── routing/            # Routing system
├── middleware/         # Middleware
├── validation/         # Validation system
├── cli/                # CLI tool
└── examples/           # Example applications
```

## Example Application

See the [examples](examples/) directory for a complete example application demonstrating all features.

## Performance

FiberNova is built on Go Fiber, one of the fastest web frameworks:

- Zero memory allocation routing
- Low memory footprint
- Express-inspired API
- Optimized for high performance

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

FiberNova is open-source software licensed under the MIT license.

## Acknowledgments

- Built on [Go Fiber](https://gofiber.io/)
- Inspired by [Laravel](https://laravel.com/)

---

Made with ❤️ for the Go community