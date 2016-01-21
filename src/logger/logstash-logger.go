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
