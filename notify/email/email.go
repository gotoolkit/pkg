package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"

	"github.com/pkg/errors"
)

// Options represents mail client options
type Options struct {
	From, To               string
	Subject                string
	Server, User, Password string
	Port                   int
}

// Client implements the notify.Sender
type Client struct {
	opt Options
}

// New returns a new mail client
func New(opt Options) *Client {
	return &Client{opt: opt}
}

// Send implements the notify.Sender interface
func (c *Client) Send(body []byte) error {
	auth := smtp.PlainAuth("", c.opt.User, c.opt.Password, c.opt.Server)
	msg, err := c.buildMessage(body)
	if err != nil {
		return errors.Wrapf(err, "failed to build message")
	}
	addr := c.parseAddr()
	err = smtp.SendMail(addr, auth, c.opt.From, []string{c.opt.To}, msg)
	if err != nil {
		return errors.Wrapf(err, "failed to send email")
	}
	return nil
}
func (c *Client) parseAddr() string {
	port := 25
	if c.opt.Port > 0 {
		port = c.opt.Port
	}
	return fmt.Sprintf("%s:%d", c.opt.Server, port)
}
func (c *Client) buildMessage(body []byte) ([]byte, error) {
	header, err := c.buildHeader()
	if err != nil {
		return nil, err
	}
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(body)
	return []byte(message), nil
}

func (c *Client) buildHeader() (map[string]string, error) {
	from, err := mail.ParseAddress(c.opt.From)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build header from address %s", from)
	}
	to, err := mail.ParseAddressList(c.opt.To)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build header to addresses %s", from)
	}
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = parseAddressToString(to)
	header["Subject"] = c.opt.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
	return header, nil
}
func parseMailSlice(mails []*mail.Address) []string {
	var mailSlice []string
	for _, email := range mails {
		mailSlice = append(mailSlice, email.Address)
	}
	return mailSlice
}

func parseAddressToString(mails []*mail.Address) string {
	var mailBuffer bytes.Buffer
	for _, email := range mails {
		mailBuffer.WriteString(email.String())
		mailBuffer.WriteString(",")
	}
	return mailBuffer.String()
}
