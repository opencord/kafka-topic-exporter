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

	Acceptpktstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "access_accept_pkt_stats",
			Help: "Number of access accept packets sent from the server",
		},
		[]string{"AuthState"},
	)
	Rejectpktstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "access_reject_pkt_stats",
			Help: "Number of access reject packets sent from the server",
		},
		[]string{"AuthState"},
	)
	Challengepktstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "access_challenge_pkt_stats",
			Help: "Number of access challenge packets sent from the server",
		},
		[]string{"AuthState"},
	)
	Reqstpktstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "access_reqst_pkt_stats",
			Help: "Number of access request packets sent to the server",
		},
		[]string{"AuthState"},
	)
	InvalidValidatorstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "access_invalid_validator_pkt_stats",
			Help: "Number of access response packets received from the server with an invalid validator",
		},
		[]string{"AuthState"},
	)
	UnknownTypestats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "unknown_radiustype_pkt_stats",
			Help: "Number of packets of an unknown RADIUS type received from the accounting server",
		},
		[]string{"AuthState"},
	)
	PendingRequeststats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "access_pndng_response_pkt_stats",
			Help: "Number of access request packets pending a response from the server",
		},
		[]string{"AuthState"},
	)
	DroppedPktstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dropped_pkt_rcvd_stats",
			Help: "Number of dropped packets received from the accounting server",
		},
		[]string{"AuthState"},
	)
	MalformedRespstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "access_malformed_response_pkt_stats",
			Help: "Number of malformed access response packets received from the server",
		},
		[]string{"AuthState"},
	)
	Unknownserverstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "unknown_server_response_pkt_stats",
			Help: "Number of packets received from an unknown server",
		},
		[]string{"AuthState"},
	)
	RoundtripPktstats = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "roundtrip_pkt_time_stats",
			Help: "Roundtrip packet time to the accounting server in Miliseconds",
		},
		[]string{"AuthState"},
	)
)

func exportAuthStatKPI(kpi AuthStatKPI) {

	Acceptpktstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.AcceptPkt)

	Rejectpktstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.RejectPkt)

	Challengepktstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.ChallengePkt)

	Reqstpktstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.AccessPkt)

	InvalidValidatorstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.InvalidValid)

	UnknownTypestats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.UnknownType)

	PendingRequeststats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.PendingRqst)

	DroppedPktstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.DroppedPkt)

	MalformedRespstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.MalformedPkt)

	Unknownserverstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.UnknownsrvrPkt)

	RoundtripPktstats.WithLabelValues(
		"kpi.AuthState",
	).Set(kpi.RoundTripPkt)

}
