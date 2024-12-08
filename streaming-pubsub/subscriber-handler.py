import signal
import time

from dapr.clients import DaprClient
from threading import Lock


start_time = 0
counter = 0
counter_lock = Lock()
total_messages = 0
TOPIC_NAME = 'my_topic'
PUBSUB_NAME = 'my_pubsub'
done = False


def process_message(message):
    global counter
    global total_messages
    global start_time

    data = message.data()

    if counter == 0:
        start_time = time.time()
        total_messages = int(message.data()["cnt"])
    counter += 1

    print(f'Subscriber received: id={data["id"]}, message="{data["message"]}", cnt="{data["cnt"]}"',
          flush=True)

    if counter == total_messages:
        print(f'Time taken to receive {counter} messages: {time.time() - start_time}', flush=True)
        return "done"

    return "continue"

def signal_handler(sig, frame):
    global terminate
    print("\nCtrl+C received! Shutting down gracefully...", flush=True)
    terminate = True

def main():
    global terminate
    # Register the signal handler for Ctrl+C
    signal.signal(signal.SIGINT, signal_handler)

    with DaprClient() as client:
        close_fn = client.subscribe_with_handler(
            pubsub_name=PUBSUB_NAME,
            topic=TOPIC_NAME,
            handler_fn=process_message
        )

        # Main loop to wait for Ctrl+C
        while not terminate:
            time.sleep(1)

        print('Closing subscription...')
        close_fn()


if __name__ == '__main__':
    main()
