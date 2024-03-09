package service

import (
	"database/sql"
	"fmt"
	"net/smtp"
	"os"
)

type MailService struct {
	db *sql.DB

	from string
	pass string

	host string
	port string
}

func NewMailService(db *sql.DB) *MailService {
	return &MailService{
		db: db,
		from: os.Getenv("EMAIL"),
		pass: os.Getenv("EMAIL_PASSWORD"),
		host: os.Getenv("SMTP_HOST"),
		port: os.Getenv("SMTP_PORT"),
	}
}

func (s *MailService) sendActivationLink(to []string, link string) error {
	subject := "Account activation"
	body := fmt.Sprintf(`
		<div>
			<h1>To activate your account, follow the link</h1>
			<a href=%s>%s</a>
		</div>
	`, link, link)

	msg := []byte("Subject: " + subject + "\r\n" +
	"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
	"\r\n" + body)

	auth := smtp.PlainAuth("", s.from, s.pass, s.host)

	if err := smtp.SendMail(s.host + ":" + s.port, auth, s.from, to, msg); err != nil {
		return err
	}

	return nil
}
