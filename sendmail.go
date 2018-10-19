package sendmail

import (
	"encoding/json"
	"fmt"
	"github.com/go-mail/mail"
	"io"
)

type Config struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (c *Config) Update(n *Config) {
	if n.Server != "" {
		c.Server = n.Server
	}
	if n.Port != 0 {
		c.Port = n.Port
	}
	if n.User != "" {
		c.User = n.User
		c.Password = n.Password
	}
}

func NewConfig(r io.Reader) (*Config, error) {

	c := Config{
		Server:   "",
		Port:     465,
		User:     "",
		Password: "",
	}

	if err := json.NewDecoder(r).Decode(&c); err != nil {
		err = fmt.Errorf("Failed to decode JSON: %v", err)
		return &c, err
	}
	return &c, nil
}

func NewMail(from *string, to, cc []string, subject, body *string, html bool) *mail.Message {
	m := mail.NewMessage()

	m.SetHeader("From", *from)
	m.SetHeader("To", to...)
	if len(cc) != 0 {
		m.SetHeader("Cc", cc...)
	}
	m.SetHeader("Subject", *subject)

	if html == true {
		m.SetBody("text/html", *body)
	} else {
		m.SetBody("text/plain", *body)
	}
	return m
}

func Send(m *mail.Message, c *Config) error {
	var d *mail.Dialer

	if c.User == "" {
		d = &mail.Dialer{Host: c.Server, Port: c.Port}
	} else {
		d = mail.NewDialer(c.Server, c.Port, c.User, c.Password)
	}

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
