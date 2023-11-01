FROM golang:alpine as builder

RUN apk add make

RUN mkdir /app
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go mod tidy

RUN make build
LABEL maintainer="Andreas Ohlén"

##########

FROM alpine:latest
WORKDIR /root/
COPY config.yaml /root/config.yaml
COPY --from=builder /app/bin/home-monitor /root/home-monitor
EXPOSE 7071
CMD ["/root/home-monitor"]