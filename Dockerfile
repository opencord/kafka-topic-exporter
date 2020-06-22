# Copyright 2018-present Open Networking Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# docker build -t opencord/kafka-topic-exporter:latest .
# docker build -t 10.128.22.1:30500/opencord/kafka-topic-exporter:latest .

FROM golang:1.14-stretch as builder
RUN mkdir -p /go/src/gerrit.opencord.org/kafka-topic-exporter
WORKDIR /go/src/gerrit.opencord.org/kafka-topic-exporter
ADD . /go/src/gerrit.opencord.org/kafka-topic-exporter
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o main .

FROM alpine:3.8
WORKDIR /go/src/gerrit.opencord.org/kafka-topic-exporter/
COPY --from=builder /go/src/gerrit.opencord.org/kafka-topic-exporter/main .
# FIXME this should be mounted by the helm charts
ADD config/conf.yaml /etc/config/conf.yaml
ENTRYPOINT ["./main"]
