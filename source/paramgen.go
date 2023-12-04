// Code generated by paramgen. DO NOT EDIT.
// Source: github.com/ConduitIO/conduit-connector-sdk/tree/main/cmd/paramgen

package source

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func (Config) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
		"acks": {
			Default:     "",
			Description: "acks defines a boolean for using acks",
			Type:        sdk.ParameterTypeBool,
			Validations: []sdk.Validation{},
		},
		"caCert": {
			Default:     "",
			Description: "caCert is the Pulsar broker's certificate.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"clientCert": {
			Default:     "",
			Description: "clientCert is the Pulsar client's certificate.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"clientKey": {
			Default:     "",
			Description: "clientKey is the Pulsar client's private key.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"insecureSkipVerify": {
			Default:     "",
			Description: "insecureSkipVerify defines whether to validate the broker's certificate chain and host name. If 'true', accepts any certificate presented by the server and any host name in that certificate.",
			Type:        sdk.ParameterTypeBool,
			Validations: []sdk.Validation{},
		},
		"readFromBeginning": {
			Default:     "",
			Description: "readFromBeginning determines from whence the consumer group should begin consuming when it finds a partition without a committed offset. If this options is set to true it will start with the first message in that partition.",
			Type:        sdk.ParameterTypeBool,
			Validations: []sdk.Validation{},
		},
		"subscriptionName": {
			Default:     "",
			Description: "subscriptionName defines a name for the subscription.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{},
		},
		"topic": {
			Default:     "",
			Description: "topic is the Pulsar topic.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
		"url": {
			Default:     "",
			Description: "url the Apache Pulsar server address.",
			Type:        sdk.ParameterTypeString,
			Validations: []sdk.Validation{
				sdk.ValidationRequired{},
			},
		},
	}
}