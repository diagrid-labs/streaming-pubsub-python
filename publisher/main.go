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
	"sync"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

type benchmarkConfig struct {
	pubsubName            string
	topicName             string
	numCycles             int
	numGoRoutinesPerCycle int
	numMessagesPerCall    int
	totalMsgCnt           int
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
	b.numCycles, err = strconv.Atoi(os.Getenv("NUM_CYCLES"))
	if err != nil {
		return fmt.Errorf("Failed to parse NUM_CYCLES: %v", err)
	}
	fmt.Println("b.numCycles", b.numCycles)

	b.numGoRoutinesPerCycle, err = strconv.Atoi(os.Getenv("NUM_GOROUTINES_PER_CYCLE"))
	if err != nil {
		return fmt.Errorf("Failed to parse NUM_GOROUTINES_PER_CYCLE: %v", err)
	}

	b.numMessagesPerCall, err = strconv.Atoi(os.Getenv("NUM_MESSAGES_PER_CALL"))
	if err != nil {
		return fmt.Errorf("Failed to parse NUM_MESSAGES_PER_CALL: %v", err)
	}

	b.totalMsgCnt = b.numCycles * b.numGoRoutinesPerCycle * b.numMessagesPerCall

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

	if cfg.numCycles > 0 {
		publishBenchmarkData(ctx, client, cfg)
	} else {
		publishDataContinuously(ctx, client, cfg)
	}

	fmt.Println("data published")
}

func publishBenchmarkData(ctx context.Context, client dapr.Client, cfg benchmarkConfig) {
	for i := range cfg.numCycles {
		sendBatch(ctx, client, cfg, i)
		fmt.Printf("Batch %d completed\n", i)
	}
}

func publishDataContinuously(ctx context.Context, client dapr.Client, cfg benchmarkConfig) {
	for {
		sendBatch(ctx, client, cfg, 0)
		time.Sleep(1 * time.Second)
	}
}

func sendBatch(ctx context.Context, client dapr.Client, cfg benchmarkConfig, batchNumber int) {
	var wg sync.WaitGroup

	for i := range cfg.numGoRoutinesPerCycle {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var publishEventsData []interface{}
			for j := 0; j < cfg.numMessagesPerCall; j++ {
				message := map[string]string{
					"id":      fmt.Sprintf("id-%d-%d-%d", batchNumber, i, j),
					"message": fmt.Sprintf("data-%d-%d-%d", batchNumber, i, j),
					"cnt":     strconv.Itoa(cfg.totalMsgCnt),
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
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}
