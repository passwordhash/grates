package main

import (
	"crypto/tls"
	gomail "gopkg.in/mail.v2"
)

func main() {
	//from := "iam@it-yaroslav.ru"
	//password := "Qwertycv13"
	//
	//to := []string{"yaroslav215@icloud.com"}
	//host := "smtp.yandex.ru"
	//port := "465"
	//
	//message := []byte("This is a test email message.")
	//
	//auth := smtp.PlainAuth("", from, password, host)
	//err := smtp.SendMail(host+":"+port, auth, from, to, message)
	//if err != nil {
	//	logrus.Info(err)
	//}
	//
	//fmt.Println("hello")
	m := gomail.NewMessage()
	m.SetHeader("From", "iam@it-yaroslav.ru")
	m.SetHeader("To", "yaroslav215@icloud.com")
	m.SetHeader("Subject", "asfasfasdf")
	m.SetBody("text/plain", "sfasf34534534553")

	d := gomail.NewDialer("smtp.yandex.ru", 465, "iam@it-yaroslav.ru", "Qwertycv13")

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
