package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Sender    string
	Username  string
	Password  string
	Templates map[string]EmailTemplate
}

type EmailTemplate struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailData struct {
	Name string
}

func (ES *EmailSender) LoadTemplates(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open templates file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	templates := make(map[string]EmailTemplate)
	if err := decoder.Decode(&templates); err != nil {
		return fmt.Errorf("failed to decode templates: %w", err)
	}

	ES.Templates = templates
	return nil
}

func (ES *EmailSender) SendAuthEmail(templateName string, receiver, name string) error {
	if ES.Templates == nil {
		return fmt.Errorf("templates not loaded")
	}

	emailTemplate, ok := ES.Templates[templateName]
	if !ok {
		return fmt.Errorf("template '%s' not found", templateName)
	}

	data := EmailData{Name: name}

	subjectTpl, err := template.New("subject").Parse(emailTemplate.Subject)
	if err != nil {
		return fmt.Errorf("failed to parse subject template: %w", err)
	}
	var subjectBuffer bytes.Buffer
	if err := subjectTpl.Execute(&subjectBuffer, data); err != nil {
		return fmt.Errorf("failed to execute subject template: %w", err)
	}

	bodyTpl, err := template.New("body").Parse(emailTemplate.Body)
	if err != nil {
		return fmt.Errorf("failed to parse body template: %w", err)
	}
	var bodyBuffer bytes.Buffer
	if err := bodyTpl.Execute(&bodyBuffer, data); err != nil {
		return fmt.Errorf("failed to execute body template: %w", err)
	}

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, ES.Username, ES.Password)

	m := gomail.NewMessage()
	m.SetHeader("From", ES.Sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", subjectBuffer.String())
	m.SetBody("text/plain", bodyBuffer.String())

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (ES *EmailSender) SendUserEmail(receiver string, metadata map[string]string) error {
	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, ES.Username, ES.Password)

	m := gomail.NewMessage()
	m.SetHeader("From", ES.Sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", metadata["subject"])
	m.SetBody("text/plain", metadata["body"])

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
