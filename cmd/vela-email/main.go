// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"time"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/go-vela/vela-email/version"
)

func main() {
	// capture application version information.
	pluginVersion := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(pluginVersion, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	cmd := &cli.Command{
		Name:      "vela-email",
		Usage:     "Vela Email plugin for sending Vela build information to a user's email.",
		Copyright: "Copyright 2022 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		Action:  run,
		Version: pluginVersion.Semantic(),
		Flags:   flags(),
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(_ context.Context, cmd *cli.Command) error {
	switch cmd.String("log.level") {
	case "t", "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	if cmd.IsSet("ci") {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: false,
			PadLevelText:  true,
		})
	}

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-email",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/email/",
		"registry": "https://hub.docker.com/r/target/vela-email",
	}).Info("Vela Email Plugin")

	// create the plugin
	p := &Plugin{
		// sendType configuration
		SendType: cmd.String("sendtype"),
		// auth configuration
		Auth: cmd.String("auth"),

		// email configuration
		Email: &email.Email{
			ReplyTo:     cmd.StringSlice("replyto"),
			From:        cmd.String("from"),
			To:          cmd.StringSlice("to"),
			Bcc:         cmd.StringSlice("bcc"),
			Cc:          cmd.StringSlice("cc"),
			Subject:     cmd.String("subject"),
			Text:        []byte(cmd.String("text")),
			HTML:        []byte(cmd.String("html")),
			Sender:      cmd.String("sender"),
			ReadReceipt: cmd.StringSlice("readreceipt"),
		},

		// email filename configuration
		EmailFilename: cmd.String("filename"),

		// attachment configuration
		Attachment: &email.Attachment{
			Filename: cmd.String("attachment"),
		},

		// smtp configuration
		SMTPHost: &SMTPHost{
			Host:     cmd.String("host"),
			Port:     cmd.String("port"),
			Username: cmd.String("username"),
			Password: cmd.String("password"),
		},

		// tls configuration
		TLSConfig: &tls.Config{
			ServerName:         cmd.String("host"),
			InsecureSkipVerify: cmd.Bool("skipverify"), //nolint:gosec // ignore false positive
		},

		// User Friendly Build configuration
		BuildEnv: &BuildEnv{
			BuildCreated:  time.Unix(int64(cmd.Int("build-created")), 0).UTC().String(),
			BuildEnqueued: time.Unix(int64(cmd.Int("build-enqueued")), 0).UTC().String(),
			BuildFinished: time.Unix(int64(cmd.Int("build-finished")), 0).UTC().String(),
			BuildStarted:  time.Unix(int64(cmd.Int("build-started")), 0).UTC().String(),
		},
	}

	// validates the plugin
	if err := p.Validate(); err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
