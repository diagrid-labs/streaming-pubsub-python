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
	"strconv"
	"sync"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	// Configuration variables
	pubsubName            = "my_pubsub"
	topicName             = "my_topic"
	numCycles             = 10
	numGoRoutinesPerCycle = 10
	numMessagesPerCall    = 10
	totalMsgCnt           = numCycles * numGoRoutinesPerCycle * numMessagesPerCall
)

func main() {
	ctx := context.Background()

	// Initialize Dapr client
	client, err := dapr.NewClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to create Dapr client: %v", err))
	}
	defer client.Close()

	// Publish messages in cycles
	for i := range numCycles {
		sendBatch(ctx, client, i)
	}

	fmt.Println("data published")
}

func sendBatch(ctx context.Context, client dapr.Client, batchNumber int) {
	var wg sync.WaitGroup

	for i := range numGoRoutinesPerCycle {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var publishEventsData []interface{}
			for j := 0; j < numMessagesPerCall; j++ {
				message := map[string]string{
					"id":      fmt.Sprintf("id-%d-%d-%d", batchNumber, i, j),
					"message": fmt.Sprintf("data-%d-%d-%d", batchNumber, i, j),
					"cnt":     strconv.Itoa(totalMsgCnt),
				}

				data, err := json.Marshal(message)
				if err != nil {
					fmt.Printf("Error marshalling message: %v\n", err)
					continue
				}

				publishEventsData = append(publishEventsData, data)
			}

			// Publish multiple events
			if res := client.PublishEvents(ctx, pubsubName, topicName, publishEventsData, dapr.PublishEventsWithContentType("application/json")); res.Error != nil {
				panic(res.Error)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	fmt.Printf("Batch %d completed\n", batchNumber)
}
