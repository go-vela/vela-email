// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"

	"github.com/aymerick/douceur/inliner"
	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)

var (
	// ErrorMissingEmailToParam is returned when the plugin is missing the To email parameter.
	ErrorMissingEmailToParam = errors.New("missing email parameter: To")

	// ErrorMissingEmailFromParam is returned when the plugin is missing the From email parameter.
	ErrorMissingEmailFromParam = errors.New("missing email parameter: From")

	// ErrorEmptyAttach is returned when the plugin finds the provided attachment to be empty.
	ErrorEmptyAttach = errors.New("attachment provided is empty")

	// ErrorMissingSMTPParam is returned when the plugin is missing a smtp host or port parameter.
	ErrorMissingSMTPParam = errors.New("missing smtp parameter (host/port)")
)

// Plugin represents the configuration loaded for the plugin.
type (
	Plugin struct {
		// Email arguments loaded for the plugin
		Email *email.Email
		// Attachment arguments loaded for the plugin
		Attachment *email.Attachment
		// SmtpHost arguments loaded for the plugin
		SMTPHost *SMTPHost
		// TLSConfig arguments loaded for the plugin
		TLSConfig *tls.Config
		// SendType arguments loaded for the plugin
		SendType string
		// Auth arguments loaded for the plugin
		Auth string
		// Readable build time environment variables
		BuildEnv *BuildEnv
	}

	// SMTPHost struct.
	SMTPHost struct {
		Host     string
		Port     string
		Username string
		Password string
	}

	// User friendly readable Build Environment Variables.
	BuildEnv struct {
		BuildCreated  string
		BuildEnqueued string
		BuildFinished string
		BuildStarted  string
	}
)

// Validate checks the plugin parameters needed from the
// user are provided. If the email subject or HTML/text are
// not provided, defaults are set.
func (p *Plugin) Validate() error {
	logrus.Trace("entered plugin.Validate")
	defer logrus.Trace("exited plugin.Validate")

	logrus.Info("Validating Parameters...")

	if len(p.Attachment.Filename) != 0 {
		fileInfo, err := os.Stat(p.Attachment.Filename)
		if errors.Is(err, os.ErrNotExist) {
			return os.ErrNotExist
		}

		if fileInfo.Size() == 0 {
			return ErrorEmptyAttach
		}

		file, err := os.Open(p.Attachment.Filename)
		if err != nil {
			return err
		}

		p.Email, err = email.NewEmailFromReader(file)
		if err != nil {
			return err
		}

		if len(p.Email.To) > 0 {
			p.Email.To = stringToSlice(p.Email.To)
		}

		if len(p.Email.Cc) > 0 {
			p.Email.Cc = stringToSlice(p.Email.Cc)
		}

		if len(p.Email.Bcc) > 0 {
			p.Email.Bcc = stringToSlice(p.Email.Bcc)
		}
	}

	if len(p.Email.To) == 0 {
		return ErrorMissingEmailToParam
	}

	if len(p.Email.From) == 0 {
		return ErrorMissingEmailFromParam
	}

	if len(p.SMTPHost.Host) == 0 || len(p.SMTPHost.Port) == 0 {
		return ErrorMissingSMTPParam
	}

	// set defaults
	if len(p.Email.Subject) == 0 {
		p.Email.Subject = DefaultSubject
	}

	if len(p.Email.HTML) == 0 && len(p.Email.Text) == 0 {
		p.Email.HTML = []byte(DefaultHTMLBody)
	}

	return nil
}

// Creates an environment map for the plugin to use and adds
// any environment variables in the os environment as well as
// some user friendly build timestamps.
func (p *Plugin) Environment() map[string]string {
	logrus.Trace("entered plugin.Environment")
	defer logrus.Trace("exited plugin.Environment")

	logrus.Info("Setting up Environment...")

	envMap := map[string]string{}

	for _, v := range os.Environ() {
		splitV := strings.Split(v, "=")
		if strings.HasPrefix(splitV[0], "VELA_") {
			envMap[splitV[0]] = strings.Join(splitV[1:], "=")
		}
	}

	envMap["BuildCreated"] = p.BuildEnv.BuildCreated
	envMap["BuildEnqueued"] = p.BuildEnv.BuildEnqueued
	envMap["BuildFinished"] = p.BuildEnv.BuildFinished
	envMap["BuildStarted"] = p.BuildEnv.BuildStarted

	return envMap
}

// Parses subject and body of email to inject environment
// variables. Uses provided authentication type and send type
// to send the email.
func (p *Plugin) Exec() error {
	logrus.Trace("entered plugin.Execute")
	defer logrus.Trace("exited plugin.Execute")

	logrus.Debug("Parsing Subject...")

	subject, err := p.injectEnv(p.Email.Subject)
	if err != nil {
		return err
	}

	p.Email.Subject = subject

	if len(p.Email.HTML) > 0 {
		logrus.Debug("Parsing HTML...")

		body, err := p.injectEnv(string(p.Email.HTML))
		if err != nil {
			return err
		}

		logrus.Debug("Parsing CSS...")

		body, err = inliner.Inline(body)
		if err != nil {
			return err
		}

		p.Email.HTML = []byte(body)
	} else {
		logrus.Debug("Parsing Text...")

		body, err := p.injectEnv(string(p.Email.Text))
		if err != nil {
			return err
		}

		p.Email.Text = []byte(body)
	}

	var auth smtp.Auth

	switch strings.ToLower(p.Auth) {
	case "plainauth":
		logrus.Info("Using login authentication from smtp/PlainAuth...")

		auth = smtp.PlainAuth("", p.SMTPHost.Username, p.SMTPHost.Password, p.SMTPHost.Host)
	case "loginauth":
		logrus.Info("Using login authentication from loginauth/LoginAuth...")

		auth = LoginAuth(p.SMTPHost.Username, p.SMTPHost.Password)
	default:
		logrus.Info("Using no login authentication...")

		auth = nil
	}

	host := p.SMTPHost.Host + ":" + p.SMTPHost.Port

	switch strings.ToLower(p.SendType) {
	case "starttls":
		logrus.Info("Sending email with StartTLS...")

		if err := p.Email.SendWithStartTLS(host, auth, p.TLSConfig); err != nil {
			return fmt.Errorf("error sending with StartTLS: %w", err)
		}
	case "tls":
		logrus.Info("Sending email with TLS...")

		if err := p.Email.SendWithTLS(host, auth, p.TLSConfig); err != nil {
			return fmt.Errorf("error sending with TLS: %w", err)
		}
	case "plain":
		fallthrough
	default:
		logrus.Info("Sending email with Plain...")

		if err := p.Email.Send(host, auth); err != nil {
			return fmt.Errorf("error sending with Plain: %w", err)
		}
	}

	logrus.Info("Plugin finished")

	return nil
}

// Injects environment variables into email template.
func (p *Plugin) injectEnv(str string) (string, error) {
	logrus.Trace("entered plugin.InjectEnv")
	defer logrus.Trace("exited plugin.InjectEnv")

	// Inject to subject
	buffer := new(bytes.Buffer)

	// parse string to template
	t := template.Must(template.New("input").Parse(str))

	err := t.Execute(buffer, p.Environment())

	return buffer.String(), err
}

// splits a string of emails and returns them as a slice.
func stringToSlice(s []string) []string {
	var slice []string

	for _, e := range s {
		if len(e) > 0 {
			temp := strings.Split(e, ", ")
			slice = append(slice, temp...)
		}
	}

	return slice
}
