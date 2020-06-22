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
	"strconv"
	"strings"

	"gerrit.opencord.org/kafka-topic-exporter/common/logger"
	"github.com/golang/protobuf/proto"
	"github.com/opencord/device-management-interface/go/dmi"
	"github.com/opencord/voltha-protos/go/voltha"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	// voltha kpis
	volthaOltTxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_tx_bytes_total",
			Help: "Number of total bytes transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	volthaOltRxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_rx_bytes_total",
			Help: "Number of total bytes received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	volthaOltTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_tx_packets_total",
			Help: "Number of total packets transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)
	volthaOltRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_rx_packets_total",
			Help: "Number of total packets received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOltTxErrorPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_tx_error_packets_total",
			Help: "Number of total transmitted packets error",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOltRxErrorPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_rx_error_packets_total",
			Help: "Number of total received packets error",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOltTxBroadcastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_tx_broadcast_packets_total",
			Help: "Number of total broadcast packets transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOltTxUnicastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_tx_unicast_packets_total",
			Help: "Number of total unicast packets transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOltTxMulticastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_tx_multicast_packets_total",
			Help: "Number of total multicast packets transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOltRxBroadcastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_rx_broadcast_packets_total",
			Help: "Number of total broadcast packets received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOltRxUnicastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_rx_unicast_packets_total",
			Help: "Number of total unicast packets received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	volthaOltRxMulticastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_olt_rx_multicast_packets_total",
			Help: "Number of total multicast packets received",
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

	VolthaOnuTransmtOpticalPower = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_transmit_optical_power",
			Help: "ONU transmited optical power",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "pon_id", "port_number", "title"},
	)

	// FEC parameters
	volthaOnuFecCorrectedCodewordsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_fec_corrected_code_words",
			Help: "Number of total code words corrected",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuFecCodewordsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_fec_code_words_total",
			Help: "Number of total code words",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuFecCorrectedBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_fec_corrected_bytes_total",
			Help: "Number of total corrected bytes",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuFecSecondsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_fec_corrected_fec_seconds_total",
			Help: "Number of fec seconds total",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuFecUncorrectablewordsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_fec_uncorrectable_words_total",
			Help: "Number of fec uncorrectable words",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	//Etheret UNI

	volthaEthernetUniSingleCollisionTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_single_collision_frame_counter",
			Help: "successfully transmitted frames but delayed by exactly one collision.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)

	volthaEthernetUniMacLayerTramsmitErrorTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_internal_mac_rx_error_counter",
			Help: "transmission failed due to an internal MAC sublayer transmit error.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)

	volthaEthernetUniMultiCollisionTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_multiple_collisions_frame_counter",
			Help: "successfully transmitted frames but delayed by multiple collisions.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)

	volthaEthernetUniFramestooLongTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_frames_too_long",
			Help: "frames that exceeded the maximum permitted frame size.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)

	volthaEthernetUniAlignmentErrorTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_alignment_error_counter",
			Help: "frames that were not an integral number of octets in length and did not pass the FCS check.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	volthaEthernetUniCarrierErrorTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_carrier_sense_error_counter",
			Help: "number of times that carrier sense was lost or never asserted when attempting to transmit a frame.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	volthaEthernetUniExcessiveCollisionErrorTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_excessive_collision_counter",
			Help: "frames whose transmission failed due to excessive collisions.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	volthaEthernetUniDeferredTxTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_deferred_tx_counter",
			Help: "frames whose first transmission attempt was delayed because the medium was busy.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	volthaEthernetUniLateCollisionTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_late_collision_counter",
			Help: "number of times that a collision was detected later than 512 bit times into the transmission of a packet.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	volthaEthernetUniBufferOverflowsRxErrorTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_buffer_overflows_on_rx",
			Help: "number of times that the receive buffer overflowed.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	volthaEthernetUniFcsErrorTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_fcs_errors",
			Help: " frames failed the frame check sequence (FCS) check.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	volthaEthernetUniSqeErrorTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_sqe_counter",
			Help: "number of times that the SQE test error message was generated by the PLS sublayer",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	volthaEthernetUniBufferOverflowsTxErrorTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_ethernet_uni_buffer_overflows_on_tx",
			Help: " number of times that the transmit buffer overflowed.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "interface_id", "port_number", "title"},
	)
	//Ethernet_Bridge_Port

	volthaOnuBridgePortTxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_tx_bytes_total",
			Help: "Number of total bytes transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortRxBytesTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_rx_bytes_total",
			Help: "Number of total bytes received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_tx_packets_total",
			Help: "Number of total packets transmitted",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_rx_packets_total",
			Help: "Number of total packets received",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_64octetTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_64_octets_Txpackets",
			Help: "packets (including bad packets) that were 64 octets long",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_65_127_octetTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_65_to_127_octet_Txpackets",
			Help: "packets (including bad packets) that were 65..127 octets long, excluding framing bits but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_128_255_octetTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_128_to_255_octet_Txpackets",
			Help: "packets (including bad packets) received that were 128..255 octets long, excluding framing bits but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_256_511_octetTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_256_to_511_octet_Txpackets",
			Help: "packets (including bad packets) received that were 256..511 octets long, excluding framing bits but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_512_1023_octetTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_512_to_1023_octet_Txpackets",
			Help: "packets (including bad packets) received that were 512..1 023 octets long, excluding framing bits but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_1024_1518_octetTxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_1024_to_1518_octet_Txpackets",
			Help: "packets (including bad packets) received that were 1024..1518 octets long, excluding framing bits, but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortTxMulticastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_multicast_Txpackets",
			Help: "packets received that were directed to a multicast address.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortTxBroadcastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_broadcast_Txpackets",
			Help: "packets received that were directed to the broadcast address.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortTxOversizePacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_oversize_Txpackets",
			Help: " packets received that were longer than 1518 octets",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortTxCrcErrorPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_crc_errored_Txpackets",
			Help: "Packets with CRC errors",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortTxUndersizePacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_undersize_Txpackets",
			Help: "Packets received that were less than 64 octets long, but were otherwise well formed",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePortTxDropEventsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_Txdrop_events",
			Help: "total number of events in which packets were dropped due to a lack of resources. ",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePort_64octetRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_64_octets_Rxpackets",
			Help: "packets (including bad packets) that were 64 octets long",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_65_127_octetRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_65_to_127_octet_Rxpackets",
			Help: "packets (including bad packets) that were 65..127 octets long, excluding framing bits but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_128_255_octetRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_128_to_255_octet_packets",
			Help: "packets (including bad packets) received that were 128..255 octets long, excluding framing bits but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_256_511_octetRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_256_to_511_octet_Rxpackets",
			Help: "packets (including bad packets) received that were 256..511 octets long, excluding framing bits but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_512_1023_octetRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_512_to_1023_octet_Rxpackets",
			Help: "packets (including bad packets) received that were 512..1 023 octets long, excluding framing bits but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePort_1024_1518_octetRxPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_1024_to_1518_octet_Rxpackets",
			Help: "packets (including bad packets) received that were 1024..1518 octets long, excluding framing bits, but including FCS.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortRxMulticastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_multicast_Rxpackets",
			Help: "packets received that were directed to a multicast address.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortRxBroadcastPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_broadcast_Rxpackets",
			Help: "packets received that were directed to the broadcast address.",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortRxOversizePacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_oversize_Rxpackets",
			Help: " packets received that were longer than 1518 octets",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortRxCrcErrorPacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_crc_errored_Rxpackets",
			Help: "Packets with CRC errors",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)
	volthaOnuBridgePortRxUndersizePacketsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_undersize_Rxpackets",
			Help: "Packets received that were less than 64 octets long, but were otherwise well formed",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
	)

	volthaOnuBridgePortRxDropEventsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "voltha_onu_bridge_port_Rxdrop_events",
			Help: "total number of events in which packets were dropped due to a lack of resources. ",
		},
		[]string{"logical_device_id", "serial_number", "device_id", "title"},
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
	/* The device metrics will be removed in future and device
	   metrics defined in VOL-3255 will be supported
	*/
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
	//OLT Device Metrics
	//TODO: Check if component level temperatures are supported by Devices,If not remove in later versions of exporter
	oltDeviceCpuTemp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_cpu_temperature",
			Help: "cpu temperature",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)

	oltDeviceCpuUsagePercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_cpu_usage_percentage",
			Help: "usage of cpu",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)
	oltDeviceFanSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_fan_speed",
			Help: "fan speed",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)
	oltDeviceDiskTemp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_disk_temp",
			Help: "disk temperature",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)
	oltDeviceDiskUsagePercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_disk_usage_percent",
			Help: "disk usage percentage",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)
	oltDeviceRamTemp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_ram_temp",
			Help: "RAM temperature",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)
	oltDeviceRamUsagePercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_ram_usage_percentage",
			Help: "RAM usage percentage",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)

	oltDevicePowerUsagePercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_power_usage_percentage",
			Help: "power usage percentage",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)

	oltDeviceInnerSurroundTemp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_inner_surrounding_temperature",
			Help: "inner surrounding temperature",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)

	oltDevicePowerUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "olt_device_power_usage",
			Help: "power usage",
		},
		[]string{"deviceuuid", "componentuuid", "componentname"},
	)
)

var oltDeviceMetrics = map[dmi.MetricNames]*prometheus.GaugeVec{
	dmi.MetricNames_METRIC_CPU_TEMP:               oltDeviceCpuTemp,
	dmi.MetricNames_METRIC_CPU_USAGE_PERCENTAGE:   oltDeviceCpuUsagePercent,
	dmi.MetricNames_METRIC_FAN_SPEED:              oltDeviceFanSpeed,
	dmi.MetricNames_METRIC_DISK_TEMP:              oltDeviceDiskTemp,
	dmi.MetricNames_METRIC_DISK_USAGE_PERCENTAGE:  oltDeviceDiskUsagePercent,
	dmi.MetricNames_METRIC_RAM_TEMP:               oltDeviceRamTemp,
	dmi.MetricNames_METRIC_RAM_USAGE_PERCENTAGE:   oltDeviceRamUsagePercent,
	dmi.MetricNames_METRIC_POWER_USAGE_PERCENTAGE: oltDevicePowerUsagePercent,
	dmi.MetricNames_METRIC_INNER_SURROUNDING_TEMP: oltDeviceInnerSurroundTemp,
	dmi.MetricNames_METRIC_POWER_USAGE:            oltDevicePowerUsage,
}

func exportVolthaEthernetPonStats(data *voltha.MetricInformation) {

	volthaOltTxBytesTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["TxBytes"]))

	volthaOltRxBytesTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["RxBytes"]))

	volthaOltTxPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["TxPackets"]))

	volthaOltRxPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["RxPackets"]))

	volthaOltTxErrorPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["TxErrorPackets"]))

	volthaOltRxErrorPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["RxErrorPackets"]))

	volthaOltTxBroadcastPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["TxBcastPackets"]))

	volthaOltTxUnicastPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["TxUcastPackets"]))

	volthaOltTxMulticastPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["TxMcastPackets"]))

	volthaOltRxBroadcastPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["RxBcastPackets"]))

	volthaOltRxUnicastPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["RxUcastPackets"]))

	volthaOltRxMulticastPacketsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		"NA", // InterfaceID
		"NA", // PonID
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["RxMcastPackets"]))
}

func exportVolthaOnuEthernetBridgePortStats(data *voltha.MetricInformation) {

	if (data.GetMetadata().GetContext()["upstream"]) == "True" {
		// ONU. Extended Ethernet statistics.
		volthaOnuBridgePortTxPacketsTotal.WithLabelValues(
			data.Metadata.GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["packets"]))

		volthaOnuBridgePortTxBytesTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.GetMetadata().GetSerialNo(),
			data.GetMetadata().GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["octets"]))

		volthaOnuBridgePort_64octetTxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["64_octets"]))

		volthaOnuBridgePort_65_127_octetTxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["65_to_127_octets"]))

		volthaOnuBridgePort_128_255_octetTxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["128_to_255_octets"]))

		volthaOnuBridgePort_256_511_octetTxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["256_to_511_octets"]))

		volthaOnuBridgePort_512_1023_octetTxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["512_to_1023_octets"]))

		volthaOnuBridgePort_1024_1518_octetTxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["1024_to_1518_octets"]))

		volthaOnuBridgePortTxMulticastPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["multicast_packets"]))

		volthaOnuBridgePortTxBroadcastPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["broadcast_packets"]))

		volthaOnuBridgePortTxOversizePacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["oversize_packets"]))

		volthaOnuBridgePortTxCrcErrorPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["crc_errored_packets"]))

		volthaOnuBridgePortTxUndersizePacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["undersize_packets"]))

		volthaOnuBridgePortTxDropEventsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["drop_events"]))

	} else {

		// ONU. Extended Ethernet statistics.
		volthaOnuBridgePortRxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.GetMetadata().GetSerialNo(),
			data.GetMetadata().GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["packets"]))

		volthaOnuBridgePortRxBytesTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.GetMetadata().GetSerialNo(),
			data.GetMetadata().GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["octets"]))

		volthaOnuBridgePort_64octetRxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["64_octets"]))

		volthaOnuBridgePort_65_127_octetRxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["65_to_127_octets"]))

		volthaOnuBridgePort_128_255_octetRxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["128_to_255_octets"]))

		volthaOnuBridgePort_256_511_octetRxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["256_to_511_octets"]))

		volthaOnuBridgePort_512_1023_octetRxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["512_to_1023_octets"]))

		volthaOnuBridgePort_1024_1518_octetRxPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["1024_to_1518_octets"]))

		volthaOnuBridgePortRxMulticastPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["multicast_packets"]))

		volthaOnuBridgePortRxBroadcastPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["broadcast_packets"]))

		volthaOnuBridgePortRxOversizePacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["oversize_packets"]))

		volthaOnuBridgePortRxCrcErrorPacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["crc_errored_packets"]))

		volthaOnuBridgePortRxUndersizePacketsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["undersize_packets"]))

		volthaOnuBridgePortRxDropEventsTotal.WithLabelValues(
			data.GetMetadata().GetLogicalDeviceId(),
			data.Metadata.GetSerialNo(),
			data.Metadata.GetDeviceId(),
			data.GetMetadata().GetTitle(),
		).Add(float64(data.GetMetrics()["drop_events"]))

	}
}

func exportVolthaOnuPonOpticalStats(data *voltha.MetricInformation) {
	VolthaOnuTransmtOpticalPower.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		"NA", // PonID,
		"NA", //PortNumber
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["transmit_power"]))

	VolthaOnuReceivedOpticalPower.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		"NA", // PonID,
		"NA", //PortNumber
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["receive_power"]))
}
func exportVolthaOnuFecStats(data *voltha.MetricInformation) {
	volthaOnuFecCorrectedCodewordsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["corrected_code_words"]))

	volthaOnuFecCodewordsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["total_code_words"]))

	volthaOnuFecCorrectedBytesTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["corrected_bytes"]))

	volthaOnuFecSecondsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["fec_seconds"]))

	volthaOnuFecUncorrectablewordsTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["uncorrectable_code_words"]))
}
func exportVolthaOnuEthernetUniStats(data *voltha.MetricInformation) {

	volthaEthernetUniSingleCollisionTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["single_collision_frame_counter"]))

	volthaEthernetUniMacLayerTramsmitErrorTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["internal_mac_rx_error_counter"]))

	volthaEthernetUniMultiCollisionTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["multiple_collisions_frame_counter"]))

	volthaEthernetUniFramestooLongTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["frames_too_long"]))
	volthaEthernetUniAlignmentErrorTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["alignment_error_counter"]))

	volthaEthernetUniCarrierErrorTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["carrier_sense_error_counter"]))
	volthaEthernetUniExcessiveCollisionErrorTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["excessive_collision_counter"]))

	volthaEthernetUniDeferredTxTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["deferred_tx_counter"]))

	volthaEthernetUniLateCollisionTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["late_collision_counter"]))

	volthaEthernetUniBufferOverflowsRxErrorTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()[""]))

	volthaEthernetUniFcsErrorTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["fcs_errors"]))

	volthaEthernetUniSqeErrorTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["sqe_counter"]))

	volthaEthernetUniBufferOverflowsTxErrorTotal.WithLabelValues(
		data.GetMetadata().GetLogicalDeviceId(),
		data.GetMetadata().GetSerialNo(),
		data.GetMetadata().GetDeviceId(),
		data.GetMetadata().GetContext()["intf_id"],
		data.GetMetadata().GetContext()["portno"],
		data.GetMetadata().GetTitle(),
	).Set(float64(data.GetMetrics()["buffer_overflows_on_tx"]))

}

func exportVolthaKPIevent2(kpi *voltha.KpiEvent2) {
	for _, data := range kpi.GetSliceData() {
		switch title := data.GetMetadata().GetTitle(); title {
		case "ETHERNET_NNI", "PON_OLT":
			exportVolthaEthernetPonStats(data)
		case "Ethernet_Bridge_Port_History":
			exportVolthaOnuEthernetBridgePortStats(data)
		case "PON_Optical":
			exportVolthaOnuPonOpticalStats(data)
		case "Ethernet_UNI_History":
			exportVolthaOnuEthernetUniStats(data)
		case "FEC_History":
			exportVolthaOnuFecStats(data)
		case "UNI_Status":
			//  Do nothing.

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

func exportDeviceKPI(kpi *dmi.Metric) {

	if metrics, ok := oltDeviceMetrics[kpi.GetMetricId()]; ok {
		metrics.WithLabelValues(
			kpi.GetMetricMetadata().GetDeviceUuid().GetUuid(),
			kpi.GetMetricMetadata().GetComponentUuid().GetUuid(),
			kpi.GetMetricMetadata().GetComponentName(),
		).Set(float64(kpi.GetValue().GetValue()))
	}
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
	case "voltha.events":
		event := voltha.Event{}
		err := proto.Unmarshal(data, &event)
		if err != nil {
			logger.Error("Invalid msg on voltha.kpis: %s, Unprocessed Msg: %s", err.Error(), string(data))
			break
		}
		if event.GetHeader().GetType() == voltha.EventType_KPI_EVENT2 {
			logger.Debug("KPI_EVENT2 received on voltha.events")
			kpiEvent2 := event.GetKpiEvent2()
			exportVolthaKPIevent2(kpiEvent2)
		}

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
	case "dm.metrics":
		kpi := dmi.Metric{}
		err := proto.Unmarshal(data, &kpi)
		if err != nil {
			logger.Error("Invalid msg on device topic : %s, Unprocessed Msg: %s", err.Error(), string(data))
			break
		}
		exportDeviceKPI(&kpi)
	default:
		logger.Warn("Unexpected export. Topic [%s] not supported. Should not come here", *topic)
	}
}
