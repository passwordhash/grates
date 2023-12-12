package service

import (
	"crypto/tls"
	gomail "gopkg.in/mail.v2"
)

type EmailService struct {
	D EmailDeps
}

func NewEmailService(d EmailDeps) *EmailService {
	return &EmailService{D: d}
}

func (e *EmailService) SendAuthEmail(to, name string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.D.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Grates | Подтверждение аккаунта")
	m.SetBody("text/html", e.getEmailTemplate(name))

	d := gomail.NewDialer(e.D.SmtpHost, e.D.SmtpPort, e.D.From, e.D.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}

func (e *EmailService) getEmailTemplate(name string) string {
	return `

	<!DOCTYPE html>
<html>
<head>
    <title>Подтверждение аккаунта</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f8;
            margin: 0;
            padding: 0;
        }
        .container {
            background-color: white;
            width: 60 % ;
            margin: 20px auto;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
        }
        .header {
            background-color: #8a2be2;
            color: white;
            padding: 10px 20px;
            text-align: center;
            border-top-left-radius: 8px;
            border-top-right-radius: 8px;
        }
        .content {
            padding: 20px;
            text-align: center;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            margin-top: 20px;
            background-color: #8a2be2;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            font-weight: bold;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 0.9em;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            Добро пожаловать в Grates!
        </div>
        <div class="content">
            <h1>Подтверждение Аккаунта</h1>
            <p>Привет, ` + name + `!</p>
            <p>Спасибо за регистрацию в Grates. Пожалуйста, нажмите на кнопку ниже, чтобы подтвердить свой аккаунт.</p>
            <a href="ССЫЛКА_ДЛЯ_ПОДТВЕРЖДЕНИЯ" class="button">Подтвердить Аккаунт</a>
        </div>
        <div class="footer">
            Если у вас возникли вопросы, свяжитесь с нами по адресу ` + e.D.From + `
        </div>
    </div>
</body>
</html>

`
}
