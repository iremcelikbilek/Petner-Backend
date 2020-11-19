package MailSender

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

func MailSender(subject string, code string, to string) {
	from := "app.petner@gmail.com"
	password := os.Getenv("SENDINBLUEPASS")
	header := make(map[string]string)
	header["From"] = "app.petner@gmail.com"
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	var mailBody MailBodyModel = MailBodyModel{code}
	t, fileError := template.ParseFiles("./Modules/MailSender/mailIndex.html")

	if fileError != nil {
		return
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, mailBody); err != nil {
		return
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(tpl.String()))
	auth := smtp.PlainAuth("", from, password, "smtp-relay.sendinblue.com")
	err := smtp.SendMail("smtp-relay.sendinblue.com:587", auth, from, []string{to}, []byte(message))

	if err != nil {
		fmt.Println(err)
		return
	}
}
