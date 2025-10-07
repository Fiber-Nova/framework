package middleware

import (
	"backend-meta-data/models"
	"context"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

// InitCasbin sets up the Casbin enforcer with GORM adapter and seeds base roles/policies.
func InitCasbin(db *gorm.DB) error {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return err
	}
	e, err := casbin.NewEnforcer("config/casbin_model.conf", adapter)
	if err != nil {
		return err
	}
	// Register matcher helpers with correct signature
	e.AddFunction("keyMatch2", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return false, nil
		}
		key1, ok1 := args[0].(string)
		key2, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return false, nil
		}
		return util.KeyMatch2(key1, key2), nil
	})
	e.AddFunction("regexMatch", func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return false, nil
		}
		s, ok1 := args[0].(string)
		p, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			return false, nil
		}
		return util.RegexMatch(s, p), nil
	})

	if err := e.LoadPolicy(); err != nil {
		return err
	}
	e.EnableAutoSave(true)

	// Seed base roles and policies if missing
	seed := [][]string{
		// Root: full access
		{"p", "Root", "/api/*", "(GET|POST|PUT|PATCH|DELETE|OPTIONS)"},
		// Admin: manage core resources
		{"p", "Admin", "/api/users*", "(GET|POST|PUT|PATCH|DELETE)"},
		{"p", "Admin", "/api/instrument-types*", "(GET|POST|PUT|PATCH|DELETE)"},
		{"p", "Admin", "/api/*", "GET"},
		// Inspector: read-only plus submit inspection forms
		{"p", "Inspector", "/api/inspection-forms*", "(GET|POST)"},
		{"p", "Inspector", "/api/*", "GET"},
		// All authenticated users can update their own profile
		{"p", "*", "/api/users/profile", "PATCH"},
	}
	for _, rule := range seed {
		_, _ = e.AddPolicy(rule[1], rule[2], rule[3])
	}

	Enforcer = e
	return nil
}

// CasbinMiddleware enforces RBAC on incoming requests.
func CasbinMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if Enforcer == nil {
			return c.Next()
		}
		// Resolve subject from session
		sess, err := models.Store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		sub, _ := sess.Get("username").(string)
		if sub == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		obj := string(c.OriginalURL())
		act := string(c.Method())

		ok, err := Enforcer.Enforce(sub, obj, act)
		if err != nil {
			log.Printf("casbin enforce error: %v", err)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}

// AssignRole assigns a role to a username in Casbin policies.
func AssignRole(ctx context.Context, username, role string) error {
	if Enforcer == nil {
		return nil
	}
	_, err := Enforcer.AddGroupingPolicy(username, role)
	return err
}
