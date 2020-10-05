package homie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSimplePropertyAttributes(t *testing.T) {
	prop := newProperty(nil, "test", "id", "name", TypeString)
	attributes := prop.getAttributes()
	assert.ElementsMatch(t, attributes, []TopicValuePair{
		{"test/id/$name", "name"},
		{"test/id/$datatype", "string"},
	})
}

func TestGetFullPropertyAttributes(t *testing.T) {
	prop := newProperty(nil, "test", "id", "name", TypeString).SetFormat("format").SetUnit("unit").Settable(true).SetRetained(false)
	attributes := prop.getAttributes()
	assert.ElementsMatch(t, attributes, []TopicValuePair{
		{"test/id/$name", "name"},
		{"test/id/$datatype", "string"},
		{"test/id/$format", "format"},
		{"test/id/$unit", "unit"},
		{"test/id/$settable", "true"},
		{"test/id/$retained", "false"},
	})
}

func TestBooleanValue(t *testing.T) {
	prop := newProperty(nil, "test", "id", "name", TypeBoolean)
	prop.Set(false)
	assert.Equal(t, TopicValuePair{"test/id", "false"}, prop.GetValue())
	prop.Set(true)
	assert.Equal(t, TopicValuePair{"test/id", "true"}, prop.GetValue())
}

func TestGetEmptySetterTopic(t *testing.T) {
	prop := newProperty(nil, "test", "id", "name", TypeBoolean)
	assert.Equal(t, "", prop.getSetterTopic())
}

func TestGetSetterTopic(t *testing.T) {
	prop := newProperty(nil, "test", "id", "name", TypeBoolean).Settable(true)
	assert.Equal(t, "test/id/set", prop.getSetterTopic())
}

func TestDeviceCallback(t *testing.T) {
	call := false
	onSet := func(topic, value string, dataType PropertyType) {
		call = true
		assert.Equal(t, "homie/deviceID/nodeID/test", topic)
		assert.Equal(t, "true", value)
	}
	device := NewDevice("deviceID", "device name").OnSet(onSet)
	property := device.AddNode("nodeID", "node name", "node type").AddProperty("test", "Test", TypeBoolean)
	property.Set(true)
	assert.True(t, call)
}

func TestCallbackOverride(t *testing.T) {
	call := 0
	onSetDevice := func(topic, value string, dataType PropertyType) {
		call += 10
	}
	onSetOverride := func(topic, value string, dataType PropertyType) {
		call++
		assert.Equal(t, "homie/deviceID/nodeID/test", topic)
		assert.Equal(t, "true", value)
	}
	device := NewDevice("deviceID", "device name").OnSet(onSetDevice)
	property := device.AddNode("nodeID", "node name", "node type").AddProperty("test", "Test", TypeBoolean).OnSet(onSetOverride)
	property.Set(true)
	assert.Equal(t, 1, call)
}
