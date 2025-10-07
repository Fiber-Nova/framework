package controllers

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"backend-meta-data/models"
)

var gormOnceDB *gorm.DB

func getDB() (*gorm.DB, error) {
	if gormOnceDB != nil {
		return gormOnceDB, nil
	}
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		portStr = "3306"
	}
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4,utf8", user, pass, host, portStr, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	gormOnceDB = db
	return db, nil
}

func sendSMTP(to string, subject string, htmlBody string) error {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USERNAME")
	pass := os.Getenv("SMTP_PASSWORD")
	fromName := os.Getenv("SMTP_FROM_NAME")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	useTLS := strings.ToLower(os.Getenv("SMTP_USE_TLS")) == "true"
	if host == "" || portStr == "" || fromEmail == "" {
		return fmt.Errorf("smtp config missing")
	}
	port, _ := strconv.Atoi(portStr)

	from := fromEmail
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, fromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", user, pass, host)

	// Try STARTTLS if useTLS
	if useTLS {
		c, err := smtp.Dial(addr)
		if err != nil {
			return err
		}
		defer c.Close()
		_ = c.Hello("localhost")
		if ok, _ := c.Extension("STARTTLS"); ok {
			cfg := &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}
			if err := c.StartTLS(cfg); err != nil {
				return err
			}
		}
		if user != "" {
			if ok, _ := c.Extension("AUTH"); ok {
				if err := c.Auth(auth); err != nil {
					return err
				}
			}
		}
		if err := c.Mail(from); err != nil {
			return err
		}
		if err := c.Rcpt(to); err != nil {
			return err
		}
		wc, err := c.Data()
		if err != nil {
			return err
		}
		_, _ = wc.Write([]byte(msg.String()))
		_ = wc.Close()
		return c.Quit()
	}

	// Plain send (will STARTTLS if server supports)
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg.String()))
}

// CreateMaintNotice handles POST /api/maint-notices
func CreateMaintNotice() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload struct {
			Station  string     `json:"station"`
			From     *time.Time `json:"from_time"`
			To       *time.Time `json:"to_time"`
			Until    bool       `json:"until_further"`
			ToAddr   string     `json:"to"`
			Template string     `json:"template"`
			Subject  string     `json:"subject"`
			Body     string     `json:"body"`
			SentBy   string     `json:"sent_by"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid payload")
		}
		if payload.ToAddr == "" || payload.Station == "" || payload.Subject == "" || payload.Body == "" {
			return fiber.NewError(fiber.StatusBadRequest, "missing required fields")
		}
		db, err := getDB()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "db error")
		}
		// ensure table exists
		_ = db.AutoMigrate(&models.MaintNoticeEmail{})
		n := models.MaintNoticeEmail{
			Station:      payload.Station,
			FromTime:     payload.From,
			ToTime:       payload.To,
			UntilFurther: payload.Until,
			To:           payload.ToAddr,
			Template:     payload.Template,
			Subject:      payload.Subject,
			Body:         payload.Body,
			SentBy:       payload.SentBy,
			SentAt:       time.Now(),
		}
		if err := db.Create(&n).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "save failed")
		}
		if err := sendSMTP(n.To, n.Subject, n.Body); err != nil {
			// log but still return success with warning
			fmt.Println("smtp send failed:", err)
		}
		return c.JSON(fiber.Map{"data": n})
	}
}
