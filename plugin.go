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
		Email    string
		Status   string
		Link     string
		Started  float64
		Finished float64
	}

	// Config for the plugin.
	Config struct {
		Host     string
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

func parseTo(value, authorEmail string) (string, bool) {
	ids := trimElement(strings.Split(value, ":"))

	if len(ids) > 1 {
		if email := ids[1]; email != authorEmail {
			log.Println("email not match")
			return "", false
		}
	}

	return ids[0], true
}

// Exec executes the plugin.
func (p Plugin) Exec() error {

	if len(p.Config.Host) == 0 || len(p.Config.Username) == 0 || len(p.Config.Password) == 0 || len(p.Config.To) == 0 {
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
		InsecureSkipVerify: true,
	}

	options := xmpp.Options{
		Host:          p.Config.Host,
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
		user, enable := parseTo(to, p.Build.Email)
		if !enable {
			continue
		}

		for _, value := range trimElement(message) {
			txt, err := template.RenderTrim(value, p)
			if err != nil {
				return err
			}

			talk.Send(xmpp.Chat{Remote: user, Type: "chat", Text: txt})
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
