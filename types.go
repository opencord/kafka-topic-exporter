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

// KPI Events format

type Metrics struct {
	TxBytes            float64 `json:"tx_bytes"`
	TxPackets          float64 `json:"tx_packets"`
	TxErrorPackets     float64 `json:"tx_error_packets"`
	TxBcastPackets     float64 `json:"tx_bcast_packets"`
	TxUnicastPackets   float64 `json:"tx_ucast_packets"`
	TxMulticastPackets float64 `json:"tx_mcast_packets"`
	RxBytes            float64 `json:"rx_bytes"`
	RxPackets          float64 `json:"rx_packets"`
	RxErrorPackets     float64 `json:"rx_error_packets"`
	RxBcastPackets     float64 `json:"rx_bcast_packets"`
	RxMulticastPackets float64 `json:"rx_mcast_packets"`

	// ONU Ethernet_Bridge_Port_history
	Packets            float64 `json:"packets"`
	Octets             float64 `json:"octets"`
}

type Context struct {
	InterfaceID string `json:"intf_id"`
	PonID       string `json:"pon_id"`
	PortNumber  string `json:"port_no"`

	// ONU Performance Metrics
	ParentClassId string `json:"parent_class_id"`
	ParentEntityId string `json:"parent_entity_id"`
	Upstream    string `json:"upstream"`
}

type Metadata struct {
	LogicalDeviceID string   `json:"logical_device_id"`
	Title           string   `json:"title"`
	SerialNumber    string   `json:"serial_no"`
	Timestamp       float64  `json:"ts"`
	DeviceID        string   `json:"device_id"`
	Context         *Context `json:"context"`
}

type SliceData struct {
	Metrics  *Metrics  `json:"metrics"`
	Metadata *Metadata `json:"metadata"`
}

type VolthaKPI struct {
	Type       string       `json:"type"`
	Timestamp  float64      `json:"ts"`
	SliceDatas []*SliceData `json:"slice_data"`
}

type OnosPort struct {
	PortID        string  `json:"portId"`
	RxPackets     float64 `json:"pktRx"`
	TxPackets     float64 `json:"pktTx"`
	RxBytes       float64 `json:"bytesRx"`
	TxBytes       float64 `json:"bytesTx"`
	RxPacketsDrop float64 `json:"pktRxDrp"`
	TxPacketsDrop float64 `json:"pktTxDrp"`
}

type OnosKPI struct {
	DeviceID string      `json:"deviceId"`
	Ports    []*OnosPort `json:"ports"`
}

type AuthStatKPI struct {
	AuthState           string  `json:"authenticationStatisticsState"`
	AcceptPkt           float64 `json:"acceptPacketCounter"`
	RejectPkt           float64 `json:"rejectPacketCounter"`
	ChallengePkt        float64 `json:"challengePacketCounter"`
	AccessPkt           float64 `json:"accessPacketCounter"`
	InvalidValid        float64 `json:"invalidValidatorCounter"`
	UnknownType         float64 `json:"unknownTypeCounter"`
	PendingRqst         float64 `json:"pendingRequestCounter"`
	DroppedPkt          float64 `json:"numberOfDroppedPackets"`
	MalformedPkt        float64 `json:"malformed_packet_counter"`
	UnknownsrvrPkt      float64 `json:"numberOfPacketFromUnknownServer"`
	RoundTripPkt        float64 `json:"packetRoundtripTimeInMilis"`
}
