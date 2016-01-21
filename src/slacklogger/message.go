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
