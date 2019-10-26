package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/drone/drone-template-lib/template"
	"github.com/mattn/go-xmpp"
)

type (
	// Repo information.
	Repo struct {
		FullName  string
		Namespace string
		Name      string
	}

	// Commit information.
	Commit struct {
		Sha     string
		Ref     string
		Branch  string
		Link    string
		Author  string
		Avatar  string
		Email   string
		Message string
	}

	// Build information.
	Build struct {
		Tag      string
		Event    string
		Number   int
		Status   string
		Link     string
		Started  float64
		Finished float64
		PR       string
		DeployTo string
	}

	// Config for the plugin.
	Config struct {
		Host       string
		Username   string
		To         []string
		Message    []string
		MatchEmail bool
		Debug      bool
		OAuthScope string
		OAuthToken string
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Commit Commit
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
	if len(p.Config.Username) == 0 ||
		len(p.Config.OAuthToken) == 0 ||
		len(p.Config.To) == 0 {
		return errors.New("missing config")
	}

	var message []string
	if len(p.Config.Message) > 0 {
		message = p.Config.Message
	} else {
		message = p.Message(p.Repo, p.Build, p.Commit)
	}

	xmpp.DefaultConfig = tls.Config{
		ServerName:         serverName(p.Config.Host),
		InsecureSkipVerify: false,
	}

	options := xmpp.Options{
		Host:          p.Config.Host,
		User:          p.Config.Username,
		NoTLS:         false,
		Debug:         p.Config.Debug,
		Session:       false,
		Status:        "xa",
		OAuthToken:    p.Config.OAuthToken,
		OAuthScope:    p.Config.OAuthScope,
		StatusMessage: "I for one welcome our new codebot overlords.",
	}

	talk, err := options.NewClient()

	if err != nil {
		return err
	}

	ids := parseTo(p.Config.To, p.Commit.Email, p.Config.MatchEmail)

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
func (p Plugin) Message(repo Repo, build Build, commit Commit) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		build.Status,
		build.Link,
		commit.Branch,
		commit.Message,
		commit.Author,
	)}
}
