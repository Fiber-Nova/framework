package controllers

import (
	"backend-meta-data/middleware"
	"backend-meta-data/models"
	"database/sql"
	"os"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

var Store = session.New()

func Login(dbConn *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		ldapServer := os.Getenv("LDAP_SERVER")
		ldapPort := 389
		bindDN := req.Username
		bindPassword := req.Password

		l, err := ldap.Dial("tcp", ldapServer+":"+strconv.Itoa(ldapPort))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "LDAP connection failed"})
		}
		defer l.Close()

		if err := l.Bind(bindDN, bindPassword); err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Authentication failed"})
		}
		// TODO: Query database for user ID using username
		userID := 123 // Placeholder for actual user ID
		sess, err := Store.Get(c)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Session error"})
		}
		sess.Set("username", req.Username)
		sess.Set("userID", userID)
		sess.Save()
		return c.JSON(fiber.Map{"message": "Authentication successful", "user": fiber.Map{"username": req.Username, "userID": userID}})
	}
}

func Logout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := Store.Get(c)
		if err == nil {
			sess.Destroy()
		}
		return c.JSON(fiber.Map{"message": "Logout successful"})
	}
}

func Me(dbConn *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := Store.Get(c)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Session not found"})
		}
		username := sess.Get("username")
		userID := sess.Get("userID")
		if username == nil || userID == nil {
			return c.Status(401).JSON(fiber.Map{"error": "User not logged in"})
		}
		// Query database for full user object (placeholder)
		// TODO: Replace with actual DB query
		user := fiber.Map{
			"userID":      userID,
			"username":    username,
			"email":       "user@example.com",
			"displayName": "John Doe",
		}
		return c.JSON(fiber.Map{"user": user})
	}
}

// CreateUser handles POST /api/users to insert a new user using the models.User helper
func CreateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
		}
		u := models.User{Username: req.Username, Password: req.Password, Email: req.Email, Role: req.Role}
		if err := models.CreateUser(db, &u); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		// assign role in Casbin if provided
		if req.Role != "" {
			_ = middleware.AssignRole(c.Context(), u.Username, req.Role)
		}
		u.Password = ""
		// Ensure email is included from request since reload may not catch it
		u.Email = req.Email
		// Reload user to get all fields including email
		if err := db.First(&u, u.ID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch created user"})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": u})
	}
}

// ListUsers handles GET /api/users (passwords are not returned)
func ListUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		if err := db.Order("id asc").Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch users"})
		}
		for i := range users {
			users[i].Password = ""
		}
		return c.JSON(fiber.Map{"data": users})
	}
}

// UpdateUser handles PATCH /api/users/:id
func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ensure authentication
		loggedInUser, err := models.GetLoggedInUser(c, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		// Check permission using Casbin middleware
		obj := c.Path()   // e.g. "/api/users/123"
		act := c.Method() // "PATCH"

		ok, err := middleware.Enforcer.Enforce(loggedInUser.Username, obj, act)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
		}

		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
		}

		// Get current user from DB
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}

		// Update fields if provided in request
		if req.Username != "" {
			user.Username = req.Username
		}
		if req.Email != "" {
			user.Email = req.Email
		}
		if req.Role != "" {
			user.Role = req.Role
			// Update role in Casbin if changed
			if err := middleware.AssignRole(c.Context(), user.Username, req.Role); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update role"})
			}
		}
		if req.Password != "" {
			user.Password = req.Password
		}

		if err := db.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update user"})
		}

		// Clear password before returning
		user.Password = ""
		return c.JSON(fiber.Map{"data": user})
	}
}
