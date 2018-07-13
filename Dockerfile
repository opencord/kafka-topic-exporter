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
# CMD ["-broker=voltha-kafka.default.svc.cluster.local:9092", "-topic=voltha.kpis"]