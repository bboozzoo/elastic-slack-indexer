package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
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

const (
	usage = "" +
		`Usage of 'slack-logger' command:

  slack-logger <config-file>
`
)

func main() {

	flag.Usage = func() {
		fmt.Println(usage)
		os.Exit(1)
	}

	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}

	conf, err := configFromFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	sl := slacklogger.New(conf.Token)

	sl.UpdateCache()
	sl.HandleMessages()
}
