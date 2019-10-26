package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version string

func main() {
	app := cli.NewApp()
	app.Name = "gtalk plugin"
	app.Usage = "gtalk plugin"
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Value:  "talk.google.com:443",
			Usage:  "gtalk host",
			EnvVar: "PLUGIN_HOST,GOOGLE_HOST",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "gtalk user email",
			EnvVar: "PLUGIN_USERNAME,GOOGLE_USERNAME",
		},
		cli.StringFlag{
			Name:   "oauthtoken",
			Usage:  "OAuthToken provides go-xmpp with the required OAuth2 token used to authenticate",
			EnvVar: "PLUGIN_OAUTH_TOKEN,OAUTH_TOKEN",
		},
		cli.StringFlag{
			Name:   "oauthscope",
			Usage:  "OAuthScope provides go-xmpp the required scope for OAuth2 authentication.",
			Value:  "https://www.googleapis.com/auth/googletalk",
			EnvVar: "PLUGIN_OAUTH_SCOPE,OAUTH_SCOPE",
		},
		cli.StringSliceFlag{
			Name:   "to",
			Usage:  "send message to user",
			EnvVar: "PLUGIN_TO,TO",
		},
		cli.StringSliceFlag{
			Name:   "message",
			Usage:  "gtalk message",
			EnvVar: "PLUGIN_MESSAGE",
		},
		cli.BoolFlag{
			Name:   "match.email",
			Usage:  "send message when only match email",
			EnvVar: "PLUGIN_ONLY_MATCH_EMAIL",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "repository owner and repository name",
			EnvVar: "DRONE_REPO,GITHUB_REPOSITORY",
		},
		cli.StringFlag{
			Name:   "repo.namespace",
			Usage:  "repository namespace",
			EnvVar: "DRONE_REPO_OWNER,DRONE_REPO_NAMESPACE,GITHUB_ACTOR",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA,GITHUB_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF,GITHUB_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.link",
			Usage:  "git commit link",
			EnvVar: "DRONE_COMMIT_LINK",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.author.avatar",
			Usage:  "git author avatar",
			EnvVar: "DRONE_COMMIT_AUTHOR_AVATAR",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.StringFlag{
			Name:   "pull.request",
			Usage:  "pull request",
			EnvVar: "DRONE_PULL_REQUEST",
		},
		cli.Float64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Float64Flag{
			Name:   "job.finished",
			Usage:  "job finished",
			EnvVar: "DRONE_BUILD_FINISHED",
		},
		cli.StringFlag{
			Name:   "env-file",
			Usage:  "source env file",
			Value:  ".env",
			EnvVar: "ENV_FILE",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "PLUGIN_DEBUG,DEBUG",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			FullName:  c.String("repo"),
			Namespace: c.String("repo.namespace"),
			Name:      c.String("repo.name"),
		},
		Commit: Commit{
			Sha:     c.String("commit.sha"),
			Ref:     c.String("commit.ref"),
			Branch:  c.String("commit.branch"),
			Link:    c.String("commit.link"),
			Author:  c.String("commit.author"),
			Email:   c.String("commit.author.email"),
			Avatar:  c.String("commit.author.avatar"),
			Message: c.String("commit.message"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Link:     c.String("build.link"),
			Started:  c.Float64("job.started"),
			Finished: c.Float64("job.finished"),
			PR:       c.String("pull.request"),
			DeployTo: c.String("deploy.to"),
		},
		Config: Config{
			Host:       c.String("host"),
			Username:   c.String("username"),
			To:         c.StringSlice("to"),
			Message:    c.StringSlice("message"),
			MatchEmail: c.Bool("match.email"),
			Debug:      c.Bool("debug"),
			OAuthScope: c.String("oauthscope"),
			OAuthToken: c.String("oauthtoken"),
		},
	}

	return plugin.Exec()
}
