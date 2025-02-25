// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

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
	app := &cli.App{
		Name:      "vela-email",
		HelpName:  "vela-email",
		Usage:     "Vela Email plugin for sending Vela build information to a user's email.",
		Copyright: "Copyright 2022 Target Brands, Inc. All rights reserved.",
		Authors: []*cli.Author{
			{
				Name:  "Vela Admins",
				Email: "vela@target.com",
			},
		},
		Action:   run,
		Compiled: time.Now(),
		Version:  pluginVersion.Semantic(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				EnvVars:  []string{"PARAMETER_LOG_LEVEL", "EMAIL_LOG_LEVEL"},
				FilePath: "/vela/parameters/email/log_level,/vela/secrets/email/log_level",
				Name:     "log.level",
				Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
				Value:    "info",
			},
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_FROM", "EMAIL_FROM"},
				Name:    "from",
				Usage:   "from address",
			},
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_SENDER", "EMAIL_SENDER"},
				Name:    "sender",
				Usage:   "from address (overwrites from)",
			},
			&cli.StringSliceFlag{
				EnvVars: []string{"PARAMETER_REPLYTO", "EMAIL_REPLYTO"},
				Name:    "replyto",
				Usage:   "address to reply to",
			},
			&cli.StringSliceFlag{
				EnvVars: []string{"PARAMETER_TO", "EMAIL_TO"},
				Name:    "to",
				Usage:   "to addresses (supports more than one addresses)",
			},
			&cli.StringSliceFlag{
				EnvVars: []string{"PARAMETER_BCC", "EMAIL_CC"},
				Name:    "bcc",
				Usage:   "blind carbon copy to addresses (supports more than one addresses)",
			},
			&cli.StringSliceFlag{
				EnvVars: []string{"PARAMETER_CC", "EMAIL_CC"},
				Name:    "cc",
				Usage:   "carbon copy to addresses (supports more than one addresses",
			},
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_SUBJECT", "EMAIL_SUBJECT"},
				Name:    "subject",
				Usage:   "subject of email",
			},
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_TEXT", "EMAIL_TEXT"},
				Name:    "text",
				Usage:   "body of message in text format",
			},
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_HTML", "EMAIL_HTML"},
				Name:    "html",
				Usage:   "body of message in html format",
			},
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_READRECEIPT", "EMAIL_READRECEIPT"},
				Name:    "readreceipt",
				Usage:   "request read receipts and delivery notifications",
			},
			// Attachment flag
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_ATTACHMENT", "EMAIL_ATTACHMENT"},
				Name:    "attachment",
				Usage:   "file to attach to email",
			},
			// SmtpHost flags
			&cli.StringFlag{
				EnvVars:  []string{"PARAMETER_HOST", "EMAIL_HOST"},
				Name:     "host",
				Usage:    "smtp host",
				Required: true,
			},
			&cli.StringFlag{
				EnvVars:  []string{"PARAMETER_PORT", "EMAIL_PORT"},
				Name:     "port",
				Usage:    "smtp port",
				Required: true,
			},
			&cli.StringFlag{
				EnvVars:  []string{"PARAMETER_USERNAME", "USERNAME"},
				FilePath: "/vela/parameters/email/username,/vela/secrets/email/username",
				Name:     "username",
				Usage:    "smtp host username",
			},
			&cli.StringFlag{
				EnvVars:  []string{"PARAMETER_PASSWORD", "PASSWORD"},
				FilePath: "/vela/parameters/email/username,/vela/secrets/email/username",
				Name:     "password",
				Usage:    "smtp host password",
			},
			// EmailFilename flag
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_FILENAME", "EMAIL_FILENAME"},
				Name:    "filename",
				Usage:   "file that contains email information (To, From, Subject, etc.)",
			},
			// TLSConfig flag
			&cli.BoolFlag{
				EnvVars: []string{"PARAMETER_SKIPVERIFY", "EMAIL_SKIPVERIFY"},
				Name:    "skipverify",
				Usage:   "skip tls verify",
			},
			// SendType flag
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_SENDTYPE", "EMAIL_SENDTYPE"},
				Name:    "sendtype",
				Usage:   "send type options: (Plain|StartTLS|TLS) default is set to StartTLS",
				Value:   "StartTLS",
			},
			// Auth flag
			&cli.StringFlag{
				EnvVars: []string{"PARAMETER_AUTH", "EMAIL_AUTH"},
				Name:    "auth",
				Usage:   "authentication for login type (PlainAuth|LoginAuth) default is set to nil",
			},
			// Build Flags
			&cli.IntFlag{
				EnvVars: []string{"VELA_BUILD_CREATED", "BUILD_CREATED"},
				Name:    "build-created",
				Usage:   "environment variable reference for reading in build created",
			},
			&cli.IntFlag{
				EnvVars: []string{"VELA_BUILD_ENQUEUED", "BUILD_ENQUEUED"},
				Name:    "build-enqueued",
				Usage:   "environment variable reference for reading in build enqueued",
			},
			&cli.IntFlag{
				EnvVars: []string{"VELA_BUILD_FINISHED", "BUILD_FINISHED"},
				Name:    "build-finished",
				Usage:   "environment variable reference for reading in build finished",
			},
			&cli.IntFlag{
				EnvVars: []string{"VELA_BUILD_STARTED", "BUILD_STARTED"},
				Name:    "build-started",
				Usage:   "environment variable reference for reading in build started",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	switch c.String("log.level") {
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

	if c.IsSet("ci") {
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
		SendType: c.String("sendtype"),
		// auth configuration
		Auth: c.String("auth"),

		// email configuration
		Email: &email.Email{
			ReplyTo:     c.StringSlice("replyto"),
			From:        c.String("from"),
			To:          c.StringSlice("to"),
			Bcc:         c.StringSlice("bcc"),
			Cc:          c.StringSlice("cc"),
			Subject:     c.String("subject"),
			Text:        []byte(c.String("text")),
			HTML:        []byte(c.String("html")),
			Sender:      c.String("sender"),
			ReadReceipt: c.StringSlice("readreceipt"),
		},

		// email filename configuration
		EmailFilename: c.String("filename"),

		// attachment configuration
		Attachment: &email.Attachment{
			Filename: c.String("attachment"),
		},

		// smtp configuration
		SMTPHost: &SMTPHost{
			Host:     c.String("host"),
			Port:     c.String("port"),
			Username: c.String("username"),
			Password: c.String("password"),
		},

		// tls configuration
		TLSConfig: &tls.Config{
			ServerName:         c.String("host"),
			InsecureSkipVerify: c.Bool("skipverify"), //nolint:gosec // ignore false positive
		},

		// User Friendly Build configuration
		BuildEnv: &BuildEnv{
			BuildCreated:  time.Unix(int64(c.Int("build-created")), 0).UTC().String(),
			BuildEnqueued: time.Unix(int64(c.Int("build-enqueued")), 0).UTC().String(),
			BuildFinished: time.Unix(int64(c.Int("build-finished")), 0).UTC().String(),
			BuildStarted:  time.Unix(int64(c.Int("build-started")), 0).UTC().String(),
		},
	}

	// validates the plugin
	if err := p.Validate(); err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
