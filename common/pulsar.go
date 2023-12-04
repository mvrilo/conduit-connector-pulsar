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
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/rs/zerolog"
)

// WithPulsarClientOpts lets you specify custom kafka client options (meant for
// test purposes).
func (c Config) WithPulsarClientOpts(opts *pulsar.ClientOptions) Config {
	c.pulsarClientOptions = opts
	return c
}

// PulsarClientOpts returns the kafka client options derived from the common config.
func (c Config) PulsarClientOpts(_ *zerolog.Logger) pulsar.ClientOptions {
	opts := pulsar.ClientOptions{}
	if tls := c.TLS(); tls != nil {
		// TODO
		opts.TLSAllowInsecureConnection = tls.InsecureSkipVerify
	}
	return opts
}
