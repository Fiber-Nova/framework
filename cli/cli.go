package cli

import (
	"fmt"
	"os"
	"strings"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Handler     func(args []string) error
}

// CLI manages command-line interface
type CLI struct {
	commands map[string]*Command
}

// New creates a new CLI instance
func New() *CLI {
	return &CLI{
		commands: make(map[string]*Command),
	}
}

// Register registers a new command
func (cli *CLI) Register(name string, description string, handler func(args []string) error) {
	cli.commands[name] = &Command{
		Name:        name,
		Description: description,
		Handler:     handler,
	}
}

// Run executes the CLI with the provided arguments
func (cli *CLI) Run(args []string) error {
	if len(args) < 2 {
		cli.showHelp()
		return nil
	}
	
	commandName := args[1]
	command, exists := cli.commands[commandName]
	
	if !exists {
		fmt.Printf("Unknown command: %s\n", commandName)
		cli.showHelp()
		return fmt.Errorf("unknown command: %s", commandName)
	}
	
	return command.Handler(args[2:])
}

// showHelp displays help information
func (cli *CLI) showHelp() {
	fmt.Println("FiberNova Framework - Artisan CLI")
	fmt.Println("\nAvailable commands:")
	
	for name, cmd := range cli.commands {
		fmt.Printf("  %-20s %s\n", name, cmd.Description)
	}
}

// DefaultCLI creates a CLI instance with default commands
func DefaultCLI() *CLI {
	cli := New()
	
	cli.Register("serve", "Start the development server", func(args []string) error {
		port := "3000"
		if len(args) > 0 && strings.HasPrefix(args[0], "--port=") {
			port = strings.TrimPrefix(args[0], "--port=")
		}
		fmt.Printf("Starting server on port %s...\n", port)
		return nil
	})
	
	cli.Register("make:controller", "Create a new controller", func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("controller name is required")
		}
		name := args[0]
		fmt.Printf("Creating controller: %s\n", name)
		return createController(name)
	})
	
	cli.Register("make:model", "Create a new model", func(args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("model name is required")
		}
		name := args[0]
		fmt.Printf("Creating model: %s\n", name)
		return createModel(name)
	})
	
	cli.Register("migrate", "Run database migrations", func(args []string) error {
		fmt.Println("Running migrations...")
		return nil
	})
	
	return cli
}

func createController(name string) error {
	template := fmt.Sprintf(`package controllers

import "github.com/gofiber/fiber/v2"

// %s handles HTTP requests
type %s struct{}

// Index handles GET requests
func (ctrl *%s) Index(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello from %s",
	})
}
`, name, name, name, name)
	
	filename := fmt.Sprintf("app/controllers/%s.go", strings.ToLower(name))
	os.MkdirAll("app/controllers", 0755)
	return os.WriteFile(filename, []byte(template), 0644)
}

func createModel(name string) error {
	template := fmt.Sprintf(`package models

import "github.com/CasperHK/FiberNova/database"

// %s represents a %s model
type %s struct {
	database.Model
	// Add your fields here
}
`, name, strings.ToLower(name), name)
	
	filename := fmt.Sprintf("app/models/%s.go", strings.ToLower(name))
	os.MkdirAll("app/models", 0755)
	return os.WriteFile(filename, []byte(template), 0644)
}
