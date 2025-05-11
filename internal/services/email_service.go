package services

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/go-mail/mail/v2"
)

// EmailService реализует отправку email-уведомлений через SMTP
type EmailService struct {
	host     string
	port     int
	user     string
	password string
}

// NewEmailService создает новый сервис email-уведомлений
func NewEmailService(host string, port int, user, password string) *EmailService {
	return &EmailService{host: host, port: port, user: user, password: password}
}

// Send отправляет email
func (s *EmailService) Send(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := mail.NewDialer(s.host, s.port, s.user, s.password)
	d.TLSConfig = &tls.Config{ServerName: s.host, InsecureSkipVerify: false}
	if err := d.DialAndSend(m); err != nil {
		log.Printf("SMTP error: %v", err)
		return fmt.Errorf("email sending failed")
	}
	return nil
}
