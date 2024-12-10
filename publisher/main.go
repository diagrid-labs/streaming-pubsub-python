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
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

type benchmarkConfig struct {
	pubsubName         string
	topicName          string
	numMessagesPerCall int
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

	var err error
	b.numMessagesPerCall, err = strconv.Atoi(os.Getenv("NUM_MESSAGES_PER_CALL"))
	if err != nil {
		return fmt.Errorf("Failed to parse NUM_MESSAGES_PER_CALL: %v", err)
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
		var publishEventsData []interface{}
		for j := 0; j < cfg.numMessagesPerCall; j++ {
			message := map[string]string{
				"id":      fmt.Sprintf("id-%d", j),
				"message": fmt.Sprintf("data-%d", j),
			}

			data, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Error marshalling message: %v\n", err)
				continue
			}

			publishEventsData = append(publishEventsData, data)
		}

		// Publish multiple events
		if res := client.PublishEvents(ctx, cfg.pubsubName, cfg.topicName, publishEventsData, dapr.PublishEventsWithContentType("application/json")); res.Error != nil {
			panic(res.Error)
		}
		fmt.Println("Published", len(publishEventsData), "events")

		time.Sleep(2 * time.Second)
	}
}
