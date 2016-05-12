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
