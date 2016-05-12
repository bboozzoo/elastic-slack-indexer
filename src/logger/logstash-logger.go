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
package logger

import (
	"encoding/json"
	"fmt"
	"github.com/heatxsink/go-logstash"
)

type LogstashLogger struct {
	ls *logstash.Logstash
}

func (l *LogstashLogger) Writeln(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("failed to encode to JSON: %s\n", err.Error())
		return
	}

	err = l.ls.Writeln(string(data))
	if err != nil {
		fmt.Printf("failed to send data to logstash: %s\n", err.Error())
		// try reconnecting
		l.ls.Connect()
	}

}

func NewLogstashLogger(host string, port int) *LogstashLogger {
	fmt.Printf("log to logstash at %s:%d\n", host, port)
	lsl := &LogstashLogger{
		logstash.New(host, port, 10),
	}

	_, err := lsl.ls.Connect()
	if err != nil {
		fmt.Printf("failed to connect to logstash: %s\n", err.Error())
		return nil
	}
	return lsl
}
