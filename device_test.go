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

func TestChangeRootOfEmptyDeviceAttributes(t *testing.T) {
	device := NewDevice("deviceID", "deviceName").SetRoot("unit/test")
	attributes := device.GetHomieAttributes()
	values := device.GetValues()

	assert.ElementsMatch(t, attributes, []TopicValuePair{
		{"unit/test/deviceID/$homie", "4.0.0"},
		{"unit/test/deviceID/$name", "deviceName"},
		{"unit/test/deviceID/$state", "init"},
		{"unit/test/deviceID/$nodes", ""},
		{"unit/test/deviceID/$extensions", ""},
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

func TestSetCallbackOnChangeState(t *testing.T) {
	call := false
	onSet := func(topic, value string, dataType PropertyType) {
		call = true
		assert.Equal(t, "homie/deviceID/$state", topic)
		assert.Equal(t, "sleeping", value)
	}
	device := NewDevice("deviceID", "deviceName").OnSet(onSet)
	device.SetState(StateSleeping)
	assert.True(t, call)
}

func TestUndefinedNode(t *testing.T) {
	device := NewDevice("deviceID", "deviceName")
	node := device.Node("nodeID")
	assert.Nil(t, node)
}
