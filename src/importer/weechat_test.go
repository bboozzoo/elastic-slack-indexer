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
package importer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"slacklogger"
	"testing"
)

func collectMessages(t *testing.T, channel string, data string) []*slacklogger.Message {
	i := NewWeechatImport(bytes.NewBufferString(data),
		channel)
	messages := make([]*slacklogger.Message, 0, 10)

	for {
		m, err := i.Next()
		if err != nil {
			t.Logf("error: %s", err)
			break
		} else {
			if m != nil {
				messages = append(messages, m)
			}
		}
	}
	return messages
}

func TestMessage1(t *testing.T) {

	d := `MI 20150804T17:11:09Z 000 Your nickname is test
MI 20150804T17:11:09Z 000 Warning: This room is not anonymous. 
MR 20160219T11:07:00Z 000 <foobar>  i am a pony
`
	m := collectMessages(t, "testchan", d)
	assert.Len(t, m, 1)

	m0 := m[0]

	assert.Equal(t, "@foobar", m0.User)
	assert.Equal(t, "#testchan", m0.Channel)
	assert.Equal(t, "2016-02-19T11:07:00Z", m0.Timestamp)
}

func TestMessage2(t *testing.T) {

	d := `MR 20160219T11:07:00Z 000 <foobar>  i am a pony
`
	m := collectMessages(t, "#testchan", d)
	assert.Len(t, m, 1)

	m0 := m[0]

	assert.Equal(t, "@foobar", m0.User)
	assert.Equal(t, "#testchan", m0.Channel)
	assert.Equal(t, "2016-02-19T11:07:00Z", m0.Timestamp)
}
