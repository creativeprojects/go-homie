package homie

// Setter is the signature of the callback to send data to a MQTT client
type Setter func(topic, value string, dataType PropertyType)
