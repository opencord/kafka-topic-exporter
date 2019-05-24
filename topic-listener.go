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
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
	"gerrit.opencord.org/kafka-topic-exporter/common/logger"
)

func topicListener(topic *string, master sarama.Consumer, wg sync.WaitGroup) {
	logger.Info("Starting topicListener for [%s]", *topic)
	defer wg.Done()
	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetOldest)
	if err != nil {
		logger.Error("topicListener panic")
		panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				logger.Error("%s", err)
			case msg := <-consumer.Messages():
				logger.Debug("Message on %s: %s", *topic, string(msg.Value))
				export(topic, msg.Value)
			case <-signals:
				logger.Warn("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
}
