package email

import (
	"bytes"
	"fmt"
	"github.com/stanislavCasciuc/atom-fit-go/internal/config"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
)

const userVerificationTemplPath = "./internal/lib/email/templates/verify-email.html"

func send(to []string, subject string, body string) error {
	const op = "email.send"

	emailCfg := config.Envs.Email

	m := gomail.NewMessage()
	m.SetHeader("From", emailCfg.Addr)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(emailCfg.Host, emailCfg.Port, emailCfg.Addr, emailCfg.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func SendVerifyUser(username, email, code string) error {
	const op = "email.SendVerifyUser"
	var body bytes.Buffer
	t, err := template.ParseFiles(userVerificationTemplPath)
	if err != nil {
		log.Fatal("cannot to parse email template")
	}
	t.Execute(
		&body, struct {
			Name string
			Code string
		}{Name: username, Code: code},
	)
	bodyStr := body.String()
	err = send([]string{email}, "User Verification", bodyStr)
	if err != nil {
		return err
	}

	return nil
}
