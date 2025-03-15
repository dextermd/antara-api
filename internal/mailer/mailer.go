package mailer

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
)

//go:embed templates
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string // "no-reply@yourdomain.com"
	Logger echo.Logger
}

type EmailData struct {
	AppName string
	Subject string
	Meta    any
}

func NewMailer(logger echo.Logger) Mailer {
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		logger.Fatal(err)
	}
	mailHost := os.Getenv("MAIL_HOST")
	mailUser := os.Getenv("MAIL_USER")
	mailPass := os.Getenv("MAIL_PASS")
	mailSender := os.Getenv("MAIL_SENDER")
	dialer := gomail.NewDialer(mailHost, mailPort, mailUser, mailPass)

	return Mailer{
		dialer: dialer,
		sender: mailSender,
		Logger: logger,
	}
}

func (mailer *Mailer) SendMail(recipient string, templateFile string, data EmailData) error {
	absolutePath := filepath.Join("templates/" + templateFile)

	tmpl, err := template.ParseFS(templateFS, absolutePath)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	data.AppName = os.Getenv("APP_NAME")

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	goMailMessage := gomail.NewMessage()
	goMailMessage.SetHeader("To", recipient)
	goMailMessage.SetHeader("From", mailer.sender)
	goMailMessage.SetHeader("Subject", subject.String())

	goMailMessage.SetBody("text/html", htmlBody.String())

	err = mailer.dialer.DialAndSend(goMailMessage)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	return nil
}
