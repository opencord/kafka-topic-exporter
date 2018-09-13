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
# docker build -t 10.90.0.101:30500/opencord/kafka-topic-exporter:latest .

FROM golang:latest as builder
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/Shopify/sarama
RUN CGO_ENABLED=0 GOOS=linux go build -o main . 

FROM alpine:latest  
WORKDIR /app/
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]