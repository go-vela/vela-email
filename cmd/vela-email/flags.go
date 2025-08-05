// SPDX-License-Identifier: Apache-2.0

package main

import "github.com/urfave/cli/v3"

// flags returns the CLI flags for the application.
func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "log.level",
			Value: "info",
			Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LOG_LEVEL"),
				cli.EnvVar("EMAIL_LOG_LEVEL"),
				cli.File("/vela/parameters/email/log_level"),
				cli.File("/vela/secrets/email/log_level"),
			),
		},
		&cli.StringFlag{
			Name:    "from",
			Usage:   "from address",
			Sources: cli.EnvVars("PARAMETER_FROM", "EMAIL_FROM"),
		},
		&cli.StringFlag{
			Name:    "sender",
			Usage:   "from address (overwrites from)",
			Sources: cli.EnvVars("PARAMETER_SENDER", "EMAIL_SENDER"),
		},
		&cli.StringSliceFlag{
			Name:    "replyto",
			Usage:   "address to reply to",
			Sources: cli.EnvVars("PARAMETER_REPLYTO", "EMAIL_REPLYTO"),
		},
		&cli.StringSliceFlag{
			Name:    "to",
			Usage:   "to addresses (supports more than one addresses)",
			Sources: cli.EnvVars("PARAMETER_TO", "EMAIL_TO"),
		},
		&cli.StringSliceFlag{
			Name:    "bcc",
			Usage:   "blind carbon copy to addresses (supports more than one addresses)",
			Sources: cli.EnvVars("PARAMETER_BCC", "EMAIL_CC"),
		},
		&cli.StringSliceFlag{
			Name:    "cc",
			Usage:   "carbon copy to addresses (supports more than one addresses",
			Sources: cli.EnvVars("PARAMETER_CC", "EMAIL_CC"),
		},
		&cli.StringFlag{
			Name:    "subject",
			Usage:   "subject of email",
			Sources: cli.EnvVars("PARAMETER_SUBJECT", "EMAIL_SUBJECT"),
		},
		&cli.StringFlag{
			Name:    "text",
			Usage:   "body of message in text format",
			Sources: cli.EnvVars("PARAMETER_TEXT", "EMAIL_TEXT"),
		},
		&cli.StringFlag{
			Name:    "html",
			Usage:   "body of message in html format",
			Sources: cli.EnvVars("PARAMETER_HTML", "EMAIL_HTML"),
		},
		&cli.StringFlag{
			Name:    "readreceipt",
			Usage:   "request read receipts and delivery notifications",
			Sources: cli.EnvVars("PARAMETER_READRECEIPT", "EMAIL_READRECEIPT"),
		},
		// Attachment flag
		&cli.StringFlag{
			Name:    "attachment",
			Usage:   "file to attach to email",
			Sources: cli.EnvVars("PARAMETER_ATTACHMENT", "EMAIL_ATTACHMENT"),
		},
		// SmtpHost flags
		&cli.StringFlag{
			Name:     "host",
			Usage:    "smtp host",
			Required: true,
			Sources:  cli.EnvVars("PARAMETER_HOST", "EMAIL_HOST"),
		},
		&cli.StringFlag{
			Name:     "port",
			Usage:    "smtp port",
			Required: true,
			Sources:  cli.EnvVars("PARAMETER_PORT", "EMAIL_PORT"),
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "smtp host username",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_USERNAME"),
				cli.EnvVar("USERNAME"),
				cli.File("/vela/parameters/email/username"),
				cli.File("/vela/secrets/email/username"),
			),
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "smtp host password",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_PASSWORD"),
				cli.EnvVar("PASSWORD"),
				cli.File("/vela/parameters/email/username"),
				cli.File("/vela/secrets/email/username"),
			),
		},
		// EmailFilename flag
		&cli.StringFlag{
			Name:    "filename",
			Usage:   "file that contains email information (To, From, Subject, etc.)",
			Sources: cli.EnvVars("PARAMETER_FILENAME", "EMAIL_FILENAME"),
		},
		// TLSConfig flag
		&cli.BoolFlag{
			Name:    "skipverify",
			Usage:   "skip tls verify",
			Sources: cli.EnvVars("PARAMETER_SKIPVERIFY", "EMAIL_SKIPVERIFY"),
		},
		// SendType flag
		&cli.StringFlag{
			Name:    "sendtype",
			Value:   "StartTLS",
			Usage:   "send type options: (Plain|StartTLS|TLS) default is set to StartTLS",
			Sources: cli.EnvVars("PARAMETER_SENDTYPE", "EMAIL_SENDTYPE"),
		},
		// Auth flag
		&cli.StringFlag{
			Name:    "auth",
			Usage:   "authentication for login type (PlainAuth|LoginAuth) default is set to nil",
			Sources: cli.EnvVars("PARAMETER_AUTH", "EMAIL_AUTH"),
		},
		// Build Flags
		&cli.IntFlag{
			Name:    "build-created",
			Usage:   "environment variable reference for reading in build created",
			Sources: cli.EnvVars("VELA_BUILD_CREATED", "BUILD_CREATED"),
		},
		&cli.IntFlag{
			Name:    "build-enqueued",
			Usage:   "environment variable reference for reading in build enqueued",
			Sources: cli.EnvVars("VELA_BUILD_ENQUEUED", "BUILD_ENQUEUED"),
		},
		&cli.IntFlag{
			Name:    "build-finished",
			Usage:   "environment variable reference for reading in build finished",
			Sources: cli.EnvVars("VELA_BUILD_FINISHED", "BUILD_FINISHED"),
		},
		&cli.IntFlag{
			Name:    "build-started",
			Usage:   "environment variable reference for reading in build started",
			Sources: cli.EnvVars("VELA_BUILD_STARTED", "BUILD_STARTED"),
		},
	}
}
