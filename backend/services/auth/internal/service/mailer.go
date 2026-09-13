package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

type MailService interface {
	SendCode(toEmail, code string) error
}

type SMTPMailer struct {
	host     string 
	port     int    
	username string 
	password string 
	from     string
}

func NewSMTPMailer() *SMTPMailer {
	host := getEnv("SMTP_HOST", "localhost")
	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "25"))
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	from := getEnv("SMTP_FROM", "no-reply@example.com")

	return &SMTPMailer{host, port, username, password, from}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (m *SMTPMailer) SendCode(toEmail, code string) error {
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: Your Verification Code\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
			"\r\n"+
			"Your one-time verification code is: %s\r\n",
		m.from, toEmail, code,
	))

	var auth smtp.Auth
	if m.username != "" && m.password != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}

	if m.port == 465 {
		tlsConfig := &tls.Config{ServerName: m.host}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}
		c, err := smtp.NewClient(conn, m.host)
		if err != nil {
			return err
		}
		defer c.Close()

		if auth != nil {
			if err := c.Auth(auth); err != nil {
				return err
			}
		}
		if err := c.Mail(m.from); err != nil {
			return err
		}
		if err := c.Rcpt(toEmail); err != nil {
			return err
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		if _, err := w.Write(msg); err != nil {
			return err
		}
		return w.Close()
	}

	return smtp.SendMail(addr, auth, m.from, []string{toEmail}, msg)
}
