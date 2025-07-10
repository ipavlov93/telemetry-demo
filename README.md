# Telemetry-demo

Telemetry-demo is mono repository.

Telemetry-demo consists of Go modules: telemetry-sink, telemetry-node, pkg.
Go workspace and modules has been used to simplify development.

---

## Repository Go modules

### Telemetry node (telemetry-node)

[telemetry-node documentation](./telemetry-node/README.md)

### Telemetry sink (telemetry-sink)

[telemetry-sink documentation](./telemetry-sink/README.md)

### pkg

pkg contains packages with types and with utility functions that is used by other modules.

---

## Run Prerequisites

There are several options how you can run this mono repository's apps using:

1. Go (1.24.4 or upper)
2. Docker
3. K8s cluster and basic tools (kubectl, etc.)

---

## Development

### Prerequisites

1. Go (1.24.4 or upper)
2. protoc (proto compiler)

### Protoc commands

Create or update Go gRPC generated code:

`protoc \
--proto_path=proto \
--go_out=paths=source_relative:pkg/grpc/generated \
--go-grpc_out=paths=source_relative:pkg/grpc/generated \
proto/v1/sensor/sensor.proto \
proto/v1/sensor_service/service.proto`

### Formatting

`go fmt ./...`

### goimports

To group and sort import sections example:

`goimports --local telemetry-demo -l -w .`
