package email

import (
	"config"
	"context"
	"github.com/go-gomail/gomail"
	"log"
)

const htmpl = `<!doctype html>
<title>Email verified</title>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<body>
  <h1>Email подтвержден ✅</h1>
  <p>Теперь вы можете заполнить профиль через API <code>PUT /v1/users/me/profile</code>.</p>
</body>`

type SMTPClient struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

func NewSMTPClient(cfg *config.Config) *SMTPClient {
	return &SMTPClient{Host: cfg.SMTPClient.Host,
		Port: cfg.SMTPClient.Port,
		User: cfg.SMTPClient.User,
		Pass: cfg.SMTPClient.Pass,
		From: cfg.SMTPClient.From}
}

func (s *SMTPClient) send(ctx context.Context, to, subject, html string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)
	d := gomail.NewDialer(s.Host, s.Port, s.User, s.Pass)
	return d.DialAndSend(m)
}

func (s *SMTPClient) SendUserRegistered(ctx context.Context, to string) error {
	log.Printf("sending user registered to: %v", to)
	err := s.send(ctx, to, "subject", htmpl)
	if err != nil {
		return err
	}
	return nil
}
