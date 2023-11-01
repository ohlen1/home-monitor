package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"codingminds.com/homemmonitor/config"
	"codingminds.com/homemmonitor/messagehandler"
	"codingminds.com/homemmonitor/metrics"
	"codingminds.com/homemmonitor/model"
	"github.com/gorilla/websocket"
)

const (
	CONNECTION_STATUS_DISCONNECTED = iota
	CONNECTION_STATUS_CONNECTED
)

var connectionStatus = CONNECTION_STATUS_DISCONNECTED
var cfg *config.Config

func main() {
	cfg = config.LoadConfig("config.yaml")
	initLogger(cfg.LogLevel)
	metrics.Init(*cfg)
	url := getWebSocketUrl()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	websocket.DefaultDialer.Subprotocols = []string{"graphql-transport-ws"}
	websocket.DefaultDialer.HandshakeTimeout = time.Second * 20
	headers := http.Header{}
	headers.Add("User-agent", "fortuna-home-monitor/1.0.0")
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.Tibber.ApiToken))

	c, _, err := websocket.DefaultDialer.Dial(url, headers)
	if err != nil {
		slog.Error("Dial error", "error", err)
	}
	defer c.Close()

	done := make(chan struct{})
	msg := make(chan string)

	go func() {
		msg <- createConnectionInitMessage()
		for connectionStatus == CONNECTION_STATUS_DISCONNECTED {
			slog.Info("Waiting for connection...")
			time.Sleep(100 * time.Millisecond)
		}

		msg <- createSubscribeMessage()
	}()

	go func() {
		defer close(done)

		for {
			messageType, message, err := c.ReadMessage()
			slog.Info("New message", "type", messageType)
			if err != nil {
				slog.Error("read:", err)
				return
			}
			if strings.Contains(string(message), "connection_ack") {
				slog.Info("Connection established!")
				connectionStatus = CONNECTION_STATUS_CONNECTED
				continue
			}
			slog.Debug(fmt.Sprintf("recv: %s", message))
			// TODO send message to a channel and call messageHandler in a goroutine
			messagehandler.Handle(message)
		}
	}()

	for {
		slog.Info("Starting loop")
		select {
		case <-done:
			return
		case m := <-msg:
			slog.Debug("Writing message", "message", m)
			err := c.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				slog.Debug("write:", err)
				return
			}
		case <-interrupt:
			slog.Warn("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				slog.Error("write close:", "error", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func getWebSocketUrl() string {
	getUrlRequestTemplate := `{"query":"{viewer{websocketSubscriptionUrl}}"}`
	getUrlRequest := []byte(getUrlRequestTemplate)

	req, err := http.NewRequest("POST", cfg.Tibber.Url, bytes.NewBuffer(getUrlRequest))
	if err != nil {
		slog.Error("Error creating init request", "error", err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cfg.Tibber.ApiToken))
	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	initResponse, err := client.Do(req)
	if err != nil || initResponse.StatusCode != 200 {
		slog.Error("Error executing init request", "error", err, "statusCode", initResponse.StatusCode, "status", initResponse.Status)
		os.Exit(1)
	}
	defer initResponse.Body.Close()

	// parse body
	r := new(model.InitResponseBody)
	json.NewDecoder(initResponse.Body).Decode(r)
	if initResponse.StatusCode != 200 {
		slog.Info("Init request failed", "statusCode", initResponse.StatusCode)
		os.Exit(1)
	}
	slog.Info("Init request successful", "statusCode", initResponse.StatusCode, "wsUrl", r.Data.Viewer.WebsocketSubscriptionUrl)
	return r.Data.Viewer.WebsocketSubscriptionUrl
}

func initLogger(level string) {
	var slogLevel slog.Level
	switch level {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slog.Warn("Log level not recognized. Defaulting to INFO")
		slogLevel = slog.LevelInfo
	}
	logger := *slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))
	slog.SetDefault(&logger)
	slog.Info("Log level set", "level", slogLevel)
}

func createConnectionInitMessage() string {
	msg := new(model.ConnectionInitMessage)
	msg.Type = "connection_init"
	msg.Payload.Token = "q-h07iu6ISI6KBc_wL5ZpmxKZwT-zuZkVqsGE49HZz4"
	b, _ := json.Marshal(msg)
	return string(b)
}

func createSubscribeMessage() string {
	msg := new(model.SubscribeMessage)
	msg.Id = "1"
	msg.Type = "subscribe"
	msg.Payload.Query = fmt.Sprintf("subscription {liveMeasurement(homeId: \"%s\"){%s}}", cfg.Tibber.HomeId, getMeasurementString())
	b, _ := json.Marshal(msg)
	return string(b)
}

func getMeasurementString() string {
	meas := ""
	for i, s := range cfg.Tibber.Measurements {
		meas += s
		if i < len(meas)-1 {
			meas += " "
		}
	}
	return meas
}
