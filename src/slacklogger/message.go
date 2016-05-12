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
	"fmt"
	"github.com/bobbytables/slacker"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	userOrChannelRe = regexp.MustCompile("<[@#]([UC][A-Z0-9]+)>")
)

type LoggedMessage struct {
	User      string
	Channel   string
	Timestamp string
	Text      string
}

func msgTsToTime(ts string) time.Time {
	split := strings.Split(ts, ".")
	tmstamp, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		fmt.Printf("failed to parse timestamp %s, err: %s\n", ts, err.Error())
	}

	return time.Unix(tmstamp, 0)
}

func parseMsgText(text string, c *cache) string {
	found := userOrChannelRe.FindAllStringSubmatch(text, -1)

	for _, match := range found {
		if len(match[0]) < 2 {
			continue
		}

		var replace string
		// channels start with C, user IDs start with U
		if match[1][0] == 'U' {
			// user
			replace = "@" + c.userIDToName(string(match[1]))
		} else if match[1][0] == 'C' {
			// channel
			replace = "#" + c.channelIDToName(string(match[1]))
		}

		if len(replace) > 0 {
			text = strings.Replace(text, match[0], replace, -1)
		}

		fmt.Printf("replaced text: %s\n", text)
	}
	return text
}

func loggedMsgFromSlackMsg(smsg *slacker.RTMMessage, c *cache) *LoggedMessage {

	user := "@" + c.userIDToName(smsg.User)
	channel := "#" + c.channelIDToName(smsg.Channel)
	tm := msgTsToTime(smsg.Ts)

	text := parseMsgText(smsg.Text, c)

	return &LoggedMessage{
		user,
		channel,
		tm.String(),
		text,
	}
}
