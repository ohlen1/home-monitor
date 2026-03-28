FROM --platform=$BUILDPLATFORM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apk add make

RUN mkdir /app
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go mod tidy

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH make build
LABEL maintainer="Andreas Ohlén"

##########

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bin/home-monitor /root/home-monitor
EXPOSE 7071
CMD ["/root/home-monitor"]