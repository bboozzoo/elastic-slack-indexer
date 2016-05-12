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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"logger"
	"os"
	"slacklogger"
)

type config struct {
	Token string `json:"token"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
}

func configFromFile(path string) (*config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open config file %s: %s", path, err.Error()))
	}

	fmt.Printf("config data: %s\n", data)

	conf := &config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to parse config: %s", err.Error()))
	}
	fmt.Printf("config: %s\n", conf)
	return conf, nil
}

func main() {
	var config string
	flag.StringVar(&config, "config", "", "configuration file path")
	flag.Parse()

	if config == "" {
		flag.Usage()
		os.Exit(1)
	}

	conf, err := LoadConfig(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	log := logger.NewLogstashLogger(conf.Host, conf.Port)
	if log == nil {
		panic("failed to setup logger")
	}
	sl := slacklogger.New(conf.Token)

	sl.UpdateCache()

	go func() {
		for {
			msg := sl.GetMessage()
			fmt.Printf("got msg: %s\n", msg)
			log.Writeln(msg)
		}
	}()

	sl.HandleMessages()
}
