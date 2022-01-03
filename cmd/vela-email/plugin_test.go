// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.
package main

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jordan-wright/email"
)

var (
	mockEmail = &email.Email{
		To: 	[]string{"one@mail.com"},
		From: 	"two@mail.com",
	}

	mockSMTPHost = &SmtpHost{
		Host:     "smtphost.com",
		Port:     "587",
		Username: "username",
		Password: "password",
	}	
	noAttachment = &email.Attachment{
		Filename: "",
	}

	mockBuildEnv = &BuildEnv{
		BuildCreated:   time.Unix(int64(1556720958), 0).UTC().String(),
		BuildEnqueued:  time.Unix(int64(1556720958), 0).UTC().String(),
		BuildFinished:  time.Unix(int64(1556720958), 0).UTC().String(),
		BuildStarted:   time.Unix(int64(1556720958), 0).UTC().String(),
	}

	mockPlugin = &Plugin{
		Email: mockEmail,
		SmtpHost: mockSMTPHost,
		Attachment: noAttachment,
		BuildEnv: mockBuildEnv,

	}

)

func createMockEnv() {
	os.Setenv("VELA_BUILD_CREATED", "1556720958")
	os.Setenv("VELA_BUILD_ENQUEUED", "1556720958")
	os.Setenv("VELA_BUILD_FINISHED", "1556720958")
	os.Setenv("VELA_BUILD_STARTED", "1556720958")
	os.Setenv("VELA_BUILD_AUTHOR", "octocat")
	os.Setenv("VELA_BUILD_AUTHOR_EMAIL", "octocat@github.com")
	os.Setenv("VELA_BUILD_BRANCH", "main")
	os.Setenv("VELA_BUILD_COMMIT", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	os.Setenv("VELA_BUILD_LINK", "https://vela-server.localhost/octocat/hello-world/1")
	os.Setenv("VELA_BUILD_MESSAGE", "Merge pull request #6 from octocat/patch-1")
	os.Setenv("VELA_BUILD_NUMBER", "1")
	os.Setenv("VELA_REPO_FULL_NAME", "octocat/hello-world")	
}

func TestValidateSuccess(t *testing.T) {
	tests := []struct {
		name 		string
		parameters 	Plugin
	}{
		{
			name: "return no errors: single To email",
			parameters: Plugin{
				Email: mockEmail,
				SmtpHost: mockSMTPHost,
				Attachment: noAttachment,
			},
		},
		{
			name: "return no errors: multiple To emails",
			parameters: Plugin{
				Email: &email.Email{
					To: 	[]string{"one@gmail.com", "two@comcast.net"},
					From: 	"three@email.com",
				},
				SmtpHost: mockSMTPHost,	
				Attachment: noAttachment,	
			},
		},
		{
			name: "return no errors: extra email parameters",
			parameters: Plugin{
				Email: &email.Email{
					To: 	 		[]string{"one@gmail.com", "two@comcast.net"},
					From: 	 		"three@email.com",
					ReplyTo: 		[]string{"first.last@email.com"},
					Bcc:	 		[]string{"first.last@email.com"},
					Cc: 	 		[]string{"first.last@email.com"},
					Subject: 		"subject",
					Text: 	 		[]byte(""),
					HTML: 	 		[]byte(""),
					Sender: 		"sender",
					ReadReceipt: 	[]string{"idk"},
				},
				SmtpHost: mockSMTPHost,
				Attachment: noAttachment,
			},
		},
		{
			name: "return no errors: parameters from attachment",
			parameters: Plugin{
				Email: &email.Email{
					To: 	[]string{""},
					From: 	"",
				},
				SmtpHost: mockSMTPHost,
				Attachment: &email.Attachment{
					Filename: "testdata/example1.txt",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.parameters.Validate(); err != nil {
				t.Errorf("Validate() should not have raised an error %s", err)
			}
		})
	}
}

func TestValidateErrors(t *testing.T) {
	tests := []struct {
		name 		string
		parameters 	Plugin
		wantErr 	error
	}{
		{
			name: "To missing",
			parameters: Plugin{
				Email: &email.Email{
					From: 	"two@email.com",
				},
				Attachment: noAttachment,
			},
			wantErr: ErrorMissingEmailToParam,
		},
		{
			name: "From missing",
			parameters: Plugin{
				Email: &email.Email{
					To: 	[]string{"one@email.com"},
				},
				Attachment: noAttachment,
			},
			wantErr: ErrorMissingEmailFromParam,
		},
		{
			name: "Email parameters missing from attachment",
			parameters: Plugin{
				Email: &email.Email{
					To: 	[]string{""},
					From: 	"",
				},
				SmtpHost: mockSMTPHost,
				Attachment: &email.Attachment{
					Filename: "testdata/badattachment.txt",
				},

			},
			wantErr: io.EOF,
		},
		{
			name: "Email attachment missing",
			parameters: Plugin{
				Attachment: &email.Attachment{
					Filename: "testdata/doesnotexist.txt",
				},
			},
			wantErr: os.ErrNotExist,
		},
		{
			name: "Email attachment empty",
			parameters: Plugin{
				Attachment: &email.Attachment{
					Filename: "testdata/empty.txt",
				},
			},
			wantErr: ErrorEmptyAttach,
		},
		{
			name: "SMTP host missing",
			parameters: Plugin{
				Email: mockEmail,
				SmtpHost: &SmtpHost{
					Port:     "1902",
				},
				Attachment: noAttachment,
			},
			wantErr: ErrorMissingSmtpParam,
		},
		{
			name: "SMTP port missing",
			parameters: Plugin{
				Email: mockEmail,
				SmtpHost: &SmtpHost{
					Host:     "smtphost.com",
				},
				Attachment: noAttachment,
			},
			wantErr: ErrorMissingSmtpParam,
		},
		{
			name: "SMTP username missing",
			parameters: Plugin{
				Email: mockEmail,
				SmtpHost: &SmtpHost{
					Host:     "smtphost.com",
					Port:     "1902",
				},
				Attachment: noAttachment,
			},
			wantErr: ErrorMissingSmtpUsernameParam,
		},
		{
			name: "SMTP password missing",
			parameters: Plugin{
				Email: mockEmail,
				SmtpHost: &SmtpHost{
					Host:     "smtphost.com",
					Port:     "1902",
					Username: "username",
				},
				Attachment: noAttachment,
			},
			wantErr: ErrorMissingSmtpPasswordParam,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			if err := test.parameters.Validate(); err == nil {
				t.Errorf("Validate() should have raised an error")
			}else if (err != test.wantErr) {
				t.Errorf("Validate() error = %v, wantErr = %v", err, test.wantErr)
			}
		})
	}
}

func TestInjectEnvSuccess(t *testing.T) {
	tests := []struct {
		name		string
		parameters	Plugin
	}{
		{
			name: "email using empty subject and html",
			parameters: *mockPlugin,
		},
		{
			name: "email using default subject and user text",
			parameters: Plugin {
				Email: &email.Email{
					To: 	 []string{"one@gmail.com", "two@comcast.net"},
					From: 	 "three@email.com",
					Subject:  DefaultSubject,
					Text:	  []byte("This is some text for repo: {{ .VELA_REPO_FULL_NAME }}"),
				},
				SmtpHost: mockSMTPHost,
				Attachment: noAttachment,
				BuildEnv: mockBuildEnv,
			},
		},
		{
			name: "email using user subject and html",
			parameters: Plugin {
				Email: &email.Email{
					To: 	 []string{"one@gmail.com", "two@comcast.net"},
					From: 	 "three@email.com",
					Subject:  "Commit failure on vela build: {{ .VELA_BUILD_NUMBER }}",
					Text:	  []byte("This is some text for repo: {{ .VELA_REPO_FULL_NAME }}"),
				},
				SmtpHost: mockSMTPHost,
				Attachment: noAttachment,
				BuildEnv: mockBuildEnv,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.parameters.Validate(); err != nil {
				t.Errorf("Validate() should not have raised an error: %s", err)
				t.FailNow()
			}
			createMockEnv()
			test.parameters.Env = test.parameters.Environment()	
			subject, err := test.parameters.injectEnv(test.parameters.Email.Subject)
			if err != nil {
				t.Errorf("InjectEnv(subject) should not have raised an error %s", err)
				t.FailNow()
			}
			if strings.Contains(subject, "<no value>") {
				t.Errorf("InjectEnv(subject) failed to inject all environment variables %s", subject)
			}
			
			var body string
			if len(test.parameters.Email.HTML) == 0 {
				body, err = test.parameters.injectEnv(string(test.parameters.Email.HTML))
			}else{
				body, err = test.parameters.injectEnv(string(test.parameters.Email.Text))
			}		
			if err != nil {
				t.Errorf("InjectEnv(body) should not have raised an error %s", err)
				t.FailNow()
			}
			if strings.Contains(body, "<no value>") {
				t.Errorf("InjectEnv(body) failed to inject all environment variables %s", body)
			}

		})
	}
}

func TestInjectEnvBadVar(t *testing.T) {

	tests := []struct {
		name		string
		parameters	Plugin
	}{
		{
			name: "error: using environment variable that doesnt exist",
			parameters: Plugin {
				Email: &email.Email{
					To: 	 []string{"one@gmail.com", "two@comcast.net"},
					From: 	 "three@email.com",
					Subject:  "This is a bad subject {{ .SOME_OTHER_VARIABLE }}",
				},
				SmtpHost: mockSMTPHost,
				Attachment: noAttachment,
				BuildEnv: mockBuildEnv,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.parameters.Validate(); err != nil {
				t.Errorf("Validate() should not have raised an error: %s", err)
				t.FailNow()
			}
			createMockEnv()
			os.Setenv("SOME_OTHER_VARIABLE", "check")
			test.parameters.Env = test.parameters.Environment()	
			subject, err := test.parameters.injectEnv(test.parameters.Email.Subject)

			if err != nil {
				t.Errorf("InjectEnv(subject) should not have raised an error %s", err)
				t.FailNow()
			}

			if strings.Contains(subject, "check") {
				t.Errorf("InjectEnv(subject) should not have have injected variable: SOME_OTHER_VARIABLE, but did")
			}

		})
	}
	

}