package service

import (
	"fmt"
	"net/smtp"

	"golang-api-template/internal/config"
	"golang-api-template/internal/i18n"
)

type EmailService struct {
	config *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// SendEmail sends a generic email
func (s *EmailService) SendEmail(to, subject, body string) error {
	from := s.config.Sender
	pass := s.config.Password
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	// Setup message
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// Setup SMTP authentication information.
	auth := smtp.PlainAuth("", s.config.Username, pass, s.config.Host)

	// Send the email
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

// SendPasswordResetEmail sends a localized password reset email
func (s *EmailService) SendPasswordResetEmail(to, resetLink, lang string) error {
	subject := i18n.TT(lang, "PasswordResetEmailSubject")
	body := fmt.Sprintf(i18n.TT(lang, "PasswordResetEmailBody"), resetLink)
	return s.SendEmail(to, subject, body)
}
