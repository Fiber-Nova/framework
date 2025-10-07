package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"backend-meta-data/config"
	"backend-meta-data/db"
	"backend-meta-data/middleware"
	"backend-meta-data/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	logrus "github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		fmt.Println("\033[31mBye! Bye! See you next time.\033[0m")
	}()

	middleware.InitLogger()

	logrus.Info("Welcome to the AWS Meta Data backend services.")

	// Load config
	configPath := "./.env"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	// Export AUTO_MIGRATE env for lower-level controls
	if cfg.App.AutoMigrate {
		os.Setenv("AUTO_MIGRATE", "true")
	} else {
		os.Setenv("AUTO_MIGRATE", "false")
	}

	// Connect to DB
	dbConn, err := db.ConnectDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	defer dbConn.Close()

	// GORM connection for migration
	gormDB, err := db.ConnectGormDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Error connecting to GORM DB: %v", err)
	}
	// Only run migrations if enabled or in local/dev
	if os.Getenv("AUTO_MIGRATE") != "false" && (cfg.App.Env == "local" || cfg.App.Env == "dev" || cfg.App.Env == "development") {
		if err := db.InitDBIfNeeded(gormDB); err != nil {
			log.Fatalf("Error running AutoMigrate: %v", err)
		}
	}

	app := fiber.New()

	// Enable CORS for frontend connection
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(middleware.LoggerMiddleware)

	// Register routes with both dbConn and gormDB
	routes.RegisterRoutes(app, dbConn, gormDB)

	// Initialize Casbin and attach middleware
	if err := middleware.InitCasbin(gormDB); err != nil {
		log.Fatalf("Error initializing Casbin: %v", err)
	}
	// Protect API routes; allow auth/health without RBAC
	app.Use(func(c *fiber.Ctx) error {
		p := c.Path()
		if strings.HasPrefix(p, "/api/") {
			return middleware.CasbinMiddleware()(c)
		}
		return c.Next()
	})

	// Signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("\033[31mBye! Bye! See you next time.\033[0m")
		os.Exit(0)
	}()

	// Start Fiber server on configured port
	if err := app.Listen(":" + strconv.Itoa(cfg.App.Port)); err != nil {
		logrus.Fatalf("Fiber failed to start: %v", err)
	}
	logrus.Infof("Meta-Data backend server started on port %d", cfg.App.Port)
}
