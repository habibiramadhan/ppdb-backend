// internal/core/services/email_service.go
package services

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendVerificationEmail(to string, token string, name string) error
	SendWelcomeEmail(to string, name string) error
	SendPasswordResetEmail(to string, token string, name string) error
	SendRegistrationConfirmation(to string, name string, registrationNumber string) error
}

type emailService struct {
	dialer *gomail.Dialer
	templates map[string]*template.Template
}

type EmailData struct {
	Name              string
	VerificationLink  string
	PasswordResetLink string
	RegistrationNumber string
	AppName           string
	AppURL            string
	Year              string
}

func NewEmailService() (EmailService, error) {
	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	templates := make(map[string]*template.Template)
	templatesDir := "templates/emails"
	templateFiles := []string{
		"verification.html",
		"welcome.html",
	}

	for _, file := range templateFiles {
		tmpl, err := template.ParseFiles(filepath.Join(templatesDir, file))
		if err != nil {
			return nil, fmt.Errorf("failed to parse email template %s: %v", file, err)
		}
		templates[file] = tmpl
	}

	return &emailService{
		dialer: dialer,
		templates: templates,
	}, nil
}

func (s *emailService) sendEmail(to string, subject string, templateName string, data EmailData) error {
	var body bytes.Buffer
	tmpl, exists := s.templates[templateName]
	if !exists {
		return fmt.Errorf("template %s not found", templateName)
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	return s.dialer.DialAndSend(m)
}

func (s *emailService) SendVerificationEmail(to string, token string, name string) error {
	data := EmailData{
		Name: name,
		VerificationLink: fmt.Sprintf("%s/verify-email?token=%s", 
			os.Getenv("APP_URL"), 
			token,
		),
		AppName: os.Getenv("APP_NAME"),
		AppURL:  os.Getenv("APP_URL"),
	}

	return s.sendEmail(
		to,
		"Verifikasi Email Anda",
		"verification.html",
		data,
	)
}

func (s *emailService) SendWelcomeEmail(to string, name string) error {
	data := EmailData{
		Name:    name,
		AppName: os.Getenv("APP_NAME"),
		AppURL:  os.Getenv("APP_URL"),
	}

	return s.sendEmail(
		to,
		"Selamat Datang di PPDB Online",
		"welcome.html",
		data,
	)
}

func (s *emailService) SendPasswordResetEmail(to string, token string, name string) error {
	data := EmailData{
		Name: name,
		PasswordResetLink: fmt.Sprintf("%s/reset-password?token=%s",
			os.Getenv("APP_URL"),
			token,
		),
		AppName: os.Getenv("APP_NAME"),
		AppURL:  os.Getenv("APP_URL"),
	}

	return s.sendEmail(
		to,
		"Reset Password Anda",
		"password_reset.html",
		data,
	)
}

func (s *emailService) SendRegistrationConfirmation(to string, name string, registrationNumber string) error {
	data := EmailData{
		Name:               name,
		RegistrationNumber: registrationNumber,
		AppName:           os.Getenv("APP_NAME"),
		AppURL:            os.Getenv("APP_URL"),
	}

	return s.sendEmail(
		to,
		"Konfirmasi Pendaftaran PPDB Online",
		"registration_confirmation.html",
		data,
	)
}