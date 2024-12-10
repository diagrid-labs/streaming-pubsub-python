import os
import signal
import time

from dapr.clients import DaprClient
from dapr.clients.grpc._response import TopicEventResponse


TOPIC_NAME = os.getenv('TOPIC_NAME', 'mytopic')
PUBSUB_NAME = os.getenv('PUBSUB_NAME', 'mypubsub')


def process_message(message):
    data = message.data()

    print(f'Subscriber received: id={data["id"]}, message="{data["message"]}"', flush=True)

    return TopicEventResponse('success')

def main():
    with DaprClient() as client:
        close_fn = client.subscribe_with_handler(
            pubsub_name=PUBSUB_NAME,
            topic=TOPIC_NAME,
            handler_fn=process_message
        )

        while True:
            time.sleep(1)

        print('Closing subscription...')
        close_fn()


if __name__ == '__main__':
    main()
