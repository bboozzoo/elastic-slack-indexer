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
package main

import (
	"flag"
	"logger"
	"os"
	"slacklogger"
)

var (
	// global local logger
	ll *logger.LocalLogger
)

func main() {
	var config string
	var debug bool
	flag.StringVar(&config, "config", "", "configuration file path")
	flag.BoolVar(&debug, "debug", false, "debug logging")
	flag.Parse()

	if config == "" {
		flag.Usage()
		os.Exit(1)
	}

	logger.SetupLocalLogger(logger.LocalLoggerConfig{
		Debug: debug,
	})

	ll := logger.NewLocalLogger()

	err := LoadConfig(config)
	if err != nil {
		ll.Errorf("failed to load config: %s", err)
		os.Exit(1)
	}

	err = logger.SetupElasticLogger(logger.ElasticLoggerConfig{
		Url:   C.Url,
		Index: C.Index,
	})
	if err != nil {
		ll.Errorf("failed to setup elastic search logger: %s", err)
		os.Exit(1)
	}

	sl := slacklogger.New(C.Token)

	sl.UpdateCache()

	go func() {
		for {
			msg := sl.GetMessage()
			ll.Debugf("got msg: %s\n", msg)
			ll.Debug(msg)
		}
	}()

	sl.HandleMessages()
}
