package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string
}

func Load() (*Config, error) {
	get := func(key, def string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return def
	}

	host := get("SMTP_HOST", "smtp.gmail.com")
	portStr := get("SMTP_PORT", "587")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT %q: %w", portStr, err)
	}

	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := get("SMTP_FROM", user)
	if from == "" {
		from = user
	}

	if user == "" || pass == "" {
		return nil, fmt.Errorf("SMTP_USER and SMTP_PASS must be set for Gmail SMTP")
	}

	return &Config{
		SMTPHost: host,
		SMTPPort: port,
		SMTPUser: user,
		SMTPPass: pass,
		SMTPFrom: from,
	}, nil
}