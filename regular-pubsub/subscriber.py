# ------------------------------------------------------------
# Copyright 2024 The Dapr Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------
import time
from threading import Lock

from cloudevents.sdk.event import v1
from dapr.ext.grpc import App
from dapr.clients.grpc._response import TopicEventResponse

import json

app = App()

start_time = 0
counter = 0
counter_lock = Lock()
TOPIC_NAME = 'my_topic'
PUBSUB_NAME = 'my_pubsub'


@app.subscribe(pubsub_name=PUBSUB_NAME, topic=TOPIC_NAME)
def mytopic(event: v1.Event) -> TopicEventResponse:
    global start_time

    data = json.loads(event.Data())

    with counter_lock:  # Acquire the lock before modifying counter
        global counter
        if counter == 0:
            start_time = time.time()
        counter += 1

        print(
            f'Subscriber received: id={data["id"]}, message="{data["message"]}", cnt="{data["cnt"]}"',
            flush=True)

        if counter == int(data["cnt"]):
            print(f'Time taken to receive {counter} messages: {time.time() - start_time}', flush=True)

    return TopicEventResponse('success')

app.register_health_check(lambda: print('Healthy'))
app.run(50051)
