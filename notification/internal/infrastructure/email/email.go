package email

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"net/url"

	"github.com/go-gomail/gomail"

	"github.com/Arclight-V/mtch/notification/internal/usecase/notification"
	config "github.com/Arclight-V/mtch/pkg/platform/config"
)

const emailTpl = `<!doctype html>
<title>Подтверждение email</title>
<meta charset="utf-8">
<body>
  <h1>Подтвердите email</h1>
  <p>Перейдите по ссылке:</p>
  <p><a href="{{.VerifyURL}}">{{.VerifyURL}}</a></p>
  <p>Ссылка действует 24 часа.</p>
</body>`

// TODO: move it
const verifyEmailURL = "https://localhost:8000/api/v1/auth/verify-email"

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

type EmailData struct {
	VerifyURL string
}

func makeVerifyURL(base string, token string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func renderEmailHTML(baseVerifyURL, token string) (string, error) {
	verifyURL, err := makeVerifyURL(baseVerifyURL, token)
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("email").Parse(emailTpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, EmailData{VerifyURL: verifyURL}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *SMTPClient) SendUserRegistered(ctx context.Context, vd notification.VerifyData) error {
	log.Printf("sending userservice registered to: %v", vd.Email)
	u, err := makeVerifyURL(verifyEmailURL, vd.VerifyToken)
	if err != nil {
		return err
	}
	htmpl, err := renderEmailHTML(u, vd.VerifyToken)
	if err != nil {
		return err
	}

	errSend := s.send(ctx, vd.Email, "subject", htmpl)
	if errSend != nil {
		return errSend
	}
	return nil
}
