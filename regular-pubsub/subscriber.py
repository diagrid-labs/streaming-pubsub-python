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
import os

from cloudevents.sdk.event import v1
from dapr.clients.grpc._response import TopicEventResponse
from dapr.ext.grpc import App

import json


TOPIC_NAME = os.getenv('TOPIC_NAME', 'mytopic')
PUBSUB_NAME = os.getenv('PUBSUB_NAME', 'mypubsub')

app = App()

@app.subscribe(pubsub_name=PUBSUB_NAME, topic=TOPIC_NAME)
def mytopic(event: v1.Event) -> TopicEventResponse:
    data = json.loads(event.Data())

    print(
        f'Subscriber received: id={data["id"]}, message="{data["message"]}"', flush=True)

    return TopicEventResponse('success')

app.register_health_check(lambda: print('Healthy'))
print('Starting server...')
app.run(50051)
