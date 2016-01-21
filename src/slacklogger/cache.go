package slacklogger

import (
	"fmt"
	"github.com/bobbytables/slacker"
)

type cache struct {
	users    map[string]string
	channels map[string]string
}

func newCache() cache {
	return cache{
		make(map[string]string),
		make(map[string]string),
	}
}

func (c *cache) updateChannels(sch []*slacker.Channel) {
	for _, channel := range sch {
		fmt.Printf("adding channel mapping: %s -> %s\n",
			channel.ID, channel.Name)
		c.channels[channel.ID] = channel.Name
	}

}

func (c *cache) updateUsers(su []*slacker.User) {
	for _, user := range su {
		fmt.Printf("adding user mapping: %s -> %s\n",
			user.ID, user.Name)
		c.users[user.ID] = user.Name
	}
}

func nameFromMap(entry string, m map[string]string) string {
	name, ok := m[entry]
	if !ok {
		return "(unknown)"
	}
	return name
}

func (c *cache) channelIDToName(id string) string {
	return nameFromMap(id, c.channels)
}

func (c *cache) userIDToName(id string) string {
	return nameFromMap(id, c.users)
}
