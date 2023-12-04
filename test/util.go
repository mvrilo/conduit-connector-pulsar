// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/conduitio/conduit-connector-pulsar/common"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

// timeout is the default timeout used in tests when interacting with pulsar.
const timeout = 5 * time.Second

// T reports when failures occur.
// testing.T and testing.B implement this interface.
type T interface {
	// Fail indicates that the test has failed but
	// allowed execution to continue.
	Fail()
	// FailNow indicates that the test has failed and
	// aborts the test.
	FailNow()
	// Name returns the name of the running (sub-) test or benchmark.
	Name() string
	// Logf formats its arguments according to the format, analogous to Printf, and
	// records the text in the error log.
	Logf(string, ...any)
	// Cleanup registers a function to be called when the test (or subtest) and all its
	// subtests complete. Cleanup functions will be called in last added,
	// first called order.
	Cleanup(func())
}

func ConfigMap(t T) map[string]string {
	lastSlash := strings.LastIndex(t.Name(), "/")
	topic := t.Name()[lastSlash+1:] + uuid.NewString()
	t.Logf("using topic: %v", topic)
	return map[string]string{
		"url":   "pulsar://localhost:6650",
		"topic": topic,
	}
}

func SourceConfigMap(t T) map[string]string {
	m := ConfigMap(t)
	m["readFromBeginning"] = "true"
	return m
}

func DestinationConfigMap(t T) map[string]string {
	m := ConfigMap(t)
	return m
}

func ParseConfigMap[C any](t T, cfg map[string]string) C {
	is := is.New(t)
	is.Helper()

	var out C
	err := sdk.Util.ParseConfig(cfg, &out)
	is.NoErr(err)

	return out
}

func NewClient(t T, cfg common.Config) pulsar.Client {
	cli, err := pulsar.NewClient(
		pulsar.ClientOptions{
			URL: cfg.URL,
		},
	)
	is := is.New(t)
	is.NoErr(err)
	return cli
}

func Consume(t T, cfg common.Config, limit int) []pulsar.Message {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	is := is.New(t)
	is.Helper()

	cl, err := pulsar.NewClient(
		pulsar.ClientOptions{
			URL: cfg.URL,
		},
	)
	is.NoErr(err)
	defer cl.Close()

	consumer, err := cl.Subscribe(pulsar.ConsumerOptions{Topic: cfg.Topic})
	is.NoErr(err)
	defer consumer.Close()

	i := 0
	var records []pulsar.Message
	for i < limit {
		i++
		msg, err := consumer.Receive(ctx)
		is.NoErr(err)
		err = consumer.Ack(msg)
		is.NoErr(err)
		records = append(records, msg)
	}
	return records[:limit]
}

func Produce(t T, cfg common.Config, records []*pulsar.ProducerMessage, timeoutOpt ...time.Duration) {
	timeout := timeout // copy default timeout
	if len(timeoutOpt) > 0 {
		timeout = timeoutOpt[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	is := is.New(t)
	is.Helper()

	cl, err := pulsar.NewClient(
		pulsar.ClientOptions{
			URL: cfg.URL,
		},
	)
	is.NoErr(err)
	if cl != nil {
		defer cl.Close()
	}

	producer, err := cl.CreateProducer(pulsar.ProducerOptions{
		Topic: cfg.Topic,
	})
	is.NoErr(err)
	if producer != nil {
		defer producer.Close()
	}

	var results []*pulsar.ProducerMessage
	callback := func(msgid pulsar.MessageID, msg *pulsar.ProducerMessage, err error) {
		is.NoErr(err)
		results = append(results, msg)
	}

	for _, r := range records {
		producer.SendAsync(ctx, r, callback)
	}
}

func GeneratePulsarRecords(from, to int) []*pulsar.ProducerMessage {
	recs := make([]*pulsar.ProducerMessage, 0, to-from+1)
	for i := from; i <= to; i++ {
		recs = append(recs, &pulsar.ProducerMessage{
			Key:   fmt.Sprintf("test-key-%d", i),
			Value: fmt.Sprintf("test-payload-%d", i),
		})
	}
	return recs
}

func GenerateSDKRecords(from, to int) []sdk.Record {
	recs := GeneratePulsarRecords(from, to)
	sdkRecs := make([]sdk.Record, len(recs))
	for i, rec := range recs {
		metadata := sdk.Metadata{}
		// metadata.SetCreatedAt(rec.Timestamp)

		sdkRecs[i] = sdk.Util.Source.NewRecordCreate(
			[]byte(uuid.NewString()),
			metadata,
			sdk.RawData(rec.Key),
			sdk.RawData(rec.Payload),
		)
	}
	return sdkRecs
}

func Certificates(t T) (clientCert, clientKey, caCert string) {
	is := is.New(t)
	is.Helper()

	// get test dir
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled // we don't need other values
	testDir := path.Dir(filename)

	readFile := func(file string) string {
		bytes, err := os.ReadFile(path.Join(testDir, file))
		is.NoErr(err)
		return string(bytes)
	}

	clientCert = readFile("client.cer.pem")
	clientKey = readFile("client.key.pem")
	caCert = readFile("server.cer.pem")
	return
}

func ConfigWithIntegrationTestOptions(cfg common.Config) common.Config {
	return cfg.WithPulsarClientOpts(&pulsar.ClientOptions{
		URL: cfg.URL,
	})
}
