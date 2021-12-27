package smtp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/buroz/grpc-clean-example/pkg/config"
)

type SMTPClient struct {
	config *config.SmtpConfig

	client *smtp.Client
	from   mail.Address
}

func NewSMTPClient(conf *config.SmtpConfig) SMTPClient {
	return SMTPClient{
		config: conf,
	}
}

func (c *SMTPClient) createConn(host, servername string, isSecure bool) (net.Conn, error) {
	if isSecure {
		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		return tls.Dial("tcp", servername, tlsconfig)
	}

	return net.Dial("tcp", servername)
}

func (c *SMTPClient) Connect() error {
	// c.Lock()
	// defer c.Unlock()

	servername := fmt.Sprintf("%v:%d", c.config.Host, c.config.Port)

	host, _, err := net.SplitHostPort(servername)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, host)

	conn, err := c.createConn(host, servername, c.config.Secure)
	if err != nil {
		return err
	}

	cl, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = cl.Auth(auth); err != nil {
		return err
	}

	c.client = cl
	c.from = mail.Address{
		Name:    "",
		Address: c.config.SenderEmail,
	}

	return nil
}

func (c *SMTPClient) Send(to string, subject string, body string) error {
	/*
		toArr := make([]mail.Address, len(to))

		for _, addr := range to {
			toArr = append(toArr, mail.Address{"", addr})
		}
	*/

	addr := mail.Address{
		Name:    "",
		Address: to,
	}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = c.from.String()
	headers["To"] = addr.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// To && From
	if err := c.client.Mail(c.from.Address); err != nil {
		return err
	}

	if err := c.client.Rcpt(addr.Address); err != nil {
		return err
	}

	// Data
	w, err := c.client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	// c.c.Quit()

	return nil
}
