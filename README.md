# Telemetry-demo

Telemetry-demo is mono repository.

Telemetry-demo consists of Go modules: telemetry-sink, telemetry-node, pkg.

Go workspace and modules has been used to simplify development.

## Demo branch

The actual branch for demo (with latest changes) `branch-for-demo`: https://github.com/ipavlov93/telemetry-demo/tree/branch-for-demo

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
3. K8s cluster and basic tools (helm, kubectl, etc.)

You can find how to configure and run apps in corresponding documentation.
Notice: .env file variable GRPC_SERVER_SOCKET would be reset during deploy to K8s cluster. 

---

## Other directories

### proto

/proto directory contains proto files.

### k8s

/k8s directory contains helm charts and templates.

## FAQ

1. Why interfaces haven't been defined for each app's component with respect dependency inversion principle ?
Answer: I've share the following idea: don't create redundant contracts (intefaces) without usage, only for testing purpose. (I am going to add them with tests in future).
Components that accepts dependencis should respect dependency inversion principle.

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
