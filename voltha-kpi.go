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

func VOLTHAListener(topic *string, master sarama.Consumer, wg sync.WaitGroup) {
	fmt.Println("Starting VOLTHAListener")
	defer wg.Done()
	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Println("VOLTHAListener panic")
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
				// fmt.Println(string(msg.Value))

				kpi := VolthaKPI{}

				err := json.Unmarshal(msg.Value, &kpi)

				if err != nil {
					log.Fatal(err)
				}

				exportVolthaKPI(kpi)

			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
}
