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
	"github.com/prometheus/client_golang/prometheus"
)

var (
	volthaTxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_tx_bytes_total",
			Help: "Number of total bytes transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	volthaRxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_rx_bytes_total",
			Help: "Number of total bytes received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	volthaTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_tx_packets_total",
			Help: "Number of total packets transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	volthaRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_rx_packets_total",
			Help: "Number of total packets received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaTxErrorPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_tx_error_packets_total",
			Help: "Number of total transmitted packets error",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaRxErrorPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_rx_error_packets_total",
			Help: "Number of total received packets error",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
)

func exportVolthaKPI(kpi VolthaKPI) {

	for _, data := range kpi.SliceDatas {
		switch title := data.Metadata.Title; title {
		case "Ethernet", "PON":
			volthaTxBytesTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.TxBytes)

			volthaRxBytesTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.RxBytes)

			volthaTxPacketsTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.TxPackets)

			volthaRxPacketsTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.RxPackets)

			volthaTxErrorPacketsTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.TxErrorPackets)

			volthaRxErrorPacketsTotal.WithLabelValues(
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

		case "Ethernet_Bridge_Port_History":
			if data.Metadata.Context.Upstream == "True" {
				// ONU. Extended Ethernet statistics.
				volthaTxPacketsTotal.WithLabelValues(
					data.Metadata.LogicalDeviceID,
					data.Metadata.SerialNumber,
					data.Metadata.DeviceID,
					"NA", // InterfaceID
					"NA", // PonID
					"NA", // PortNumber
					data.Metadata.Title,
				).Add(data.Metrics.Packets)

				volthaTxBytesTotal.WithLabelValues(
					data.Metadata.LogicalDeviceID,
					data.Metadata.SerialNumber,
					data.Metadata.DeviceID,
					"NA", // InterfaceID
					"NA", // PonID
					"NA", // PortNumber
					data.Metadata.Title,
				).Add(data.Metrics.Octets)
			} else {
				// ONU. Extended Ethernet statistics.
				volthaRxPacketsTotal.WithLabelValues(
					data.Metadata.LogicalDeviceID,
					data.Metadata.SerialNumber,
					data.Metadata.DeviceID,
					"NA", // InterfaceID
					"NA", // PonID
					"NA", // PortNumber
					data.Metadata.Title,
				).Add(data.Metrics.Packets)

				volthaRxBytesTotal.WithLabelValues(
					data.Metadata.LogicalDeviceID,
					data.Metadata.SerialNumber,
					data.Metadata.DeviceID,
					"NA", // InterfaceID
					"NA", // PonID
					"NA", // PortNumber
					data.Metadata.Title,
				).Add(data.Metrics.Octets)
			}

		case "Ethernet_UNI_History":
			// ONU. Do nothing.

		case "FEC_History":
			// ONU. Do Nothing.

			volthaTxBytesTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.TxBytes)

			volthaRxBytesTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.RxBytes)

			volthaTxPacketsTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.TxPackets)

			volthaRxPacketsTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.RxPackets)

			volthaTxErrorPacketsTotal.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.TxErrorPackets)

			volthaRxErrorPacketsTotal.WithLabelValues(
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

		case "voltha.internal":
			// Voltha Internal. Do nothing.
		}
	}
}
