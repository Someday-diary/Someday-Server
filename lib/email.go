package lib

import (
	"bytes"
	"html/template"
	"os"

	"gopkg.in/gomail.v2"
)

type EmailRequest struct {
	from    string
	to      string
	subject string
	body    string
}

func NewRequest(to, from, subject, body string) *EmailRequest {
	return &EmailRequest{
		to:      to,
		subject: subject,
		body:    body,
		from:    from,
	}
}

func (r *EmailRequest) SendEmail() error {
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("smtp_id"), os.Getenv("smtp_pwd"))

	m := gomail.NewMessage()

	m.SetHeader("From", r.from)
	m.SetHeader("To", r.to)
	m.SetHeader("Subject", r.subject)
	m.SetBody("text/html", r.body)

	err := d.DialAndSend(m)

	return err
}

func (r *EmailRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
