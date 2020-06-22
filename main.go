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
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"gerrit.opencord.org/kafka-topic-exporter/common/logger"
	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
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
	prometheus.MustRegister(volthaOltTxBytesTotal)
	prometheus.MustRegister(volthaOltRxBytesTotal)
	prometheus.MustRegister(volthaOltTxPacketsTotal)
	prometheus.MustRegister(volthaOltRxPacketsTotal)
	prometheus.MustRegister(volthaOltTxErrorPacketsTotal)
	prometheus.MustRegister(volthaOltRxErrorPacketsTotal)
	prometheus.MustRegister(volthaOltTxBroadcastPacketsTotal)
	prometheus.MustRegister(volthaOltTxUnicastPacketsTotal)
	prometheus.MustRegister(volthaOltTxMulticastPacketsTotal)
	prometheus.MustRegister(volthaOltRxBroadcastPacketsTotal)
	prometheus.MustRegister(volthaOltRxUnicastPacketsTotal)
	prometheus.MustRegister(volthaOltRxMulticastPacketsTotal)

	prometheus.MustRegister(VolthaOnuLaserBiasCurrent)
	prometheus.MustRegister(volthaOnuTemperature)
	prometheus.MustRegister(VolthaOnuPowerFeedVoltage)
	prometheus.MustRegister(VolthaOnuMeanOpticalLaunchPower)
	prometheus.MustRegister(VolthaOnuReceivedOpticalPower)
	prometheus.MustRegister(VolthaOnuTransmtOpticalPower)

	prometheus.MustRegister(volthaOnuFecCorrectedCodewordsTotal)
	prometheus.MustRegister(volthaOnuFecCodewordsTotal)
	prometheus.MustRegister(volthaOnuFecCorrectedBytesTotal)
	prometheus.MustRegister(volthaOnuFecSecondsTotal)
	prometheus.MustRegister(volthaOnuFecUncorrectablewordsTotal)

	prometheus.MustRegister(volthaEthernetUniSingleCollisionTotal)
	prometheus.MustRegister(volthaEthernetUniMacLayerTramsmitErrorTotal)
	prometheus.MustRegister(volthaEthernetUniMultiCollisionTotal)
	prometheus.MustRegister(volthaEthernetUniFramestooLongTotal)
	prometheus.MustRegister(volthaEthernetUniAlignmentErrorTotal)
	prometheus.MustRegister(volthaEthernetUniCarrierErrorTotal)
	prometheus.MustRegister(volthaEthernetUniExcessiveCollisionErrorTotal)
	prometheus.MustRegister(volthaEthernetUniDeferredTxTotal)
	prometheus.MustRegister(volthaEthernetUniLateCollisionTotal)
	prometheus.MustRegister(volthaEthernetUniBufferOverflowsRxErrorTotal)
	prometheus.MustRegister(volthaEthernetUniFcsErrorTotal)
	prometheus.MustRegister(volthaEthernetUniSqeErrorTotal)
	prometheus.MustRegister(volthaEthernetUniBufferOverflowsTxErrorTotal)

	prometheus.MustRegister(volthaOnuBridgePortRxBytesTotal)
	prometheus.MustRegister(volthaOnuBridgePortRxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_64octetRxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_65_127_octetRxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_128_255_octetRxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_256_511_octetRxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_512_1023_octetRxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_1024_1518_octetRxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortRxMulticastPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortRxBroadcastPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortRxOversizePacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortRxCrcErrorPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortRxUndersizePacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortRxDropEventsTotal)

	prometheus.MustRegister(volthaOnuBridgePortTxBytesTotal)
	prometheus.MustRegister(volthaOnuBridgePortTxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_64octetTxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_65_127_octetTxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_128_255_octetTxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_256_511_octetTxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_512_1023_octetTxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePort_1024_1518_octetTxPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortTxMulticastPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortTxBroadcastPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortTxOversizePacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortTxCrcErrorPacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortTxUndersizePacketsTotal)
	prometheus.MustRegister(volthaOnuBridgePortTxDropEventsTotal)

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

	prometheus.MustRegister(onosBngUpTxBytes)
	prometheus.MustRegister(onosBngUpTxPackets)
	prometheus.MustRegister(onosBngUpDropBytes)
	prometheus.MustRegister(onosBngUpDropPackets)
	prometheus.MustRegister(onosBngControlPackets)
	prometheus.MustRegister(onosBngDownRxBytes)
	prometheus.MustRegister(onosBngDownRxPackets)
	prometheus.MustRegister(onosBngDownTxBytes)
	prometheus.MustRegister(onosBngDownTxPackets)
	prometheus.MustRegister(onosBngDownDropPackets)
	prometheus.MustRegister(onosBngDownDropBytes)

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

	//device metrics
	//TODO: Check if component level temperatures are supported by Devices,If not remove in later versions of exporter
	prometheus.MustRegister(oltDeviceCpuTemp)
	prometheus.MustRegister(oltDeviceCpuUsagePercent)
	prometheus.MustRegister(oltDeviceFanSpeed)
	prometheus.MustRegister(oltDeviceDiskTemp)
	prometheus.MustRegister(oltDeviceDiskUsagePercent)
	prometheus.MustRegister(oltDeviceRamTemp)
	prometheus.MustRegister(oltDeviceRamUsagePercent)
	prometheus.MustRegister(oltDevicePowerUsagePercent)
	prometheus.MustRegister(oltDeviceInnerSurroundTemp)
	prometheus.MustRegister(oltDevicePowerUsage)
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
