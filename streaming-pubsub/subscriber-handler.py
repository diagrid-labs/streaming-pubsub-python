import os
import signal
import time

from dapr.clients import DaprClient
from dapr.clients.grpc._response import TopicEventResponse


TOPIC_NAME = os.getenv('TOPIC_NAME', 'mytopic')
PUBSUB_NAME = os.getenv('PUBSUB_NAME', 'mypubsub')

terminate = False

def process_message(message):
    data = message.data()

    # Your processing logic here
    print(f'Subscriber received: id={data["id"]}, message="{data["message"]}"', flush=True)

    return TopicEventResponse('success') # other options are 'retry' and 'drop'

def signal_handler(sig, frame):
    global terminate
    print("\nCtrl+C received! Shutting down gracefully...", flush=True)
    terminate = True

def main():
    global terminate

    # Register signal handlers for graceful shutdown
    signal.signal(signal.SIGINT, signal_handler)  # Handle Ctrl+C
    signal.signal(signal.SIGTERM, signal_handler)  # Handle termination signal

    with DaprClient() as client:
        close_fn = client.subscribe_with_handler(
            pubsub_name=PUBSUB_NAME,
            topic=TOPIC_NAME,
            handler_fn=process_message
        )

        try:
            while not terminate:
                time.sleep(1)
        finally:
            print('Closing subscription...', flush=True)
            close_fn()


if __name__ == '__main__':
    main()
