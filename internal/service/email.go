package service

import (
	"crypto/tls"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"grates/internal/repository"
	"grates/pkg/utils"
)

var AlreadyConfirmedErr = fmt.Errorf("email already confirmed")

type EmailService struct {
	D    EmailDeps
	repo repository.EmailRepository
}

func NewEmailService(repo repository.EmailRepository, d EmailDeps) *EmailService {
	return &EmailService{D: d, repo: repo}
}

// ReplaceConfirmationEmail если письмо уже существует, удаляет его, создает новое и отправляет
func (e *EmailService) ReplaceConfirmationEmail(userId int, to, name string) error {
	hash, err := utils.GenerateHash(32)
	if err != nil {
		hash = utils.RandStringBytesRmndr(32)
	}

	if err = e.repo.ReplaceEmail(userId, hash); err != nil {
		return err
	}

	return e.SendAuthEmail(to, name, hash)
}

// ConfirmEmail по переданному hash'у подтверждает аккаунт
func (e *EmailService) ConfirmEmail(hash string) error {
	if isConirmed, _ := e.repo.IsConfirmed(hash); isConirmed {
		return AlreadyConfirmedErr
	}

	return e.repo.ConfirmEmail(hash)
}

// sendAuthEmail отправляет письмо на почту, интергируя в него name, hash
func (e *EmailService) SendAuthEmail(to, name, hash string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.D.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Grates | Подтверждение аккаунта")
	m.SetBody("text/html", e.getEmailTemplate(name, hash))

	d := gomail.NewDialer(e.D.SmtpHost, e.D.SmtpPort, e.D.From, e.D.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}

// getEmailTemplate возвращает html-верстку письма
func (e *EmailService) getEmailTemplate(name, hash string) string {
	confirmUrl := fmt.Sprintf("http://%s/auth/confirm?hash=%s", e.D.BaseUrl, hash)
	return `

<!DOCTYPE html>
<html>
<head>
    <title>Подтверждение аккаунта</title>
    <meta charset="utf-8">
</head>
<body style="font-family: Arial, sans-serif;
background-color: #f4f4f8;
margin: 0;
padding: 0;">
    <div class="container" style="background-color: white;
    width: 60% ;
    margin: 20px auto;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);">
        <header style="background-color: #8a2be2;
        color: white;
        padding: 10px 20px;
        text-align: center;
        border-top-left-radius: 8px;
        border-top-right-radius: 8px;">
            Добро пожаловать в Grates!
        </header>
        <div class="content" style="padding: 20px; text-align: center;">
            <h1>Подтверждение Аккаунта</h1>
            <p>Привет, ` + name + ` !</p>
            <p>Спасибо за регистрацию в Grates. Пожалуйста, нажмите на кнопку ниже, чтобы подтвердить свой аккаунт.</p>
                <a href="` + confirmUrl + `" target="_blank" class="submit button" style="cursor: pointer;
                display: inline-block;
                padding: 10px 20px;
                margin-top: 20px;
                background-color: #8a2be2;
                color: white;
                text-decoration: none;
                border-radius: 4px;
                border: 0;
                font-weight: bold;">Подтвердить Аккаунт</a>
        <footer style="text-align: center; margin-top: 30px; font-size: 0.9em; color: #666;">
            Если у вас возникли вопросы, свяжитесь с нами по адресу ` + e.D.From + `
        </footer>
</body>
</html>

`
}
