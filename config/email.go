package config

import (
	"bytes"
	"crypto/tls"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"text/template"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func EmailServer() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()

	parsePort, errParsePort := strconv.Atoi(GoDotEnvVariable("MAIL_PORT"))
	if errParsePort != nil {
		return nil, errParsePort
	}
	// SMTP Server
	server.Host = GoDotEnvVariable("MAIL_HOST")
	server.Port = parsePort
	server.Username = GoDotEnvVariable("MAIL_USERNAME")
	server.Password = GoDotEnvVariable("MAIL_PASSWORD")
	server.Encryption = mail.EncryptionSTARTTLS

	server.Authentication = mail.AuthAuto

	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, errSmtpClient := server.Connect()

	if errSmtpClient != nil {
		log.Fatal(errSmtpClient)
	}

	return smtpClient, nil

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
