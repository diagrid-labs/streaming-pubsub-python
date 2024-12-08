import time

from dapr.clients import DaprClient
from dapr.clients.grpc.subscription import StreamInactiveError
from dapr.common.pubsub.subscription import StreamCancelledError

start_time = 0
counter = 0
total_messages = 0
TOPIC_NAME = 'my_topic'
PUBSUB_NAME = 'my_pubsub'

def process_message(message):
    global counter
    global total_messages
    global start_time

    data = message.data()

    if counter == 0:
        start_time = time.time()
        total_messages = int(data["cnt"])
    counter += 1


    print(f'Subscriber received: id={data["id"]}, message="{data["message"]}", cnt="{data["cnt"]}"',
          flush=True)

    if counter == total_messages:
        print(f'Time taken to receive {counter} messages: {time.time() - start_time}', flush=True)
        return "done"

    return "continue"

def main():
    with DaprClient() as client:
        global counter

        try:
            subscription = client.subscribe(
                pubsub_name=PUBSUB_NAME, topic=TOPIC_NAME
            )
        except Exception as e:
            print(f'Error occurred: {e}')
            return

        try:
            for message in subscription:
                if message is None:
                    print('No message received. The stream might have been cancelled.')
                    continue

                try:
                    response_status = process_message(message)

                    subscription.respond_success(message)
                    if response_status == 'done':
                        break
                except StreamInactiveError:
                    print('Stream is inactive. Retrying...')
                    time.sleep(1)
                    continue
                except StreamCancelledError:
                    print('Stream was cancelled')
                    break
                except Exception as e:
                    print(f'Error occurred during message processing: {e}')

        finally:
            print('Closing subscription...')
            subscription.close()


if __name__ == '__main__':
    main()




# import time
#
# from dapr.clients import DaprClient
# from dapr.clients.grpc._response import TopicEventResponse
#
# start_time = 0
# counter = 0
# total_messages = 0
# topic_name = 'my_topic'
# pubsub_name = 'my_pubsub'
#
#
#
# def process_message(message):
#     global counter
#     global start_time
#     global total_messages
#
#     data = message.data()
#
#     if counter == 0:
#         start_time = time.time()
#         total_messages = int(message.data()["cnt"])
#     counter += 1
#
#
#     print(f'Subscriber received: id={data["id"]}, message="{data["message"]}", cnt="{data["cnt"]}"',
#           flush=True)
#
#     if counter == int(data["cnt"]):
#         print(f'Time taken to receive {counter} messages: {time.time() - start_time}', flush=True)
#
#     return TopicEventResponse('success')
#
#
# def main():
#     with DaprClient() as client:
#         close_fn = client.subscribe_with_handler(
#             pubsub_name=pubsub_name,
#             topic=topic_name,
#             handler_fn=process_message,
#         )
#
#         # Keep the program running until all messages have been processed
#         while counter < total_messages:
#             time.sleep(1)
#
#         print('Closing subscription...')
#         close_fn()
#
#
# if __name__ == '__main__':
#     main()
