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
	"gerrit.opencord.org/kafka-topic-exporter/common/logger"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var (
	// voltha kpis
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

	// optical parameters
	VolthaOnuLaserBiasCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_laser_bias_current",
			Help: "ONU Laser bias current value",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOnuTemperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_temperature",
			Help: "ONU temperature value",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	VolthaOnuPowerFeedVoltage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_power_feed_voltage",
			Help: "ONU power feed voltage",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	VolthaOnuMeanOpticalLaunchPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_mean_optical_launch_power",
			Help: "ONU mean optical launch power",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	VolthaOnuReceivedOpticalPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_received_optical_power",
			Help: "ONU received optical power",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	// onos kpis
	onosTxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onos_tx_bytes_total",
			Help: "Number of total bytes transmitted",
		},
		[]string{"device_id", "port_id"},
	)
	onosRxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onos_rx_bytes_total",
			Help: "Number of total bytes received",
		},
		[]string{"device_id", "port_id"},
	)
	onosTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onos_tx_packets_total",
			Help: "Number of total packets transmitted",
		},
		[]string{"device_id", "port_id"},
	)
	onosRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onos_rx_packets_total",
			Help: "Number of total packets received",
		},
		[]string{"device_id", "port_id"},
	)

	onosTxDropPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onos_tx_drop_packets_total",
			Help: "Number of total transmitted packets dropped",
		},
		[]string{"device_id", "port_id"},
	)

	onosRxDropPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onos_rx_drop_packets_total",
			Help: "Number of total received packets dropped",
		},
		[]string{"device_id", "port_id"},
	)

	// onos.aaa kpis
	onosaaaRxAcceptResponses = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_accept_responses",
			Help: "Number of access accept packets received from the server",
		})
	onosaaaRxRejectResponses = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_reject_responses",
			Help: "Number of access reject packets received from the server",
		})
	onosaaaRxChallengeResponses = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_challenge_response",
			Help: "Number of access challenge packets received from the server",
		})
	onosaaaTxAccessRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_tx_access_requests",
			Help: "Number of access request packets sent to the server",
		})
	onosaaaRxInvalidValidators = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_invalid_validators",
			Help: "Number of access response packets received from the server with an invalid validator",
		})
	onosaaaRxUnknownType = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_unknown_type",
			Help: "Number of packets of an unknown RADIUS type received from the accounting server",
		})
	onosaaaPendingRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_pending_responses",
			Help: "Number of access request packets pending a response from the server",
		})
	onosaaaRxDroppedResponses = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_dropped_responses",
			Help: "Number of dropped packets received from the accounting server",
		})
	onosaaaRxMalformedResponses = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_malformed_responses",
			Help: "Number of malformed access response packets received from the server",
		})
	onosaaaRxUnknownserver = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_from_unknown_server",
			Help: "Number of packets received from an unknown server",
		})
	onosaaaRequestRttMillis = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_request_rttmillis",
			Help: "Roundtrip packet time to the accounting server in Miliseconds",
		})
	onosaaaRequestReTx = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_request_re_tx",
			Help: "Number of access request packets retransmitted to the server",
		})

	// ONOS BNG kpis

	// --------------------- BNG UPSTREAM STATISTICS -----------------------------------------
	onosBngUpTxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngUpTxBytes",
			Help: "onosBngUpTxBytes",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)
	onosBngUpTxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngUpTxPackets",
			Help: "onosBngUpTxPackets",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)

	onosBngUpRxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngUpRxBytes",
			Help: "onosBngUpRxBytes",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)

	onosBngUpRxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngUpRxPackets",
			Help: "onosBngUpRxPackets",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)
	onosBngUpDropBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngUpDropBytes",
			Help: "onosBngUpDropBytes",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)
	onosBngUpDropPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngUpDropPackets",
			Help: "onosBngUpDropPackets",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)

	// --------------------- BNG CONTROL STATISTICS ------------------------------------------
	onosBngControlPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngControlPackets",
			Help: "onosBngControlPackets",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)

	// -------------------- BNG DOWNSTREAM STATISTICS ----------------------------------------
	onosBngDownTxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngDownTxBytes",
			Help: "onosBngDownTxBytes",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)
	onosBngDownTxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngDownTxPackets",
			Help: "onosBngDownTxPackets",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)

	onosBngDownRxBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngDownRxBytes",
			Help: "onosBngDownRxBytes",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)
	onosBngDownRxPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngDownRxPackets",
			Help: "onosBngDownRxPackets",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)

	onosBngDownDropBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngDownDropBytes",
			Help: "onosBngDownDropBytes",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)
	onosBngDownDropPackets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "onosBngDownDropPackets",
			Help: "onosBngDownDropPackets",
		},
		[]string{"mac_address", "ip", "session_id", "s_tag", "c_tag", "onu_serial", "type"},
	)

	deviceLaserBiasCurrent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "device_laser_bias_current",
			Help: "Device Laser Bias Current",
		},
		[]string{"port_id"},
	)
	deviceTemperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "device_temperature",
			Help: "Device Temperature",
		},
		[]string{"port_id"},
	)
	deviceTxPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "device_tx_power",
			Help: "Device Tx Power",
		},
		[]string{"port_id"},
	)
	deviceVoltage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "device_voltage",
			Help: "Device Voltage",
		},
		[]string{"port_id"},
	)

	onosaaaRxEapolLogoff = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_eapol_Logoff",
			Help: "Number of EAPOL logoff messages received resulting in disconnected state",
		})
	onosaaaTxEapolResIdentityMsg = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_tx_eapol_Res_IdentityMsg",
			Help: "Number of authenticating transitions due to EAP response or identity message",
		})
	onosaaaTxAuthSuccess = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_tx_auth_Success",
			Help: "Number of authenticated transitions due to successful authentication",
		})
	onosaaaTxAuthFailure = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_tx_auth_Failure",
			Help: "Number of transitions to held due to authentication failure",
		})
	onosaaaTxStartReq = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_tx_start_Req",
			Help: "Number of transitions to connecting due to start request",
		})
	onosaaaEapPktTxAuthChooseEap = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_eap_Pkt_tx_auth_choosing_Eap",
			Help: "Number of EAP request packets sent due to the authenticator choosing the EAP method",
		})
	onosaaaTxRespnotNak = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_tx_Resp_not_Nak",
			Help: "Number of transitions to response (received response other that NAK)",
		})
	onosaaaEapolFramesTx = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_eapol_frames_tx",
			Help: "Number of EAPOL frames transmitted",
		})
	onosaaaAuthStateIdle = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_auth_state_idle",
			Help: "Number of state machine status as Idle",
		})
	onosaaaRequestIdFramesTx = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_tx_request_id_frames",
			Help: "Number of request ID EAP frames transmitted",
		})
	onosaaaRequestEapFramesTx = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_tx_request_eap_frames",
			Help: "Number of request EAP frames transmitted",
		})
	onosaaaInvalidPktType = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_invalid_pkt_type",
			Help: "Number of EAPOL frames received with invalid frame(Packet) type",
		})
	onosaaaInvalidBodyLength = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_invalid_body_length",
			Help: "Number of EAPOL frames received with invalid body length",
		})
	onosaaaValidEapolFramesRx = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_valid_eapol_frames",
			Help: "Number of valid EAPOL frames received",
		})
	onosaaaPendingResSupplicant = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_pending_response_supplicant",
			Help: "Number of request pending response from supplicant",
		})
	onosaaaRxResIdEapFrames = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "onosaaa_rx_res_id_eap_frames",
			Help: "Number of response ID EAP frames received",
		})
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

		case "PON_Optical":
			VolthaOnuLaserBiasCurrent.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.LaserBiasCurrent)

			volthaOnuTemperature.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.Temperature)

			VolthaOnuPowerFeedVoltage.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.PowerFeedVoltage)

			VolthaOnuMeanOpticalLaunchPower.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.MeanOpticalLaunchPower)

			VolthaOnuReceivedOpticalPower.WithLabelValues(
				data.Metadata.LogicalDeviceID,
				data.Metadata.SerialNumber,
				data.Metadata.DeviceID,
				data.Metadata.Context.InterfaceID,
				data.Metadata.Context.PonID,
				data.Metadata.Context.PortNumber,
				data.Metadata.Title,
			).Set(data.Metrics.ReceivedOpticalPower)

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

func exportOnosKPI(kpi OnosKPI) {

	for _, data := range kpi.Ports {

		onosTxBytesTotal.WithLabelValues(
			kpi.DeviceID,
			data.PortID,
		).Set(data.TxBytes)

		onosRxBytesTotal.WithLabelValues(
			kpi.DeviceID,
			data.PortID,
		).Set(data.RxBytes)

		onosTxPacketsTotal.WithLabelValues(
			kpi.DeviceID,
			data.PortID,
		).Set(data.TxPackets)

		onosRxPacketsTotal.WithLabelValues(
			kpi.DeviceID,
			data.PortID,
		).Set(data.RxPackets)

		onosTxDropPacketsTotal.WithLabelValues(
			kpi.DeviceID,
			data.PortID,
		).Set(data.TxPacketsDrop)

		onosRxDropPacketsTotal.WithLabelValues(
			kpi.DeviceID,
			data.PortID,
		).Set(data.RxPacketsDrop)
	}
}

func exportImporterKPI(kpi ImporterKPI) {
	deviceLaserBiasCurrent.WithLabelValues(
		kpi.PortId,
	).Set(kpi.LaserBiasCurrent)
	deviceTemperature.WithLabelValues(
		kpi.PortId,
	).Set(kpi.Temperature)
	deviceTxPower.WithLabelValues(
		kpi.PortId,
	).Set(kpi.TxPower)
	deviceVoltage.WithLabelValues(
		kpi.PortId,
	).Set(kpi.Voltage)
}

func exportOnosAaaKPI(kpi OnosAaaKPI) {

	onosaaaRxAcceptResponses.Set(kpi.RxAcceptResponses)

	onosaaaRxRejectResponses.Set(kpi.RxRejectResponses)

	onosaaaRxChallengeResponses.Set(kpi.RxChallengeResponses)

	onosaaaTxAccessRequests.Set(kpi.TxAccessRequests)

	onosaaaRxInvalidValidators.Set(kpi.RxInvalidValidators)

	onosaaaRxUnknownType.Set(kpi.RxUnknownType)

	onosaaaPendingRequests.Set(kpi.PendingRequests)

	onosaaaRxDroppedResponses.Set(kpi.RxDroppedResponses)

	onosaaaRxMalformedResponses.Set(kpi.RxMalformedResponses)

	onosaaaRxUnknownserver.Set(kpi.RxUnknownserver)

	onosaaaRequestRttMillis.Set(kpi.RequestRttMillis)

	onosaaaRequestReTx.Set(kpi.RequestReTx)

	onosaaaRxEapolLogoff.Set(kpi.RxEapolLogoff)

	onosaaaTxEapolResIdentityMsg.Set(kpi.TxEapolResIdentityMsg)

	onosaaaTxAuthSuccess.Set(kpi.TxAuthSuccess)

	onosaaaTxAuthFailure.Set(kpi.TxAuthFailure)

	onosaaaTxStartReq.Set(kpi.TxStartReq)

	onosaaaEapPktTxAuthChooseEap.Set(kpi.EapPktTxAuthChooseEap)

	onosaaaTxRespnotNak.Set(kpi.TxResponseNotNak)

	onosaaaEapolFramesTx.Set(kpi.EapolFramesTx)

	onosaaaAuthStateIdle.Set(kpi.AuthStateIdle)

	onosaaaRequestIdFramesTx.Set(kpi.RequestIdFramesTx)

	onosaaaRequestEapFramesTx.Set(kpi.RequestEapFramesTx)

	onosaaaInvalidPktType.Set(kpi.InvalidPktType)

	onosaaaInvalidBodyLength.Set(kpi.InvalidBodyLength)

	onosaaaValidEapolFramesRx.Set(kpi.ValidEapolFramesRx)

	onosaaaPendingResSupplicant.Set(kpi.PendingResSupplicant)

	onosaaaRxResIdEapFrames.Set(kpi.RxResIdEapFrames)
}

func exportOnosBngKPI(kpi OnosBngKPI) {
	logger.WithFields(log.Fields{
		"Mac":             kpi.Mac,
		"Ip":              kpi.Ip,
		"PppoeSessionId":  strconv.Itoa(kpi.PppoeSessionId),
		"AttachmentType":  kpi.AttachmentType,
		"STag":            strconv.Itoa(kpi.STag),
		"CTag":            strconv.Itoa(kpi.CTag),
		"onuSerialNumber": kpi.OnuSerialNumber,
	}).Trace("Received OnosBngKPI message")

	if kpi.UpTxBytes != nil {
		onosBngUpTxBytes.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.UpTxBytes)
	}
	if kpi.UpTxPackets != nil {
		onosBngUpTxPackets.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.UpTxPackets)
	}
	if kpi.UpRxBytes != nil {
		onosBngUpRxBytes.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.UpRxBytes)
	}
	if kpi.UpRxPackets != nil {
		onosBngUpRxPackets.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.UpRxPackets)
	}

	if kpi.UpDropBytes != nil {
		onosBngUpDropBytes.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.UpDropBytes)
	}
	if kpi.UpDropPackets != nil {
		onosBngUpDropPackets.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.UpDropPackets)
	}

	if kpi.ControlPackets != nil {
		onosBngControlPackets.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.ControlPackets)
	}

	if kpi.DownTxBytes != nil {
		onosBngDownTxBytes.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.DownTxBytes)
	}
	if kpi.DownTxPackets != nil {
		onosBngDownTxPackets.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.DownTxPackets)
	}

	if kpi.DownRxBytes != nil {
		onosBngDownRxBytes.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.DownRxBytes)
	}
	if kpi.DownRxPackets != nil {
		onosBngDownRxPackets.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.DownRxPackets)
	}

	if kpi.DownDropBytes != nil {
		onosBngDownDropBytes.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.DownDropBytes)
	}
	if kpi.DownDropPackets != nil {
		onosBngDownDropPackets.WithLabelValues(
			kpi.Mac,
			kpi.Ip,
			strconv.Itoa(kpi.PppoeSessionId),
			strconv.Itoa(kpi.STag),
			strconv.Itoa(kpi.CTag),
			kpi.OnuSerialNumber,
			kpi.AttachmentType,
		).Set(*kpi.DownDropPackets)
	}
}

func export(topic *string, data []byte) {
	switch *topic {
	case "voltha.kpis":
		kpi := VolthaKPI{}
		err := json.Unmarshal(data, &kpi)
		if err != nil {
			logger.Error("Invalid msg on voltha.kpis: %s, Unprocessed Msg: %s", err.Error(), string(data))
			break
		}
		exportVolthaKPI(kpi)
	case "onos.kpis":
		kpi := OnosKPI{}
		err := json.Unmarshal(data, &kpi)
		if err != nil {
			logger.Error("Invalid msg on onos.kpis: %s, Unprocessed Msg: %s", err.Error(), string(data))
			break
		}
		exportOnosKPI(kpi)
	case "importer":
		kpi := ImporterKPI{}
		strData := string(data)
		idx := strings.Index(strData, "{")
		strData = strData[idx:]

		var m map[string]interface{}
		err := json.Unmarshal([]byte(strData), &m)
		if err != nil {
			logger.Error("Invalid msg on importer: %s", err.Error())
			logger.Debug("Unprocessed Msg: %s", strData)
			break
		}
		if val, ok := m["TransceiverStatistics"]; ok {
			stats := val.(map[string]interface{})
			kpi.LaserBiasCurrent = stats["BiasCurrent"].(map[string]interface{})["Reading"].(float64)
			kpi.Temperature = stats["Temperature"].(map[string]interface{})["Reading"].(float64)
			kpi.TxPower = stats["TxPower"].(map[string]interface{})["Reading"].(float64)
			kpi.Voltage = stats["Voltage"].(map[string]interface{})["Reading"].(float64)
		} else {
			logger.Error("Optical stats (TransceiverStatistics) information missing [topic=importer")
			logger.Debug("Unprocessed Msg: %s", strData)
			break
		}
		kpi.PortId = m["Id"].(string)
		exportImporterKPI(kpi)
	case "onos.aaa.stats.kpis":
		kpi := OnosAaaKPI{}
		err := json.Unmarshal(data, &kpi)
		if err != nil {
			logger.Error("Invalid msg on onos.aaa.stats.kpis: %s, Unprocessed Msg: %s", err.Error(), string(data))
			break
		}
		exportOnosAaaKPI(kpi)
	case "bng.stats":
		kpi := OnosBngKPI{}
		err := json.Unmarshal(data, &kpi)
		if err != nil {
			logger.Error("Invalid msg on bng.stats: %s, Unprocessed Msg: %s", err.Error(), string(data))
			break
		}
		exportOnosBngKPI(kpi)
	default:
		logger.Warn("Unexpected export. Topic [%s] not supported. Should not come here", *topic)
	}
}
