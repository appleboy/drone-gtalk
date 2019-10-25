package main

import (
	"crypto/tls"
	"errors"
	"fmt"
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
		Tag      string
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
		Host       string
		Username   string
		Password   string
		To         []string
		Message    []string
		MatchEmail bool
		Debug      bool
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

func serverName(host string) string {
	return strings.Split(host, ":")[0]
}

func parseTo(to []string, authorEmail string, matchEmail bool) []string {
	var emails []string
	var ids []string
	attachEmail := true

	for _, value := range trimElement(to) {
		idArray := trimElement(strings.Split(value, ":"))

		// check match author email
		if len(idArray) > 1 {
			if email := idArray[1]; email != authorEmail {
				continue
			}

			emails = append(emails, idArray[0])
			attachEmail = false
			continue
		}

		ids = append(ids, idArray[0])
	}

	if matchEmail == true && attachEmail == false {
		return emails
	}

	for _, value := range emails {
		ids = append(ids, value)
	}

	return ids
}

// Exec executes the plugin.
func (p Plugin) Exec() error {
	if len(p.Config.Host) == 0 || len(p.Config.Username) == 0 || len(p.Config.Password) == 0 || len(p.Config.To) == 0 {
		return errors.New("missing google config")
	}

	var message []string
	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build)
	}

	xmpp.DefaultConfig = tls.Config{
		ServerName:         serverName(p.Config.Host),
		InsecureSkipVerify: false,
	}

	options := xmpp.Options{
		Host:          p.Config.Host,
		User:          p.Config.Username,
		NoTLS:         false,
		Debug:         true,
		Session:       false,
		Status:        "xa",
		StatusMessage: "I for one welcome our new codebot overlords.",
	}

	talk, err := options.NewClient()

	if err != nil {
		return err
	}

	ids := parseTo(p.Config.To, p.Build.Email, p.Config.MatchEmail)

	// send message.
	for _, user := range ids {
		for _, value := range trimElement(message) {
			txt, err := template.RenderTrim(value, p)
			if err != nil {
				return err
			}

			if _, err := talk.Send(xmpp.Chat{Remote: user, Type: "chat", Text: txt}); err != nil {
				return err
			}
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
