package controllers

import (
	"backend-meta-data/auth"
	"backend-meta-data/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthController handles user authentication
func AuthController(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		var user models.User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid username or password"})
		}

		if user.Password != req.Password {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid username or password"})
		}

		// TODO: Set session or return JWT token
		return c.JSON(fiber.Map{"message": "Login successful", "user": user})
	}
}

// OIDCAuthMiddleware validates JWT from Azure AD or ADFS
func OIDCAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// TODO: Replace with your Azure AD/ADFS public key or JWKS validation
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method and return public key
			// For demo, accept any method and no key (DO NOT use in production)
			return []byte("your_secret_key"), nil
		})
		if err != nil || !token.Valid {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}
		// Extract user info from claims (e.g., email, name, AD account)
		c.Locals("user", claims)
		return c.Next()
	}
}

// SSOController handles Windows AD SSO via LDAP
func SSOController(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get username from request header (e.g., REMOTE_USER, for IIS/AD integration)
		adUsername := c.Get("X-AD-Username")
		if adUsername == "" {
			return c.Status(401).JSON(fiber.Map{"error": "No AD username detected"})
		}

		// Use a service account to bind to LDAP and check if user exists
		ldapURL := os.Getenv("LDAP_URL")
		baseDN := os.Getenv("LDAP_BASE_DN")
		serviceUser := os.Getenv("LDAP_SERVICE_USER")
		servicePass := os.Getenv("LDAP_SERVICE_PASS")

		ok, err := auth.AuthenticateAD(ldapURL, baseDN, serviceUser, servicePass)
		if err != nil || !ok {
			return c.Status(500).JSON(fiber.Map{"error": "LDAP service bind failed"})
		}

		// Optionally, check if user exists in AD (not authenticating password here)
		// For demo, assume user exists if header is present

		// Find user in local DB
		var user models.User
		if err := db.Where("username = ?", adUsername).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "User not found in system"})
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"id":       user.ID,
		})
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Token generation failed"})
		}

		return c.JSON(fiber.Map{"token": tokenString, "user": user})
	}
}

// ADLoginHandler handles Active Directory login
func ADLoginHandler(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	ldapURL := os.Getenv("LDAP_URL")
	baseDN := os.Getenv("LDAP_BASE_DN")

	l, err := ldap.DialURL(ldapURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "LDAP connection failed"})
	}
	defer l.Close()

	userDN := "CN=" + req.Username + "," + baseDN
	if err := l.Bind(userDN, req.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "AD authentication failed"})
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT generation failed"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}
