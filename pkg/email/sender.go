package email

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Sender   string
	Password string
}

func (ES *EmailSender) SendEmail(data map[string]string, receiver string) error {
	d := gomail.NewDialer("smtp.gmail.com", 587, ES.Sender, ES.Password)

	m := gomail.NewMessage()
	m.SetHeader("From", ES.Sender)
	m.SetHeader("To", receiver)

	// Set subject and body from the map
	subject, ok := data["subject"]
	if !ok {
		return fmt.Errorf("missing 'subject' in data")
	}
	body, ok := data["body"]
	if !ok {
		return fmt.Errorf("missing 'body' in data")
	}

	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Optionally attach a file if provided in the map
	if filePath, ok := data["attachment"]; ok {
		m.Attach(filePath)
	}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
