# Conduit Connector Pulsar

The Apache Pulsar connector is one of [Conduit](https://github.com/ConduitIO/conduit) builtin plugins. It provides both, a
source and a destination connector for [Apache Pulsar](https://pulsar.apache.org).

## How to build?

Run `make build` to build the connector.

## Testing

Run `make test` to run all the unit and integration tests. Tests require Docker to be installed and running. The command
will handle starting and stopping docker containers for you.

Tests will run twice, once against an Apache Pulsar instance.

## Source

A Pulsar source connector is represented by a single consumer. By virtue of that, a source's
logical position is the respective consumer's offset in Pulsar. Internally, though, we're not saving the offset as the
position: instead, we're saving the consumer group ID, since that's all which is needed for Pulsar to find the offsets
for our consumer.

A source is getting associated with a consumer group ID the first time the `Read()` method is called.

### Configuration

| name                 | description                                                                                                                                                                                                  | required | default value             |
|----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|---------------------------|
| `url`             | URLs is the address of Apache Pulsar server.                                                                                                                                                              | true     |                           |
| `topic`              | Topic is the Pulsar topic from which records will be read.                                                                                                                                                   | true     |                           |
| `subscriptionName`   | A subscription name.                                                                                                                                                                                           | false    | `conduit-connector-pulsar` |
| `groupID`            | Defines the consumer group ID.                                                                                                                                                                               | false    |                           |
| `clientCert`         | A certificate for the Pulsar client, in PEM format. If provided, the private key needs to be provided too.                                                                                                    | false    |                           |
| `clientKey`          | A private key for the Pulsar client, in PEM format. If provided, the certificate needs to be provided too.                                                                                                    | false    |                           |
| `caCert`             | The Pulsar broker's certificate, in PEM format.                                                                                                                                                               | false    |                           |
| `insecureSkipVerify` | Controls whether a client verifies the server's certificate chain and host name. If `true`, accepts any certificate presented by the server and any host name in that certificate.                           | false    | `false`                   |

## Destination

The destination connector sends records to Pulsar.

### Configuration

There's no global, connector configuration. Each connector instance is configured separately.

| name                 | description                                                                                                                                                                                                                                                          | required | default value             |
|----------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|---------------------------|
| `url`            | URLs is the address of Apache Pulsar server.                                                                                                                                                                            | true     |                           |
| `topic`              | Topic is the Pulsar topic into which records will be written.                                                                                                                                                                                                         | true     |                           |
| `clientID`           | A Pulsar client ID.                                                                                                                                                                                                                                                   | false    | `conduit-connector-pulsar` |
| `acks`               | Acks defines the number of acknowledges from partition replicas required before receiving a response to a produce request. `none` = fire and forget, `one` = wait for the leader to acknowledge the writes, `all` = wait for the full ISR to acknowledge the writes. | false    | `all`                     |
| `deliveryTimeout`    | Message delivery timeout.                                                                                                                                                                                                                                            | false    |                           |
| `clientCert`         | A certificate for the Pulsar client, in PEM format. If provided, the private key needs to be provided too.                                                                                                                                                            | false    |                           |
| `clientKey`          | A private key for the Pulsar client, in PEM format. If provided, the certificate needs to be provided too.                                                                                                                                                            | false    |                           |
| `caCert`             | The Pulsar broker's certificate, in PEM format.                                                                                                                                                                                                                       | false    |                           |
| `insecureSkipVerify` | Controls whether a client verifies the server's certificate chain and host name. If `true`, accepts any certificate presented by the server and any host name in that certificate.                                                                                   | false    | `false`                   |

### Output format

The output format can be adjusted using configuration options provided by the connector SDK:

- `sdk.record.format`: used to choose the format
- `sdk.record.format.options`: used to configure the specifics of the chosen format

See [this article](https://conduit.io/docs/connectors/output-formats) for more info
on configuring the output format.

### Batching

Batching can also be configured using connector SDK provided options:

- `sdk.batch.size`: maximum number of records in batch before it gets written to the destination (defaults to 0, no batching)
- `sdk.batch.delay`: maximum delay before an incomplete batch is written to the destination (defaults to 0, no limit)
