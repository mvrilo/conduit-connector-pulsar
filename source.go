// Copyright © 2022 Meroxa, Inc.
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

package pulsar

import (
	"context"
	"fmt"

	"github.com/conduitio/conduit-connector-pulsar/source"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

const (
	// MetadataPulsarTopic is the metadata key for storing the kafka topic
	MetadataPulsarTopic = "kafka.topic"
)

type Source struct {
	sdk.UnimplementedSource

	consumer source.Consumer
	config   source.Config
}

func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	return source.Config{}.Parameters()
}

func (s *Source) Configure(_ context.Context, cfg map[string]string) error {
	var config source.Config

	err := sdk.Util.ParseConfig(cfg, &config)
	if err != nil {
		return err
	}
	err = config.Validate()
	if err != nil {
		return err
	}

	s.config = config
	return nil
}

func (s *Source) Open(ctx context.Context, _ sdk.Position) error {
	err := s.config.TryDial(ctx)
	if err != nil {
		return fmt.Errorf("failed to dial broker: %w", err)
	}

	if s.config.SubscriptionName == "" {
		// this must be the first run of the connector, create a new group ID
		s.config.SubscriptionName = "conduit-connector-pulsar"
		sdk.Logger(ctx).Info().Str("subscriptionName", s.config.SubscriptionName).Msg("assigning source subscription name")
	}

	s.consumer, err = source.NewPulsarConsumer(ctx, s.config)
	if err != nil {
		return fmt.Errorf("failed to create Pulsar consumer: %w", err)
	}

	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	rec, err := s.consumer.Consume(ctx)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("failed getting a record: %w", err)
	}

	metadata := sdk.Metadata{MetadataPulsarTopic: rec.Topic()}
	metadata.SetCreatedAt(rec.EventTime())

	return sdk.Util.Source.NewRecordCreate(
		source.Position{
			// GroupID:   s.config.GroupID,
			// Offset:    rec.Offset(),
			Topic:     rec.Topic(),
			Partition: rec.ID().PartitionIdx(),
		}.ToSDKPosition(),
		metadata,
		sdk.RawData(rec.Key()),
		sdk.RawData(rec.Payload()),
	), nil
}

func (s *Source) Teardown(ctx context.Context) error {
	if s.consumer != nil {
		err := s.consumer.Close(ctx)
		if err != nil {
			return fmt.Errorf("failed closing Pulsar consumer: %w", err)
		}
	}
	return nil
}
