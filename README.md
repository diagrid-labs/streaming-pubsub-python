# Streaming PubSub in the Dapr Python SDK

These examples demonstrate how to use the Dapr Python SDK to create a subscriber application that listens to a topic and receives messages from a publisher application. The subscriber application can be run in two modes: regular subscription and streaming subscription.

## Running in self-hosted mode
### Prerequisites for running locally
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
This example uses Redis as the pub/sub broker. Ensure Redis is installed and running locally.
Also make sure the details of the default Redis component in the `../components` directory match your local Redis installation. Refer to the Dapr pub/sub documentation for details.

### Run the example: self-hosted mode

1. Run the publisher  
The publisher is written in Go, to make use of goroutines
```bash
export PUBSUB_NAME=mypubsub
export TOPIC_NAME=mytopic
export NUM_MESSAGES_PER_CALL=10
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
dapr run --app-id python-subscriber --resources-path=../components -- python3 subscriber-handler.py
```

## Running in Kubernetes


### Prerequisites for running in Kubernetes
- A Kubernetes cluster

### Run the example

Run the setup script:
```bash
./setup.sh
```

The script will set the context to `kind-kind`, but you can override this by running it with the desired context as the first argument:
```
./setup.sh mycontext
```