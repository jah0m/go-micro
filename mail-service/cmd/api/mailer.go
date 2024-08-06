package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Mail represents an email service
type Mail struct {
	Domain      string `json:"domain"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Encryption  string `json:"encryption"`
	FromAddress string `json:"from_address"`
	FromName    string `json:"from_name"`
}

// Message represents an email message
type Message struct {
	From        string         `json:"from"`
	FromName    string         `json:"from_name"`
	To          string         `json:"to"`
	Subject     string         `json:"subject"`
	Attachments []string       `json:"attachments"`
	Data        any            `json:"data"`
	DataMap     map[string]any `json:"data_map"`
}

// SendSMTPMessage sends an email message using SMTP
func (m *Mail) SendSMTPMessage(msg Message) error {
	// Set the default from address and name if not provided
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainMessage(msg)
	if err != nil {
		return err
	}

	// Create a new SMTP client
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryptionType(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Connect to the SMTP server
	client, err := server.Connect()
	if err != nil {
		return err
	}

	// Create a new email message
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	// Send the email message
	err = email.Send(client)
	if err != nil {
		return err
	}

	return nil
}

// buildHTMLMessage builds an HTML email message
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}
	return formattedMessage, nil
}

// buildPlainMessage builds a plain text email message
func (m *Mail) buildPlainMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()
	return plainMessage, nil
}

// inlineCSS inlines CSS styles into an HTML email message
func (m *Mail) inlineCSS(s string) (string, error) {
	options := &premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}
	p, err := premailer.NewPremailerFromString(s, options)
	if err != nil {
		return "", err
	}

	html, err := p.Transform()
	if err != nil {
		return "", err
	}
	return html, nil
}

func (m *Mail) getEncryptionType(encryption string) mail.Encryption {
	switch encryption {
	case "ssl":
		return mail.EncryptionSSL
	case "tls":
		return mail.EncryptionTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
