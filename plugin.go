package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/appleboy/drone-facebook/template"
	"github.com/mattn/go-xmpp"
)

type (
	// Repo information.
	Repo struct {
		Owner string
		Name  string
	}

	// Build information.
	Build struct {
		Event    string
		Number   int
		Commit   string
		Message  string
		Branch   string
		Author   string
		Status   string
		Link     string
		Started  float64
		Finished float64
	}

	// Config for the plugin.
	Config struct {
		Username string
		Password string
		To       []string
		Message  []string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

func trimElement(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.Trim(value, " ")
		if len(value) == 0 {
			continue
		}
		newKeys = append(newKeys, value)
	}

	return newKeys
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.Username) == 0 || len(p.Config.Password) == 0 || len(p.Config.To) == 0 {
		log.Println("missing google config")

		return errors.New("missing google config")
	}

	var message []string
	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	xmpp.DefaultConfig = tls.Config{
		ServerName:         "talk.google.com",
		InsecureSkipVerify: true,
	}

	options := xmpp.Options{
		Host:          "talk.google.com:443",
		User:          p.Config.Username,
		Password:      p.Config.Password,
		NoTLS:         false,
		Debug:         false,
		Session:       false,
		Status:        "xa",
		StatusMessage: "I for one welcome our new codebot overlords.",
	}

	talk, err := options.NewClient()

	if err != nil {
		log.Println(err.Error())

		return err
	}

	// send message.
	for _, to := range trimElement(p.Config.To) {
		for _, value := range trimElement(message) {
			txt, err := template.RenderTrim(value, p)
			if err != nil {
				return err
			}

			talk.Send(xmpp.Chat{Remote: to, Type: "chat", Text: txt})
		}
	}

	return nil
}

// Message is plugin default message.
func (p Plugin) Message(repo Repo, build Build) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		build.Status,
		build.Link,
		build.Branch,
		build.Message,
		build.Author,
	)}
}
