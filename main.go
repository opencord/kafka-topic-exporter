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
	"log"
	"strings"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

var topics []interface{}
var broker string

func kafkaInit(broker string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	var wg sync.WaitGroup

	wg.Add(2) // we are spinning up two thread and we need to wait for them to exit before stopping the kafka connection

	master, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		Panic("kafkaInit panic")
		panic(err)
	}
	defer func() {
		Debug("kafkaInit close connection")
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()
	for _, topic := range topics {
		t := topic.(string)
		go topicListener(&t, master, wg)
	}

	wg.Wait()
}

func runServer() {
	Debug("Starting Server")
	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
}

func loadConfigFile() map[string] interface{} {
	// this file path is configmap mounted in pod yaml
	jsonFile, err := os.Open("/etc/config/config.json")
	if err != nil {
		log.Fatal("Config file open failure - ",err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	confData := make(map[string]interface{})
	err = json.Unmarshal([]byte(byteValue), &confData)
	if err != nil {
		log.Fatal("Config file parsing failed - ", err)
	}
	return confData
}

func init() {

	// read configuration from configmap
	confData := loadConfigFile()

	// read kafka topics from configmap
	kafkaBroker := confData["broker"].(map[string]interface{})
	broker = kafkaBroker["host"].(string)
	topics = kafkaBroker["topics"].([]interface{})

	// read logger configuration from configmap
	logConf := confData["logger"].(map[string]interface{})
	logLevel := logConf["log_level"].(string)
	logBroker := logConf["host"].(string)

	// logger setup
	Setup(logBroker, strings.ToUpper(logLevel))

	Info("Connecting to broker: [%s]", broker)

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
	go kafkaInit(broker)
	runServer()
}
