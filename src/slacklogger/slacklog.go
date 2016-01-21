package slacklogger

import (
	"fmt"
	"github.com/bobbytables/slacker"
)

type SlackLogger struct {
	client *slacker.APIClient
	scache cache
}

func (s *SlackLogger) updateChannelsCache() {
	sch, err := s.client.ChannelsList()
	if err != nil {
		panic(fmt.Sprintf("failed to fetch channels list: %s", err.Error()))
	}
	s.scache.updateChannels(sch)
}

func (s *SlackLogger) updateUsersCache() {
	su, err := s.client.UsersList()
	if err != nil {
		panic(fmt.Sprintf("failed to fetch users list: %s", err.Error()))
	}
	s.scache.updateUsers(su)
}

func New(token string) *SlackLogger {
	return &SlackLogger{
		slacker.NewAPIClient(token, ""),
		newCache(),
	}
}

func (s *SlackLogger) UpdateCache() {
	fmt.Println("cache update")
	s.updateChannelsCache()
	s.updateUsersCache()
}

func (s *SlackLogger) HandleMessages() {
	rtmStart, err := s.client.RTMStart()
	if err != nil {
		panic(err)
	}

	broker := slacker.NewRTMBroker(rtmStart)
	broker.Connect()

	for {
		event := <-broker.Events()

		if event.Type == "message" {
			msg, err := event.Message()
			if err != nil {
				panic(err)
			}

			lm := loggedMsgFromSlackMsg(msg, &s.scache)

			fmt.Printf("%s\n", lm)
		} else if event.Type == "channel_create" {
			fmt.Printf("new channel %s\n")
		} else if event.Type == "team_joined" {
			fmt.Printf("user joined a team\n")
		}
	}
}
