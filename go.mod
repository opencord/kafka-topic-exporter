module gerrit.opencord.org/kafka-topic-exporter

go 1.16

require (
	github.com/Shopify/sarama v1.32.0
	github.com/gfremex/logrus-kafka-hook v0.0.0-20180109031623-f62e125fcbfe
	github.com/golang/protobuf v1.5.2
	github.com/opencord/device-management-interface v1.4.0
	github.com/opencord/voltha-protos/v5 v5.2.4
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.6.0 // indirect
	github.com/prometheus/procfs v0.0.3 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.2.3
)
