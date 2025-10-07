package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var tmplNameRe = regexp.MustCompile(`(?i)Template\s*Name\s*:\s*(.+)`) // capture after 'Template Name:'

func resolveMaintTemplatesDir() string {
	candidates := []string{
		filepath.Join(".", "src", "backend", "email_templates", "maint_notice"),
		filepath.Join(".", "backend", "email_templates", "maint_notice"),
		filepath.Join(".", "email_templates", "maint_notice"),
	}
	for _, p := range candidates {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			return p
		}
	}
	return candidates[0]
}

func ListMaintNoticeTemplates() fiber.Handler {
	return func(c *fiber.Ctx) error {
		base := resolveMaintTemplatesDir()
		entries, err := os.ReadDir(base)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to read templates"})
		}
		type item struct {
			Name  string `json:"name"`
			Label string `json:"label"`
		}
		out := make([]item, 0, len(entries))
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			ext := strings.ToLower(filepath.Ext(name))
			if ext != ".html" && ext != ".htm" {
				continue
			}
			label := name
			if b, err := os.ReadFile(filepath.Join(base, name)); err == nil {
				if m := tmplNameRe.FindSubmatch(b); len(m) > 1 {
					label = strings.TrimSpace(string(m[1]))
				}
			}
			out = append(out, item{Name: name, Label: label})
		}
		return c.JSON(fiber.Map{"data": out})
	}
}

// ListTemplates currently proxies to maintenance notice templates; extend as needed for other template groups
func ListTemplates() fiber.Handler { return ListMaintNoticeTemplates() }

func GetMaintNoticeTemplate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		name := c.Params("name")
		base := filepath.Base(name) // prevent path traversal
		if base == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "missing template name"})
		}
		if ext := strings.ToLower(filepath.Ext(base)); ext != ".html" && ext != ".htm" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "unsupported template extension"})
		}
		dir := resolveMaintTemplatesDir()
		path := filepath.Join(dir, base)
		b, err := os.ReadFile(path)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "template not found"})
		}
		return c.Type("text/html; charset=utf-8").Send(b)
	}
}
