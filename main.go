package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	txBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tx_bytes_total",
			Help: "Number of total bytes transmitted, partitioned by device_id, port_type and port_id",
		},
		[]string{"device_id", "port_type", "port_id"},
	)
	rxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rx_bytes_total",
			Help: "Number of total bytes received, partitioned by device_id, port_type and port_id",
		},
		[]string{"device_id", "port_type", "port_id"},
	)
	txPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tx_packets_total",
			Help: "Number of total packets transmitted, partitioned by device_id, port_type and port_id",
		},
		[]string{"device_id", "port_type", "port_id"},
	)
	rxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rx_packets_total",
			Help: "Number of total packets received, partitioned by device_id, port_type and port_id",
		},
		[]string{"device_id", "port_type", "port_id"},
	)

	// config
	broker  = flag.String("broker", "voltha-kafka.default.svc.cluster.local:9092", "The Kafka broker")
	topic   = flag.String("topic", "voltha.kpis", "The Kafka topic")
)

var brokers []string
var messageCountStart int

func prefixToLabels(prefix string) (string, string, string) {
	var p = strings.Split(prefix, ".")
	var deviceId, portType, portId string = "", "", ""
	if len(p) == 5 {
		// format is voltha.openolt.000130af0b0b2c51.pon.0
		deviceId = p[2]
		portType = p[3]
		portId = p[4]
	}
	if len(p) == 4 {
		// fomrat is voltha.openolt.000130af0b0b2c51nni.129
		s := p[2]
		deviceId = string(s[0 : len(s)-3])
		portType = string(s[len(s)-3:])
		portId = p[3]
	}

	return deviceId, portType, portId

}

func interfaceToFloat(unk interface{}) float64 {
	switch i := unk.(type) {
	case float64:
		return i
	case float32:
		return float64(i)
	case int64:
		return float64(i)
	default:
		return math.NaN()
	}
}

func kafkaInit(brokers []string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetOldest)
	if err != nil {
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
				messageCountStart++

				var label map[string]interface{}
				json.Unmarshal(msg.Value, &label)

				var data map[string]map[string]map[string]interface{}
				json.Unmarshal(msg.Value, &data)

				var tagString = reflect.ValueOf(label["prefixes"]).MapKeys()[0].String()

				fmt.Println("tagString: "+tagString, "\n")
				fmt.Println("data: ", data["prefixes"][tagString]["metrics"], "\n")

				v, ok := data["prefixes"][tagString]["metrics"].(map[string]interface{})
				if !ok {
					// Can't assert, handle error.
					fmt.Println("Eroror")
				}
				for k, s := range v {
					fmt.Println("Type k: ", reflect.TypeOf(k))
					fmt.Println("Type: ", reflect.TypeOf(s))
					fmt.Printf("Value: %v\n", s)

					d, pt, pi := prefixToLabels(tagString)

					if k == "tx_bytes" {
						txBytesTotal.WithLabelValues(d, pt, pi).Set(interfaceToFloat(s))
					}
					if k == "rx_bytes" {
						rxBytesTotal.WithLabelValues(d, pt, pi).Set(interfaceToFloat(s))
					}
					if k == "tx_packets" {
						txPacketsTotal.WithLabelValues(d, pt, pi).Set(interfaceToFloat(s))
					}
					if k == "rx_packets" {
						rxPacketsTotal.WithLabelValues(d, pt, pi).Set(interfaceToFloat(s))
					}
				}

				// fmt.Println("data: ", data["prefixes"][tagString]["metrics"].tx_bytes, "\n")
				// var txBytesTotalValue = data["prefixes"][tagString]["metrics"]["tx_bytes"]

				// d, pt, pi := prefixToLabels(tagString)

				// txBytesTotal.WithLabelValues(d, pt, pi).Set(float64(txBytesTotalValue))

				// fmt.Println("Adding txBytesTotal metric: ", d, pt, pi, txBytesTotalValue)

			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	fmt.Println("Processed", messageCountStart, "messages")
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
	fmt.Println("Listening to topic: ", *topic)

	// register metrics within Prometheus
	prometheus.MustRegister(txBytesTotal)
	prometheus.MustRegister(rxBytesTotal)
	prometheus.MustRegister(txPacketsTotal)
	prometheus.MustRegister(rxPacketsTotal)
}

func main() {
	go kafkaInit(brokers)
	runServer()
}
