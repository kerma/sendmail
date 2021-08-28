package sendmail

import (
	"encoding/json"
	"fmt"
	"io"
	"net/mail"

	gomail "github.com/go-mail/mail"
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

func NewMail(from *string, to, cc []*mail.Address, subject, body, attachment string, html bool) *gomail.Message {
	var toAddr, ccAddr []string
	m := gomail.NewMessage()

	m.SetHeader("From", *from)

	for _, t := range to {
		toAddr = append(toAddr, m.FormatAddress(t.Address, t.Name))
	}
	m.SetHeader("To", toAddr...)

	if len(cc) != 0 {
		for _, c := range to {
			ccAddr = append(toAddr, m.FormatAddress(c.Address, c.Name))
		}
		m.SetHeader("Cc", ccAddr...)
	}

	m.SetHeader("Subject", subject)

	if html == true {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	if attachment != "" {
		m.Attach(attachment)
	}

	return m
}

func Send(m *gomail.Message, c *Config) error {
	var d *gomail.Dialer

	if c.User == "" {
		d = &gomail.Dialer{Host: c.Server, Port: c.Port}
	} else {
		d = gomail.NewDialer(c.Server, c.Port, c.User, c.Password)
	}

	return d.DialAndSend(m)
}
