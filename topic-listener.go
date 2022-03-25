// Copyright 2018 Open Networking Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"gerrit.opencord.org/kafka-topic-exporter/common/logger"
	"github.com/Shopify/sarama"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	HandleFunc func(topic *string, data []byte)
}


func topicListener(ctx context.Context, topics []string, consGrp sarama.ConsumerGroup, wg sync.WaitGroup) {
	logger.Info("Starting topicListener for [%s]", topics)
	defer wg.Done()

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := Consumer{
		HandleFunc: export,
	}

	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := consGrp.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case <-ctx.Done():
			logger.Warn("terminating: context cancelled")
			return
		case <-signals:
			logger.Warn("Interrupt is detected")
			return
		}
	}

}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29

	if consumer.HandleFunc == nil {
		logger.Error("No handler for consumer ")
		return fmt.Errorf("no handler for consumer")
	}

	for message := range claim.Messages() {
		topic := string(message.Topic)
		consumer.HandleFunc(&topic, message.Value)
		session.MarkMessage(message, "")
	}

	return nil
}
