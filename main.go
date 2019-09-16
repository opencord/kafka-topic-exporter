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
	"gerrit.opencord.org/kafka-topic-exporter/common/logger"
	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func kafkaInit(broker BrokerInfo) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	var wg sync.WaitGroup

	master, err := sarama.NewConsumer([]string{broker.Host}, config)

	if err != nil {
		logger.Panic("kafkaInit panic")
		panic(err)
	}
	defer func() {
		logger.Debug("kafkaInit close connection")
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	// read topics from config
	topics := broker.Topics

	// we are spinning threads for each topic, we need to wait for
	// them to exit before stopping the kafka connection
	wg.Add(len(topics))

	for _, topic := range topics {
		t := topic
		go topicListener(&t, master, wg)
	}

	wg.Wait()
}

func runServer(target TargetInfo) {
	if target.Port == 0 {
		logger.Warn("Prometheus target port not configured, using default 8080")
		target.Port = 8080
	}
	logger.Debug("Starting HTTP Server on %d port", target.Port)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+strconv.Itoa(target.Port), nil)
	if err != nil {
		logger.Error("HTTP Server Error: %s", err.Error())
       }
}

func init() {
	// register metrics within Prometheus
	prometheus.MustRegister(volthaTxBytesTotal)
	prometheus.MustRegister(volthaRxBytesTotal)
	prometheus.MustRegister(volthaTxPacketsTotal)
	prometheus.MustRegister(volthaRxPacketsTotal)
	prometheus.MustRegister(volthaTxErrorPacketsTotal)
	prometheus.MustRegister(volthaRxErrorPacketsTotal)

	prometheus.MustRegister(VolthaOnuLaserBiasCurrent)
	prometheus.MustRegister(volthaOnuTemperature)
	prometheus.MustRegister(VolthaOnuPowerFeedVoltage)
	prometheus.MustRegister(VolthaOnuMeanOpticalLaunchPower)
	prometheus.MustRegister(VolthaOnuReceivedOpticalPower)

	prometheus.MustRegister(onosTxBytesTotal)
	prometheus.MustRegister(onosRxBytesTotal)
	prometheus.MustRegister(onosTxPacketsTotal)
	prometheus.MustRegister(onosRxPacketsTotal)
	prometheus.MustRegister(onosTxDropPacketsTotal)
	prometheus.MustRegister(onosRxDropPacketsTotal)

	prometheus.MustRegister(onosaaaRxAcceptResponses)
	prometheus.MustRegister(onosaaaRxRejectResponses)
	prometheus.MustRegister(onosaaaRxChallengeResponses)
	prometheus.MustRegister(onosaaaTxAccessRequests)
	prometheus.MustRegister(onosaaaRxInvalidValidators)
	prometheus.MustRegister(onosaaaRxUnknownType)
	prometheus.MustRegister(onosaaaPendingRequests)
	prometheus.MustRegister(onosaaaRxDroppedResponses)
	prometheus.MustRegister(onosaaaRxMalformedResponses)
	prometheus.MustRegister(onosaaaRxUnknownserver)
	prometheus.MustRegister(onosaaaRequestRttMillis)
	prometheus.MustRegister(onosaaaRequestReTx)

	prometheus.MustRegister(onosPppoeUpTermBytes)
	prometheus.MustRegister(onosPppoeUpTermPackets)
	prometheus.MustRegister(onosPppoeUpDropBytes)
	prometheus.MustRegister(onosPppoeUpDropPackets)
	prometheus.MustRegister(onosPppoeUpControlPackets)
	prometheus.MustRegister(onosPppoeDownRxBytes)
	prometheus.MustRegister(onosPppoeDownRxPackets)
	prometheus.MustRegister(onosPppoeDownTxBytes)
	prometheus.MustRegister(onosPppoeDownTxPackets)

	prometheus.MustRegister(deviceLaserBiasCurrent)
	prometheus.MustRegister(deviceTemperature)
	prometheus.MustRegister(deviceTxPower)
	prometheus.MustRegister(deviceVoltage)

	prometheus.MustRegister(onosaaaRxEapolLogoff)
	prometheus.MustRegister(onosaaaTxEapolResIdentityMsg)
	prometheus.MustRegister(onosaaaTxAuthSuccess)
	prometheus.MustRegister(onosaaaTxAuthFailure)
	prometheus.MustRegister(onosaaaTxStartReq)
	prometheus.MustRegister(onosaaaEapPktTxAuthChooseEap)
	prometheus.MustRegister(onosaaaTxRespnotNak)

        prometheus.MustRegister(onosaaaEapolFramesTx)
	prometheus.MustRegister(onosaaaAuthStateIdle)
	prometheus.MustRegister(onosaaaRequestIdFramesTx)
	prometheus.MustRegister(onosaaaRequestEapFramesTx)
	prometheus.MustRegister(onosaaaInvalidPktType)
	prometheus.MustRegister(onosaaaInvalidBodyLength)
	prometheus.MustRegister(onosaaaValidEapolFramesRx)
	prometheus.MustRegister(onosaaaPendingResSupplicant)
	prometheus.MustRegister(onosaaaRxResIdEapFrames)
}

func loadConfigFile() Config {
	m := Config{}
	// this file path is configmap mounted in pod yaml
	yamlFile, err := ioutil.ReadFile("/etc/config/conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return m
}

func main() {
	// load configuration
	conf := loadConfigFile()

	// logger setup
	logger.Setup(conf.Logger.Host, strings.ToUpper(conf.Logger.LogLevel))
	logger.Info("Connecting to broker: [%s]", conf.Broker.Host)

	go kafkaInit(conf.Broker)
	runServer(conf.Target)
}
