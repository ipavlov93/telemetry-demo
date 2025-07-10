# Telemetry node

Go app designed as two-stage pipeline using workers and channel.

---

## Run prerequisites

### env file

1. Create copy of [.env.example](.env.example) file.
2. Set values depends on your environment.

### Run docker container

#### Build image
`
cd ../../
docker build --no-cache -f telemetry-node/docker/Dockerfile -t telemetry-node:latest .
`

#### Run container

`
docker run --env-file ./telemetry-node/.env telemetry-node
`

#### Create and push tag

Build docker image using docker build from previous section.
Push image to docker registry (DockerHub by default).

[//]: # (here some steps omitted related to docker push tag to docker registry)

### Deploy to K8s cluster

1. Pull latest docker image for this demo:

`
   docker pull 93catdog/telemetry-node:latest
`

2. Create env config map. Ensure that name is the same as in telemetry-node-deployment.yaml and .env file exists with given path.

`
kubectl create configmap telemetry-node-env \
--from-env-file=.env
`

3. Apply deployment:

`
kubectl apply -f k8s/telemetry-node-deployment.yaml
`

---

## Run tests

There are unit few unit tests for this project.
It's recommended to run the tests using this command:
`go test ./... -race -count 1`

---

## Concept

App designed as two-stage pipeline using workers and channel.

1. Sensor data generated periodically by IntervalSensor.
    - IntervalSensor is designed as one per Run() worker.
    - Data are buffered before sent, data send occurs when buffer is full.
    - Run() worker flushes buffer on context cancellation.
2. Sensor data sent with constant rate by SensorService using gRPC protocol.
    - Data buffered before sent, data send occurs when buffer is full.
    - SensorService doesn't drain remaining messages from channel if context is cancelled.

### Components

#### Interval Sensor

Interval Sensor produces data by separate go routine (worker) with constant rate using Run().
Produced data are buffered and send to channel when buffer is full.
Buffered channel is used to prevent immediate block on channel send operation.

#### SensorService

SensorService is designed consumes data from channel and sends data to destination server using corresponding client.
SensorService doesn't retry reestablish connection on transport-level failures (when server is unreachable).

#### Sender worker

Rate Limiter (token-based setup) is used to achieve send messages with given RPS.

#### gRPC client configuration

gRPC client configured to send requests with given RPS and constant interval retry strategy. 
Retry mechanism partially resolve issue with network instability.
The default gRPC retry mechanism does not re-establish broken connections (e.g., when the server is unreachable).
RPC calls retry on failure with respect to the following status codes: Unavailable, ResourceExhausted, DeadlineExceeded.
The caller is responsible for ensuring that context lifetimes are long enough to support full retry cycles.

---

## TODO

Future improvements:
1. Enable TLS for network communication.
2. Add custom implementation of retry strategy with transport connection reestablishment and exponential backoff.
3. Add hardcoded variables in main() to config.
4. Add gRPC client configuration options to config.
5. Add IntervalSensor RPS configuration to config.
6. Add more tests.
