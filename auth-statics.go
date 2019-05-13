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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
)

func STATSListner(topic *string, master sarama.Consumer, wg sync.WaitGroup) {
	fmt.Println("Starting STATSListner")
	defer wg.Done()
	// OffsetOldest stands for the oldest offset available on the broker for a
    // partition. You can send this to a client's GetOffset method to get this
    // offset, or when calling ConsumePartition to start consuming from the
    // oldest offset that is still available on the broker.
	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Println("STATSListner panic")
		panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				fmt.Println("Reading message from consumer")
				kpi := AuthStatKPI{}

				err := json.Unmarshal(msg.Value, &kpi)
				fmt.Println("Compare message successfull")
				if err != nil {
					log.Fatal(err)
				}

				exportAuthStatKPI(kpi)

			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
}
