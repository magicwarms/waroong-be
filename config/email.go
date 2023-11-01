package config

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"text/template"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

func SendEmail(messagesInfo []mailjet.InfoMessagesV31) string {
	mailjetClient := mailjet.NewMailjetClient(GoDotEnvVariable("MJ_APIKEY_PUBLIC"), GoDotEnvVariable("MJ_APIKEY_PRIVATE"))

	messages := mailjet.MessagesV31{Info: messagesInfo}
	response, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		fmt.Println("sendEmailError", err)
	}

	fmt.Println("SendEmailResponseStatus", response.ResultsV31[0].Status, "SendEmailResponseMessageId", response.ResultsV31[0].To[0].MessageID)

	return response.ResultsV31[0].Status
}

func ParseHtmlTemplate(templateFileName string, data interface{}) (string, error) {
	_, base, _, _ := runtime.Caller(0) // Relative to the runtime Dir
	dir := path.Join(path.Dir(base))
	rootDir := filepath.Dir(dir)

	// its better to use os specific path separator
	htmlDir := path.Join(rootDir, "config", "templates/"+templateFileName)

	t, err := template.ParseFiles(htmlDir)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
