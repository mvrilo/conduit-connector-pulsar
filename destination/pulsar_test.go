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
	"testing"

	"github.com/conduitio/conduit-connector-pulsar/common"
	"github.com/conduitio/conduit-connector-pulsar/test"
	"github.com/matryer/is"
)

func TestPulsarProducer_Opts(t *testing.T) {
	is := is.New(t)

	clientCert, clientKey, caCert := test.Certificates(t)

	cfg := Config{
		Config: common.Config{
			URL:   "pulsar://localhost:6650",
			Topic: "test-topic",
			ConfigTLS: common.ConfigTLS{
				ClientCert: clientCert,
				ClientKey:  clientKey,
				CACert:     caCert,
			},
		},
	}

	_, err := NewPulsarProducer(context.Background(), cfg)
	is.NoErr(err)
}

// func TestPulsarProducer_Opts_AcksDisableIdempotentWrite(t *testing.T) {
// 	// minimal valid config
// 	cfg := Config{
// 		Config:     common.Config{URL: "test-host:8080"},
// 		BatchBytes: 512,
// 	}

// 	testCases := []struct {
// 		acks                       string
// 		wantRequiredAcks           kgo.Acks
// 		wantDisableIdempotentWrite bool
// 	}{{
// 		acks:                       "all",
// 		wantRequiredAcks:           kgo.AllISRAcks(),
// 		wantDisableIdempotentWrite: false,
// 	}, {
// 		acks:                       "one",
// 		wantRequiredAcks:           kgo.LeaderAck(),
// 		wantDisableIdempotentWrite: true,
// 	}, {
// 		acks:                       "none",
// 		wantRequiredAcks:           kgo.NoAck(),
// 		wantDisableIdempotentWrite: true,
// 	}}

// 	for _, tc := range testCases {
// 		t.Run(tc.acks, func(t *testing.T) {
// 			is := is.New(t)
// 			cfg.Acks = tc.acks
// 			_, err := NewPulsarProducer(context.Background(), cfg)
// 			is.NoErr(err)

// 			// is.Equal(p.client.OptValue(kgo.RequiredAcks), tc.wantRequiredAcks)
// 			// is.Equal(p.client.OptValue(kgo.DisableIdempotentWrite), tc.wantDisableIdempotentWrite)
// 		})
// 	}
// }
