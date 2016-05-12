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
	"config"
	"flag"
	"importer"
	"io"
	"logger"
	"os"
)

var (
	// global local logger
	ll *logger.LocalLogger
)

func main() {
	var conffile string
	var debug bool
	var channel string
	var infile string
	flag.StringVar(&conffile, "config", "", "configuration file path")
	flag.StringVar(&channel, "channel", "", "channel name")
	flag.StringVar(&infile, "infile", "", "input file path")
	flag.BoolVar(&debug, "debug", false, "debug logging")
	flag.Parse()

	logger.SetupLocalLogger(logger.LocalLoggerConfig{
		Debug: debug,
	})
	ll := logger.NewLocalLogger()

	if conffile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if infile == "" {
		ll.Fatalf("need input file, use -infile\n")
	}

	if channel == "" {
		ll.Fatalf("need channel name, use -channel\n")
	}

	if err := config.Load(conffile); err != nil {
		ll.Fatalf("failed to load configuration: %s", err)
	}

	el, err := logger.NewElasticLogger(logger.ElasticLoggerConfig{
		Url:   config.C.Url,
		Index: config.C.Index,
	})
	if err != nil {
		ll.Fatalf("failed to setup elastic search logger: %s", err)
	}

	f, err := os.OpenFile(infile, os.O_RDONLY, 0)
	if err != nil {
		ll.Fatalf("failed to open source file: %s", err)
	}
	defer f.Close()

	wi := importer.NewWeechatImport(f, channel)
	for {
		m, err := wi.Next()
		if err == io.EOF {
			// all done
			break
		}

		if err != nil {
			ll.Fatalf("error processing input file: %s", err)
		}

		if m == nil {
			continue
		}

		ll.Debugf("message: %+v", m)
		if err := el.LogMessage(m); err != nil {
			ll.Errorf("failed to log message: %s", err)
		}
	}
}
