package main

import (
	"fmt"
	"time"

	"github.com/creativeprojects/go-homie"
)

func main() {
	// defining the homie device
	device := homie.
		NewDevice("raspberry-pi", "Raspberry PI agent").
		AddNode("bme280", "BME280 on GPIO", "bme280").
		AddProperty("temperature", "Temperature", homie.TypeFloat).SetUnit("°C").Node().
		AddProperty("pressure", "Pressure", homie.TypeFloat).SetUnit("hPa").Node().
		AddProperty("humidity", "Humidity", homie.TypeFloat).SetUnit("%").Node().
		Device()

	// get the full homie definition to send to MQTT - you only need to send it once unless it's changing over time
	for _, attribute := range device.GetHomieAttributes() {
		publish(attribute.Topic, attribute.Value)
	}

	// install a global callback on the device
	device.OnSet(onSet)

	for i := 0; i <= 3; i++ {
		// new values will be published (to the console)
		device.Node("bme280").Property("temperature").Set(28 + i)
		device.Node("bme280").Property("humidity").Set(40 + i*10)
		device.Node("bme280").Property("pressure").Set(998 + i)

		device.SetState(homie.StateSleeping)
		time.Sleep(3 * time.Second)
	}

}

func onSet(topic, value string, dataType homie.PropertyType) {
	if value == "<nil>" {
		value = ""
	}
	if value == "" && dataType != homie.TypeString {
		// don't send a blank string on anything else than a string data type
		return
	}
	publish(topic, value)
}

func publish(topic, value string) {
	// display topic/value to the console
	fmt.Printf("%s\t '%s'\n", topic, value)
}
