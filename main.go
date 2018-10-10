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
	"flag"
	"fmt"
	"net/http"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	broker      = flag.String("broker", "voltha-kafka.default.svc.cluster.local:9092", "The Kafka broker")
	volthaTopic = "voltha.kpis"
	onosTopic   = "onos.kpis"

	volthaTopicPointer = &volthaTopic
	onosTopicPointer   = &onosTopic
)

var brokers []string

func kafkaInit(brokers []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	var wg sync.WaitGroup

	wg.Add(2) // we are spinning up two thread and we need to wait for them to exit before stopping the kafka connection

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println("kafkaInit panic")
		panic(err)
	}
	defer func() {
		fmt.Println("kafkaInit close connection")
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()
	go VOLTHAListener(volthaTopicPointer, master, wg)
	go ONOSListener(onosTopicPointer, master, wg)

	wg.Wait()
}

func runServer() {
	fmt.Println("Starting Server")
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
}

func init() {

	// read config from cli flags
	flag.Parse()
	brokers = make([]string, 0)
	brokers = []string{*broker}
	fmt.Println("Connecting to broker: ", brokers)
	fmt.Println("Listening to voltha on topic: ", *volthaTopicPointer)
	fmt.Println("Listening to onos on topic: ", *onosTopicPointer)

	// register metrics within Prometheus
	prometheus.MustRegister(volthaTxBytesTotal)
	prometheus.MustRegister(volthaRxBytesTotal)
	prometheus.MustRegister(volthaTxPacketsTotal)
	prometheus.MustRegister(volthaRxPacketsTotal)
	prometheus.MustRegister(volthaTxErrorPacketsTotal)
	prometheus.MustRegister(volthaRxErrorPacketsTotal)

	prometheus.MustRegister(onosTxBytesTotal)
	prometheus.MustRegister(onosRxBytesTotal)
	prometheus.MustRegister(onosTxPacketsTotal)
	prometheus.MustRegister(onosRxPacketsTotal)
	prometheus.MustRegister(onosTxDropPacketsTotal)
	prometheus.MustRegister(onosRxDropPacketsTotal)
}

func main() {
	go kafkaInit(brokers)
	runServer()
}
