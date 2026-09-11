package mail

import (
	"log"
	"net/smtp"
	"os"
)

func Send(to, subject, body string) error {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		log.Printf("почта в консоль: кому=%s | тема=%s | текст=%s", to, subject, body)
		return nil
	}

	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = user
	}

	auth := smtp.PlainAuth("", user, pass, host)
	msg := "From: " + from + "\r\nTo: " + to + "\r\nSubject: " + subject + "\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" + body
	return smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
}

func SendVerification(to, code string) error {
	return Send(to,
		"Форум — подтверждение email",
		"Ваш код подтверждения: "+code+"\n\nКод действителен 30 минут.",
	)
}