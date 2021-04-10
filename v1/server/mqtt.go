package main

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker   = "mqtt.sandyuraz.com"
	port     = 1883
	clientID = "sandissa-secret-kitty"
	username = "kitty"

	topicTemperature = "sensors/Temp"
	topicLED         = "controls/LED"
)

var (
	clientMQTT mqtt.Client
	//go:embed x509/ca.crt
	MQTTCAfile []byte
	//go:embed x509/client.crt
	MQTTClientCert []byte
	//go:embed x509/client.key
	MQTTClientKey []byte
)

/// getClient return MQTT client and makes sure it connects.
func getClient() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	opts.SetUsername(username)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	tls, err := getTLS()
	if err != nil {
		return nil, err
	}
	opts.SetTLSConfig(tls)
	client := mqtt.NewClient(opts)
	// Automatically check the connection
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func getTLS() (*tls.Config, error) {
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(MQTTCAfile)
	clientKeyPair, err := tls.X509KeyPair(MQTTClientCert, MQTTClientKey)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{clientKeyPair},
	}, nil
}

// subscribe subscribes to a topic.
func subscribe(QoS byte, topic string) {
	token := clientMQTT.Subscribe(topic, QoS, nil)
	token.Wait()
	l("Subscribed to topic: " + topic)
}

// publish publishes a payload to a topic.
func publish(QoS byte, topic string, payload string) {
	token := clientMQTT.Publish(topic, QoS, false, payload)
	token.Wait()
}

// messagePubHandler handles all payload from subscribed topics.
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	lf("Received message", params{
		"payload": string(msg.Payload()),
		"topic":   string(msg.Topic()),
	})
	mqttHandler(msg)
}

// connectHandler triggers when connection is successful
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	l("Successfully connected: " + broker)
}

// connectLostHandler triggers when connection is lost
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	l("Connection lost:" + err.Error())
}
