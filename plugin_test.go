package main

import (
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
)

func TestMissingDefaultConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestMissingUserConfig(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Username: "123456789",
			Password: "123456789",
		},
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestDefaultMessageFormat(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number:  101,
			Status:  "success",
			Link:    "https://github.com/appleboy/go-hello",
			Author:  "Bo-Yi Wu",
			Branch:  "master",
			Message: "update travis",
			Commit:  "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},
	}

	message := plugin.Message(plugin.Repo, plugin.Build)

	assert.Equal(t, []string{"[success] <https://github.com/appleboy/go-hello> (master)『update travis』by Bo-Yi Wu"}, message)
}

func TestSendMessage(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number:  101,
			Status:  "success",
			Link:    "https://github.com/appleboy/go-hello",
			Author:  "Bo-Yi Wu",
			Email:   "test@gmail.com",
			Branch:  "master",
			Message: "update travis by drone plugin",
			Commit:  "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},

		Config: Config{
			Host:     os.Getenv("GOOGLE_HOST"),
			Username: os.Getenv("GOOGLE_USERNAME"),
			Password: os.Getenv("GOOGLE_PASSWORD"),
			To:       []string{os.Getenv("GOOGLE_TO"), "中文ID:a@gmail.com", "1234567890"},
			Message:  []string{"Test Google Chat Bot From Travis or Local", "commit message: 『{{ build.message }}』", " "},
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	// disable message
	plugin.Config.Message = []string{}
	err = plugin.Exec()
	assert.Nil(t, err)
}

func TestTrimElement(t *testing.T) {
	var input, result []string

	input = []string{"1", "     ", "3"}
	result = []string{"1", "3"}

	assert.Equal(t, result, trimElement(input))

	input = []string{"1", "2"}
	result = []string{"1", "2"}

	assert.Equal(t, result, trimElement(input))
}

func TestParseTo(t *testing.T) {
	id, enable := parseTo("appleboy@gmail.com", "test@gmail.com")

	assert.Equal(t, true, enable)
	assert.Equal(t, "appleboy@gmail.com", id)

	id, enable = parseTo("appleboy@gmail.com:test2@gmail.com", "test@gmail.com")

	assert.Equal(t, false, enable)
	assert.Equal(t, "", id)

	id, enable = parseTo("appleboy@gmail.com:test@gmail.com", "test@gmail.com")

	assert.Equal(t, true, enable)
	assert.Equal(t, "appleboy@gmail.com", id)
}
