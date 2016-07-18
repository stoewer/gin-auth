// Copyright (c) 2016, German Neuroinformatics Node (G-Node)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License. See
// LICENSE file in the root of the Project.

package util

import (
	"bytes"
	"fmt"
	"github.com/G-Node/gin-auth/conf"
	"net/smtp"
	"strconv"
	"text/template"
)

// EmailDispatcher defines an interface for e-mail dispatch.
type EmailDispatcher interface {
	Send(recipient []string, message []byte) error
}

type emailDispatcher struct {
	conf *conf.SmtpCredentials
	send func(string, smtp.Auth, string, []string, []byte) error
}

// Send sets up authentication for e-mail dispatch via smtp and invokes the objects send function.
func (e *emailDispatcher) Send(recipient []string, content []byte) error {
	addr := e.conf.Host + ":" + strconv.Itoa(e.conf.Port)
	auth := smtp.PlainAuth("", e.conf.From, e.conf.Password, e.conf.Host)
	return e.send(addr, auth, e.conf.From, recipient, content)
}

// NewEmailDispatcher returns an instance of emailDispatcher.
// Dependent on the value of config.smtp.Mode the send method will
// print the e-mail content to the commandline (value "print"), do nothing (value "skip")
// or by default send an e-mail via smtp.SendMail.
func NewEmailDispatcher() EmailDispatcher {
	conf := conf.GetSmtpCredentials()
	send := smtp.SendMail
	if conf.Mode == "print" {
		send = func(addr string, auth smtp.Auth, from string, recipient []string, cont []byte) error {
			fmt.Printf("E-Mail content:\n---\n%s---\n", string(cont))
			return nil
		}
	} else if conf.Mode == "skip" {
		send = func(addr string, auth smtp.Auth, from string, recipient []string, cont []byte) error {
			return nil
		}
	}
	return &emailDispatcher{conf, send}
}

// EmailStandardFields specifies all fields required for a standard format e-mail
type EmailStandardFields struct {
	From    string
	To      string
	Subject string
	Body    string
}

// MakePlainEmailTemplate returns a bytes.Buffer containing a standard format e-mail
func MakePlainEmailTemplate(content *EmailStandardFields) *bytes.Buffer {
	var doc bytes.Buffer

	tmpl, err := template.ParseFiles(conf.GetResourceFile("templates", "emailplain.txt"))
	if err != nil {
		panic("Error parsing e-mail template: " + err.Error())
	}
	err = tmpl.Execute(&doc, content)
	if err != nil {
		panic("Error executing e-mail template: " + err.Error())
	}
	return &doc
}

// MakeEmailTemplate returns a bytes.Buffer containing a multipart format e-mail
func MakeEmailTemplate(fileName string, content interface{}) *bytes.Buffer {
	var doc bytes.Buffer

	mainFile := conf.GetResourceFile("templates", "emaillayout.txt")
	contentFile := conf.GetResourceFile("templates", fileName)
	tmpl, err := template.ParseFiles(mainFile, contentFile)
	if err != nil {
		panic("Error parsing e-mail template: " + err.Error())
	}

	err = tmpl.Execute(&doc, content)
	if err != nil {
		panic("Error executing e-mail template: " + err.Error())
	}

	return &doc
}
