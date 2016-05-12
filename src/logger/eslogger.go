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
	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v3"
)

type ElasticLogger struct {
	client *elastic.Client
	conf   ElasticLoggerConfig
}

type ElasticLoggerConfig struct {
	Url   string
	Index string
	Host  string
}

func NewElasticLogger(econf ElasticLoggerConfig) (*ElasticLogger, error) {
	var el ElasticLogger

	el.conf = econf
	var err error
	el.client, err = elastic.NewClient(elastic.SetURL(el.conf.Url))
	if err != nil {
		return nil, errors.Wrap(err, "Elastic Search setup failed")
	}

	if el.conf.Host == "" {
		el.conf.Host = "localhost"
	}

	if el.conf.Index == "" {
		el.conf.Index = "eslacki"
	}

	// check if index exists
	exists, err := el.client.IndexExists(el.conf.Index).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to check if index %s exists",
			el.conf.Index)
	}

	if !exists {
		ci, err := el.client.CreateIndex(el.conf.Index).Do()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create index %s",
				el.conf.Index)
		}

		if ci.Acknowledged == false {
			return nil, errors.Wrapf(err, "failed to create index %s, no ack from the server",
				el.conf.Index)
		}
	}

	// should be good now
	return &el, nil
}

func (el *ElasticLogger) LogMessage(msg interface{}) error {

	_, err := el.client.
		Index().
		Index(el.conf.Index).
		Type("log").
		BodyJson(msg).
		Do()

	return err
}
