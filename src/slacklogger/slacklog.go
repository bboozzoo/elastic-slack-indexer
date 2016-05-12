// The MIT License (MIT)

// Copyright (c) 2016 Maciej Borzecki

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package slacklogger

import (
	"github.com/bobbytables/slacker"
	"logger"
)

type SlackLogger struct {
	client  *slacker.APIClient
	scache  cache
	msgPipe chan LoggedMessage
	l       *logger.LocalLogger
}

func (s *SlackLogger) updateChannelsCache() {
	sch, err := s.client.ChannelsList()
	if err != nil {
		s.l.Fatalf("failed to fetch channels list: %s", err.Error())
	}
	s.scache.updateChannels(sch)
}

func (s *SlackLogger) updateUsersCache() {
	su, err := s.client.UsersList()
	if err != nil {
		s.l.Fatalf("failed to fetch users list: %s", err.Error())
	}
	s.scache.updateUsers(su)
}

func New(token string) *SlackLogger {
	return &SlackLogger{
		slacker.NewAPIClient(token, ""),
		newCache(),
		make(chan LoggedMessage),
		logger.NewLocalLogger(),
	}
}

func (s *SlackLogger) UpdateCache() {
	s.l.Info("cache update")
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

			select {
			case s.msgPipe <- *lm:
			default:
				s.l.Debug(lm)
			}

		} else if event.Type == "channel_create" {
			s.l.Info("new channel %s")
		} else if event.Type == "team_joined" {
			s.l.Info("user joined a team")
		}
	}
}

func (s *SlackLogger) GetMessage() LoggedMessage {

	msg := <-s.msgPipe
	return msg
}
