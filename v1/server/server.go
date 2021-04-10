//  ______________________
// < I'm a sandissa kitty >
//  ----------------------
//   \
//    \
//       /\_)o<
//      |      \
//      | O . O|
//      \_____/
package main

import (
	"os"
	"os/signal"
)

func main() {
	var err error

	// Initialize the database
	err = initDB()
	if err != nil {
		panic(err)
	}
	defer closeDB()

	// Initialize MQTT
	clientMQTT, err = getClient()
	if err != nil {
		panic(err)

	}
	subscribe(1, topicTemperature)
	defer clientMQTT.Disconnect(250)

	// Initialize REST
	startRouter()

	// Block until SIGTERM
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
