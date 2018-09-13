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

import "github.com/prometheus/client_golang/prometheus"

var (
	txBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tx_bytes_total",
			Help: "Number of total bytes transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	rxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rx_bytes_total",
			Help: "Number of total bytes received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	txPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tx_packets_total",
			Help: "Number of total packets transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	rxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rx_packets_total",
			Help: "Number of total packets received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	txErrorPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "tx_error_packets_total",
			Help: "Number of total transmitted packets error",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	rxErrorPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rx_error_packets_total",
			Help: "Number of total received packets error",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
)

func export(kpi KPI) {

	for _, data := range kpi.SliceDatas {
		txBytesTotal.WithLabelValues(
			data.Metadata.LogicalDeviceID,
			data.Metadata.SerialNumber,
			data.Metadata.DeviceID,
			data.Metadata.Context.InterfaceID,
			data.Metadata.Context.PonID,
			data.Metadata.Context.PortNumber,
			data.Metadata.Title,
		).Set(data.Metrics.TxBytes)

		rxBytesTotal.WithLabelValues(
			data.Metadata.LogicalDeviceID,
			data.Metadata.SerialNumber,
			data.Metadata.DeviceID,
			data.Metadata.Context.InterfaceID,
			data.Metadata.Context.PonID,
			data.Metadata.Context.PortNumber,
			data.Metadata.Title,
		).Set(data.Metrics.RxBytes)

		txPacketsTotal.WithLabelValues(
			data.Metadata.LogicalDeviceID,
			data.Metadata.SerialNumber,
			data.Metadata.DeviceID,
			data.Metadata.Context.InterfaceID,
			data.Metadata.Context.PonID,
			data.Metadata.Context.PortNumber,
			data.Metadata.Title,
		).Set(data.Metrics.TxPackets)

		rxPacketsTotal.WithLabelValues(
			data.Metadata.LogicalDeviceID,
			data.Metadata.SerialNumber,
			data.Metadata.DeviceID,
			data.Metadata.Context.InterfaceID,
			data.Metadata.Context.PonID,
			data.Metadata.Context.PortNumber,
			data.Metadata.Title,
		).Set(data.Metrics.RxPackets)

		txErrorPacketsTotal.WithLabelValues(
			data.Metadata.LogicalDeviceID,
			data.Metadata.SerialNumber,
			data.Metadata.DeviceID,
			data.Metadata.Context.InterfaceID,
			data.Metadata.Context.PonID,
			data.Metadata.Context.PortNumber,
			data.Metadata.Title,
		).Set(data.Metrics.TxErrorPackets)

		rxErrorPacketsTotal.WithLabelValues(
			data.Metadata.LogicalDeviceID,
			data.Metadata.SerialNumber,
			data.Metadata.DeviceID,
			data.Metadata.Context.InterfaceID,
			data.Metadata.Context.PonID,
			data.Metadata.Context.PortNumber,
			data.Metadata.Title,
		).Set(data.Metrics.RxErrorPackets)

		// TODO add metrics for:
		// TxBcastPackets
		// TxUnicastPackets
		// TxMulticastPackets
		// RxBcastPackets
		// RxMulticastPackets

	}
}
