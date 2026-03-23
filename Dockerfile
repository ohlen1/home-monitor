FROM golang:alpine AS builder

ENV GOOS=linux
ENV GOARCH=arm64

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
COPY --from=builder /app/bin/home-monitor /root/home-monitor
EXPOSE 7071
CMD ["/root/home-monitor"]