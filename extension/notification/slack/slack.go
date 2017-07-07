package slack

import (
	"os"
	"strings"
	"time"

	"github.com/nlopes/slack"

	"github.com/rusenask/keel/constants"
	"github.com/rusenask/keel/extension/notification"
	"github.com/rusenask/keel/types"

	log "github.com/Sirupsen/logrus"
)

const timeout = 5 * time.Second

type sender struct {
	slackClient *slack.Client
	channels    []string
	botName     string
}

func init() {
	notification.RegisterSender("slack", &sender{})
}

func (s *sender) Configure(config *notification.Config) (bool, error) {
	var token string
	// Get configuration
	if os.Getenv(constants.EnvSlackToken) != "" {
		token = os.Getenv(constants.EnvSlackToken)
	} else {
		return false, nil
	}
	if os.Getenv(constants.EnvSlackBotName) != "" {
		s.botName = os.Getenv(constants.EnvSlackBotName)
	} else {
		s.botName = "keel"
	}

	if os.Getenv(constants.EnvSlackChannels) != "" {
		channels := os.Getenv(constants.EnvSlackChannels)
		s.channels = strings.Split(channels, ",")
	} else {
		s.channels = []string{"general"}
	}

	s.slackClient = slack.New(token)

	log.WithFields(log.Fields{
		"name": "slack",
	}).Info("extension.notification.slack: sender configured")

	return true, nil
}

func (s *sender) Send(event types.EventNotification) error {
	params := slack.NewPostMessageParameters()
	params.Username = s.botName
	for _, channel := range s.channels {
		_, _, err := s.slackClient.PostMessage(channel, event.Message, params)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"channel": channel,
			}).Error("extension.notification.slack: failed to send notification")
		}
	}
	return nil
}
