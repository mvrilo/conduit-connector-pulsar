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

package common

import (
	"context"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// Config contains common configuration parameters.
type Config struct {
	// URL the Apache Pulsar server address.
	URL string `json:"url" validate:"required"`
	// Topic is the Pulsar topic.
	Topic string `json:"topic" validate:"required"`

	ConfigTLS

	pulsarClientOptions *pulsar.ClientOptions
}

// Validate executes manual validations beyond what is defined in struct tags.
func (c Config) Validate() error {
	err := c.ConfigTLS.Validate()
	return err
}

// TryDial tries to establish a connection to brokers and returns nil if it
// succeeds to connect to at least one broker.
func (c Config) TryDial(_ context.Context) error {
	cl, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               c.URL,
		OperationTimeout:  1 * time.Second,
		ConnectionTimeout: 1 * time.Second,
	})
	if cl != nil {
		defer cl.Close()
	}
	if err != nil {
		return err
	}
	prod, err := cl.CreateProducer(pulsar.ProducerOptions{Topic: "test"})
	if err != nil {
		return err
	}
	prod.Close()
	return nil
}
