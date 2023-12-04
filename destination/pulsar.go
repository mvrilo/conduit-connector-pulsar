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

package destination

import (
	"context"
	"fmt"
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-sdk/kafkaconnect"
	"github.com/goccy/go-json"
)

type PulsarProducer struct {
	client     pulsar.Client
	producer   pulsar.Producer
	keyEncoder dataEncoder
}

var _ Producer = (*PulsarProducer)(nil)

func NewPulsarProducer(ctx context.Context, cfg Config) (*PulsarProducer, error) {
	opts := cfg.PulsarClientOpts(sdk.Logger(ctx))
	// opts = append(opts, []kgo.Opt{
	// 	kgo.AllowAutoTopicCreation(),
	// 	kgo.DefaultProduceTopic(cfg.Topic),
	// 	kgo.RecordDeliveryTimeout(cfg.DeliveryTimeout),
	// 	kgo.RequiredAcks(cfg.RequiredAcks()),
	// 	kgo.ProducerBatchMaxBytes(cfg.BatchBytes),
	// }...)

	// if cfg.RequiredAcks() != kgo.AllISRAcks() {
	// 	sdk.Logger(ctx).Warn().Msgf("disabling idempotent writes because \"acks\" is set to %v", cfg.Acks)
	// 	opts = append(opts, kgo.DisableIdempotentWrite())
	// }

	cli, err := pulsar.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create pulsar client: %w", err)
	}

	producer, err := cli.CreateProducer(pulsar.ProducerOptions{
		Topic: cfg.Topic,
	})
	if err != nil {
		return nil, err
	}

	var keyEncoder dataEncoder = bytesEncoder{}
	if cfg.useKafkaConnectKeyFormat {
		keyEncoder = kafkaConnectEncoder{}
	}

	return &PulsarProducer{
		client:     cli,
		keyEncoder: keyEncoder,
		producer:   producer,
	}, nil
}

func (p *PulsarProducer) Produce(ctx context.Context, records []sdk.Record) (int, error) {
	var wg sync.WaitGroup
	results := make([]Record, 0, len(records))

	sendCallback := func(msgid pulsar.MessageID, msg *pulsar.ProducerMessage, err error) {
		results = append(results, Record{MsgID: msgid, Msg: msg, Err: err})
		wg.Done()
	}

	wg.Add(len(records))
	for i, rec := range records {
		encodedKey, err := p.keyEncoder.Encode(rec.Key)
		if err != nil {
			return 0, fmt.Errorf("could not encode key of record %v: %w", i, err)
		}
		msg := &pulsar.ProducerMessage{
			Key:     string(encodedKey),
			Payload: rec.Bytes(),
		}
		p.producer.SendAsync(ctx, msg, sendCallback)
		wg.Done()
	}
	wg.Wait()

	for i, re := range results {
		if re.Err != nil {
			return i, re.Err
		}
	}

	return len(results), nil
}

func (p *PulsarProducer) Close(_ context.Context) error {
	if p.producer != nil {
		p.producer.Close()
	}
	if p.client != nil {
		p.client.Close()
	}
	return nil
}

// dataEncoder is similar to a sdk.Encoder, which takes data and encodes it in
// a certain format. The producer uses this to encode the key of the kafka
// message.
type dataEncoder interface {
	Encode(sdk.Data) ([]byte, error)
}

// bytesEncoder is a dataEncoder that simply calls data.Bytes().
type bytesEncoder struct{}

func (bytesEncoder) Encode(data sdk.Data) ([]byte, error) {
	return data.Bytes(), nil
}

// kafkaConnectEncoder encodes the data into a kafka connect JSON with schema
// (NB: this is not the same as JSONSchema).
type kafkaConnectEncoder struct{}

func (e kafkaConnectEncoder) Encode(data sdk.Data) ([]byte, error) {
	sd := e.toStructuredData(data)
	schema := kafkaconnect.Reflect(sd)
	if schema == nil {
		// s is nil, let's write an empty struct in the schema
		schema = &kafkaconnect.Schema{
			Type:     kafkaconnect.TypeStruct,
			Optional: true,
		}
	}

	env := kafkaconnect.Envelope{
		Schema:  *schema,
		Payload: sd,
	}

	// TODO add support for other encodings than JSON
	return json.Marshal(env)
}

// toStructuredData tries its best to return StructuredData.
func (kafkaConnectEncoder) toStructuredData(d sdk.Data) sdk.Data {
	switch d := d.(type) {
	case nil:
		return nil
	case sdk.StructuredData:
		return d
	case sdk.RawData:
		// try parsing the raw data as json
		var sd sdk.StructuredData
		err := json.Unmarshal(d, &sd)
		if err != nil {
			// it's not JSON, nothing more we can do
			return d
		}
		return sd
	default:
		// should not be possible
		panic(fmt.Errorf("unknown data type: %T", d))
	}
}
