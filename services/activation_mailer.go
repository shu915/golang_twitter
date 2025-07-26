package services

import (
	"fmt"
	"net/smtp"
)

func SendActivationEmail(to, token string) error {
	activationURL := fmt.Sprintf("http://localhost:8080/activate?token=%s", token)
	subject := "アカウント有効化のご案内"
	body := fmt.Sprintf("以下のリンクをクリックしてアカウントを有効化してください：\n\n%s", activationURL)

	from := "no-reply@example.com"

	msg := []byte(
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"\n" +
		body + "\n")

	err := smtp.SendMail("mailhog:1025", nil, from, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}
