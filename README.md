# Streaming benchmark in the Dapr python-sdk

This benchmark compares the performance of regular pub/sub and streaming pub/sub features in the Dapr Python SDK. It includes a Go-based publisher to simulate high-throughput scenarios.

## Prerequisites
**1. Install Dapr CLI and Runtime**  
Ensure that you have Dapr installed on your system:
```bash
dapr --version
```
If not installed, follow the [Dapr installation guide](https://docs.dapr.io/getting-started)

2. Install Python 3.8+  
Ensure that you have Python 3.8 or higher installed on your system. If not, download and install it from the [official Python website](https://www.python.org/downloads/).

3. Install Dependencies  
Install the required Python libraries for the subscriber applications:
```bash
pip3 install -r requirements.txt
```

4. Configure Components
This benchmark uses Redis as the pub/sub broker. Ensure Redis is installed and running locally.
Also make sure the details of the default Redis component in the `../components` directory match your local Redis installation. Refer to the Dapr pub/sub documentation for details.

## Run the example

1. Run the publisher  
The publisher is written in Go, to make use of goroutines
```bash
cd publisher
dapr run --app-id python-publisher --app-protocol grpc --enable-app-health-check --resources-path=../components -- go run main.go
```

2. Run the subscriber with regular subscription
```bash
cd regular-pubsub
dapr run --app-id python-subscriber --app-protocol grpc --app-port 50051 --resources-path=../components -- python3 subscriber.py
```


3. Run the subscriber with streaming subscription
```bash
cd streaming-pubsub
dapr run --app-id python-subscriber --resources-path=../components -- python3 subscriber.py
```

## Monitoring the Benchmark
Observe the logs of both the publisher and subscriber to verify message flow and measure performance.
Compare the message throughput and latency between regular pub/sub and streaming pub/sub.
You can tweak the number of sent messages through the `numGoRoutinesPerCycle`, `numMessagesPerCall` and `totalMsgCnt` variables in the [publisher](publisher/main.go) to simulate different scenarios.