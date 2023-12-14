package service

import (
	"crypto/tls"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"grates/internal/repository"
	"grates/pkg/utils"
)

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

	return e.sendAuthEmail(to, name, hash)
}

// ConfirmEmail по переданному hash'у подтверждает аккаунт
func (e *EmailService) ConfirmEmail(hash string) error {
	return e.repo.ConfirmEmail(hash)
}

// sendAuthEmail отправляет письмо на почту, интергируя в него name, hash
func (e *EmailService) sendAuthEmail(to, name, hash string) error {
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
			cursor: pointer;
            display: inline-block;
            padding: 10px 20px;
            margin-top: 20px;
            background-color: #8a2be2;
            color: white;
            text-decoration: none;
            border-radius: 4px;
      		border: 0;
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
          <form action="` + confirmUrl + `" method="post">
    			<input type="submit" " value="Подтвердить Аккаунт" class="button" />
			</form> 
        <div class="footer">
            Если у вас возникли вопросы, свяжитесь с нами по адресу ` + e.D.From + `
        </div>
    </div>
</body>
</html>

`
}
