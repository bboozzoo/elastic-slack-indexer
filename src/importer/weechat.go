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
	"bufio"
	"github.com/pkg/errors"
	"io"
	// "logger"
	"slacklogger"
	"strings"
	"time"
)

type WeechatImporter struct {
	bin     *bufio.Reader
	channel string
}

func NewWeechatImport(in io.Reader, channel string) Importer {
	if strings.HasPrefix(channel, "#") == false {
		channel = "#" + channel
	}
	return &WeechatImporter{
		bin:     bufio.NewReader(in),
		channel: channel,
	}
}

// Try to read the next data line from the input. If error != nil, the
// mesage is always nil. Message may also be nil if the current line
// did not contain an actual message, in this case Next() should be
// called again
func (w *WeechatImporter) Next() (*slacklogger.Message, error) {
	// l := logger.NewLocalLogger()

	line, err := w.bin.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			return nil, errors.Wrap(err, "failed to read input line")
		} else {
			return nil, err
		}
	}

	// line format:
	// MI 20150804T17:11:09Z 000 Your nickname is test
	// MR 20160219T11:07:00Z 000 <foobar>  i am a pony
	if len(line) < len("MR 20160219T11:07:00Z 000 ") {
		return nil, errors.Errorf("malformed line %s", line)
	}

	// we're only interested with lines with actual message
	// content, i.e. MR heaer
	if line[0:2] != "MR" {
		return nil, nil
	}

	tlen := len("20160219T11:07:00Z")
	tm := line[3 : 3+tlen]
	t, err := time.Parse("20060102T15:04:05Z", tm)
	if err != nil {
		return nil, errors.Wrapf(err,
			"failed to parse timestamp in line %s", line)
	}

	// left < of nickname
	nmin := strings.IndexByte(line, '<')
	if nmin == -1 {
		return nil, errors.Errorf("malformed nickname in line %s",
			line)
	}
	// l.Debugf("left index: %v", nmin)

	// right > of nickname
	nmax := strings.IndexByte(line[nmin:], '>')
	if nmax == -1 {
		return nil, errors.Errorf("malformed nicname end in line %s",
			line)
	}

	// l.Debugf("right index: %v", nmax)
	nick := line[nmin+1 : nmin+nmax]

	// l.Debugf("nickname: %s", nick)

	if len(nick) == 0 {
		return nil, errors.Errorf("nickname too short in line %s",
			line)
	}

	mstart := strings.IndexByte(line, ' ')
	if mstart == -1 {
		return nil, errors.Errorf("cannot locate message start in line %s",
			line)
	}

	msg := strings.TrimSpace(line[mstart+1:])

	if strings.HasPrefix(nick, "@") == false {
		nick = "@" + nick
	}

	sm := slacklogger.Message{
		User:      nick,
		Channel:   w.channel,
		Timestamp: t.UTC().Format(time.RFC3339Nano),
		Text:      msg,
	}

	// l.Debugf("slack message: %+v", sm)
	return &sm, nil
}
