# Telemetry Sink

Go app designed as three-stage pipeline using workers and channel.

---

## Run prerequisites

### env file

1. Create copy of [.env.example](.env.example) file.
2. Set values depends on your environment.

### Run docker container

#### Build image
`
cd ../../
docker build --no-cache -f telemetry-sink/docker/Dockerfile -t telemetry-sink:latest .
`

#### Run container

`
docker run --env-file ./telemetry-sink/.env telemetry-sink
`

#### Create and push tag

Build docker image using docker build from previous section.
Push image to docker registry (DockerHub by default).

[//]: # (here some steps omitted related to docker push tag to docker registry)

### Deploy to K8s cluster

Here is simple solution. It's recommended to separate helm values for different environments.
Commands arguments you will pass can be different depends on your path. 

1. Pull latest docker image for this demo:

`
   docker pull 93catdog/telemetry-sink:demo
`

2. Create env config map. Ensure that name is the same as in telemetry-sink-deployment.yaml and .env file exists with given path.

`
kubectl create configmap telemetry-sink-env \
--from-env-file=.env
`

3. Apply deployments with values using helm:

`
cd k8s/telemetry-sink
helm install telemetry-sink ./charts -f ../values/values.yaml
`

or

`
cd k8s/telemetry-sink
helm upgrade telemetry-sink ./charts -f ../values/values.yaml
`

4. Check service is running:

`
kubectl get services -l app=telemetry-sink
`

5. Check pod is running:

`
kubectl get pods -l app=telemetry-sink
`

6. Copy app container directory (with .log.json file) to your local path:

`
kubectl cp telemetry-sink-deployment-hash:/app <your_path>
`

---

## Run tests

There are no any tests for this project.

---

## Concept

App designed as three-stage pipeline using workers and channel.
Buffered channels are used to prevent immediate block on channel send operation.

1. Server send all incoming messages to buffered channel.
2. BufferedProcessor forwards messages (reads from input channel and puts them into the out channel) only on specific events (read [BufferedProcessor](#BufferedProcessor) section).
3. JsonWriter write messages to file on each receive from input channel.

### Components

#### Server

Server represent component that implements gRPC SensorServiceServer. It has single RPC handler SendSensorValues.
SendSensorValues handler sends all incoming messages to buffered channel.

#### gRPC server configuration

gRPC server configured to process requests with given allowed bandwidth (rate limit in bytes/sec) using ByteRateLimiterInterceptor.
If data flow rate exceeds the allowed bandwidth it will drop incoming messages with following status codes: ResourceExhausted.

#### BufferedProcessor

BufferedProcessor is designed to save messages to in-memory buffer from input channel and to flush messages to out channel.
It Flushes messages on such events:
- parent context is done (via <-ctx.Done());
- input channel is closed;
- timer tick with given interval;
- buffer is full (bufferSize is reached).

#### JsonWriter

JsonWriter represent component that write telemetry messages to a JSON file on each receive from input channel using Run().
Messages are written as a separate log line to file using pkg/logger (zap.Logger).
Notice: actual logs format is different from JSON.
Current implementation doesn't drain input channel as part of graceful shutdown.

---

## TODO

Future improvements:

JsonWriter
1. Add log file rotation strategy.
2. Run should drain input channel and write to file when parent context is done.
3. JsonWriter should supply file writer (logger) as dependency.

Add tests.