package model

import "time"

type LiveMeasurementResponseBody struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Payload struct {
		Data struct {
			LiveMeasurement struct {
				Timestamp              time.Time `json:"timestamp"`
				Power                  float64   `json:"power"`
				AveragePower           float64   `json:"averagePower"`
				MinPower               float64   `json:"minPower"`
				MaxPower               float64   `json:"maxPower"`
				AccumulatedConsumption float64   `json:"accumulatedConsumption"`
				AccumulatedCost        float64   `json:"accumulatedCost"`
				CurrentL1              float64   `json:"currentL1"`
				CurrentL2              float64   `json:"currentL2"`
				CurrentL3              float64   `json:"currentL3"`
				VoltagePhase1          float64   `json:"voltagePhase1"`
				VoltagePhase2          float64   `json:"voltagePhase2"`
				VoltagePhase3          float64   `json:"voltagePhase3"`
			} `json:"liveMeasurement"`
		} `json:"data"`
	} `json:"payload"`
}

type InitResponseBody struct {
	Data struct {
		Viewer struct {
			WebsocketSubscriptionUrl string `json:"websocketSubscriptionUrl"`
		}
	}
}

type ConnectionInitMessage struct {
	Type    string `json:"type"`
	Payload struct {
		Token string `json:"token"`
	} `json:"payload"`
}

type SubscribeMessage struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Payload struct {
		Variables     struct{} `json:"variables"`
		Extensions    struct{} `json:"extensions"`
		OperationName *string  `json:"operationName"`
		Query         string   `json:"query"`
	} `json:"payload"`
}
