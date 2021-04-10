package main

import (
	"fmt"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mqttHandler(msg mqtt.Message) {
	if msg.Topic() == topicTemperature {
		temp, err := strconv.ParseFloat(string(msg.Payload()), 64)
		if err != nil {
			lerr("Bad temperature payload", err, params{
				"payload": string(msg.Payload()),
			})
			return
		}
		logTemperature(temp)
	}
}

func logTemperature(temp float64) {
	err := addTempDB(temp)
	if err != nil {
		lerr("Failed to log temperature", err, params{
			"temp": temp,
		})
		return
	}
	l(fmt.Sprintf("Logged temperature: %f\n", temp))
}
