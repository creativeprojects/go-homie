[![Go Report Card](https://goreportcard.com/badge/github.com/creativeprojects/go-homie)](https://goreportcard.com/report/github.com/creativeprojects/go-homie)

# Homie convention for Go

The [Homie convention](https://homieiot.github.io/) defines a standardized way of how IoT devices and services announce themselves and their data on the MQTT broker.

This library is MQTT implementation agnostic: it only generates and manages topics and values, not actually sending them through the wire.

See an example on how to define a Homie device:

``` go
device := homie.NewDevice("my-sensor", "MQTT ESP8266 agent")
device.
    AddNode("bme280", "BME280 via ESP8266EX", "bme280").
    AddProperty("temperature", "Temperature", homie.TypeFloat).SetUnit("°C").Node().
    AddProperty("pressure", "Pressure", homie.TypeFloat).SetUnit("hPa").Node().
    AddProperty("humidity", "Humidity", homie.TypeFloat).SetUnit("%")
```

Send the Homie attributes and or values to the MQTT client:

```go
// all homie attributes
for _, attribute := range device.GetHomieAttributes() {
    publish(attribute.Topic, attribute.Value)
}

// all property values
for _, attribute := range device.GetValues() {
    if attribute.Value != "" {
        publish(attribute.Topic, attribute.Value)
    }
}
```

Your `publish` function can use the MQTT client of your choice.

You can set a callback for sending property values for you when you set them:

```go
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

// either install a global callback on the device...
device.OnSet(onSet)

// ...or install the callback on this property only
device.Node("bmp280").Property("temperature").OnSet(onSet)

// new values will be published for you
device.Node("bme280").Property("temperature").Set(20.0)

```

## More information

See the [example](https://github.com/creativeprojects/go-homie/blob/main/example/main.go)

See [go doc](https://pkg.go.dev/github.com/creativeprojects/go-homie)
