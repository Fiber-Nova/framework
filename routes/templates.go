package routes

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func ListMaintNoticeTemplates() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Locate the templates directory relative to the project
		base := filepath.Join(".", "src", "backend", "email_templates", "maint_notice")
		entries, err := os.ReadDir(base)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read templates"})
		}
		files := make([]string, 0, len(entries))
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			ext := filepath.Ext(name)
			if ext == ".html" || ext == ".htm" {
				files = append(files, name)
			}
		}
		return c.JSON(fiber.Map{"data": files})
	}
}
