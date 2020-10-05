package homie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyDeviceAttributes(t *testing.T) {
	device := NewDevice("deviceID", "deviceName")
	attributes := device.GetHomieAttributes()
	values := device.GetValues()

	assert.ElementsMatch(t, attributes, []TopicValuePair{
		{"homie/deviceID/$homie", "4.0.0"},
		{"homie/deviceID/$name", "deviceName"},
		{"homie/deviceID/$state", "init"},
		{"homie/deviceID/$nodes", ""},
		{"homie/deviceID/$extensions", ""},
	})
	assert.Empty(t, values)
}

func TestDeviceAttributes(t *testing.T) {
	device := NewDevice("deviceID", "deviceName")
	device.
		AddNode("node1", "node1 name", "test1").
		AddProperty("prop1", "prop1 name", TypeBoolean).Set(true).Node().Device().
		AddNode("node2", "node2 name", "test2").
		AddProperty("prop2", "prop2 name", TypeInteger).Set(11)

	attributes := device.GetHomieAttributes()
	values := device.GetValues()

	assert.ElementsMatch(t, attributes, []TopicValuePair{
		{"homie/deviceID/$homie", "4.0.0"},
		{"homie/deviceID/$name", "deviceName"},
		{"homie/deviceID/$state", "init"},
		{"homie/deviceID/$nodes", "node1,node2"},
		{"homie/deviceID/$extensions", ""},
		{"homie/deviceID/node1/$name", "node1 name"},
		{"homie/deviceID/node1/$type", "test1"},
		{"homie/deviceID/node1/$properties", "prop1"},
		{"homie/deviceID/node1/prop1/$name", "prop1 name"},
		{"homie/deviceID/node1/prop1/$datatype", "boolean"},
		{"homie/deviceID/node2/$name", "node2 name"},
		{"homie/deviceID/node2/$type", "test2"},
		{"homie/deviceID/node2/$properties", "prop2"},
		{"homie/deviceID/node2/prop2/$name", "prop2 name"},
		{"homie/deviceID/node2/prop2/$datatype", "integer"},
	})
	assert.ElementsMatch(t, values, []TopicValuePair{
		{"homie/deviceID/node1/prop1", "true"},
		{"homie/deviceID/node2/prop2", "11"},
	})
}

func TestEmptyPropertySetters(t *testing.T) {
	device := NewDevice("deviceID", "deviceName")
	device.
		AddNode("node1", "node1 name", "test1").
		AddProperty("prop1", "prop1 name", TypeBoolean).Set(true).Node().Device().
		AddNode("node2", "node2 name", "test2").
		AddProperty("prop2", "prop2 name", TypeInteger).Set(11)

	props := device.GetPropertySetters()
	assert.Empty(t, props)
}

func TestPropertySetters(t *testing.T) {
	device := NewDevice("deviceID", "deviceName")
	device.
		AddNode("node1", "node1 name", "test1").
		AddProperty("prop1", "prop1 name", TypeBoolean).Settable(true).Node().Device().
		AddNode("node2", "node2 name", "test2").
		AddProperty("prop2", "prop2 name", TypeInteger).Settable(true).Node().
		AddProperty("prop3", "prop3 name", TypeInteger).Settable(true).Node()

	props := device.GetPropertySetters()
	assert.Len(t, props, 3)
	assert.NotEmpty(t, props["homie/deviceID/node1/prop1/set"])
	assert.NotEmpty(t, props["homie/deviceID/node2/prop2/set"])
	assert.NotEmpty(t, props["homie/deviceID/node2/prop3/set"])
}
