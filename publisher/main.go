/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

type benchmarkConfig struct {
	pubsubName string
	topicName  string
}

func (b *benchmarkConfig) Init() error {
	b.pubsubName = os.Getenv("PUBSUB_NAME")
	if b.pubsubName == "" {
		return fmt.Errorf("PUBSUB_NAME is not set")
	}
	b.topicName = os.Getenv("TOPIC_NAME")
	if b.topicName == "" {
		return fmt.Errorf("TOPIC_NAME is not set")
	}

	return nil
}

func main() {
	var cfg benchmarkConfig
	if err := cfg.Init(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Initialize Dapr client
	client, err := dapr.NewClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to create Dapr client: %v", err))
	}
	defer client.Close()

	publishData(ctx, client, cfg)
}

func publishData(ctx context.Context, client dapr.Client, cfg benchmarkConfig) {
	for {
		message := map[string]string{
			"id":      fmt.Sprintf("id-%d", time.Now().UnixNano()),
			"message": fmt.Sprintf("data-%d", time.Now().UnixNano()),
		}

		err := client.PublishEvent(context.Background(), cfg.pubsubName, cfg.topicName, message, dapr.PublishEventWithContentType("application/json"))
		if err != nil {
			panic(err)
		}

		fmt.Println("Published event:", message["id"])

		time.Sleep(2 * time.Second)
	}
}
