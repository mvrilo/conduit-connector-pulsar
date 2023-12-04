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

package source

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type PulsarConsumer struct {
	client   pulsar.Client
	consumer pulsar.Consumer
	acks     bool
}

var _ Consumer = (*PulsarConsumer)(nil)

func NewPulsarConsumer(ctx context.Context, cfg Config) (*PulsarConsumer, error) {
	opts := cfg.PulsarClientOpts(sdk.Logger(ctx))

	cl, err := pulsar.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create pulsar client: %w", err)
	}

	consumer, err := cl.Subscribe(pulsar.ConsumerOptions{
		Topic:            cfg.Topic,
		SubscriptionName: cfg.SubscriptionName,
		Type:             pulsar.Shared,
	})
	if err != nil {
		return nil, err
	}

	return &PulsarConsumer{
		client:   cl,
		consumer: consumer,
		acks:     cfg.Acks,
	}, nil
}

func (c *PulsarConsumer) Consume(_ context.Context) (*Record, error) {
	for {
		msg, err := c.consumer.Receive(context.TODO())
		if err != nil {
			return nil, err
		}
		if !c.acks {
			continue
		}
		if err := c.consumer.Ack(msg); err != nil {
			// TODO log error
			continue
		}
	}
}

func (c *PulsarConsumer) Ack(_ context.Context) error {
	// TODO address acks
	return nil
}

func (c *PulsarConsumer) Close(_ context.Context) error {
	err := c.consumer.Unsubscribe()
	c.consumer.Close()
	c.client.Close()
	return err
}
