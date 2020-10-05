package homie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEmptyNodeAttributes(t *testing.T) {
	node := newNode(nil, "test", "nodeID", "nodeName", "nodeType")
	attributes := node.getAttributes()
	values := node.getValues()
	assert.ElementsMatch(t, attributes, []TopicValuePair{
		{"test/nodeID/$name", "nodeName"},
		{"test/nodeID/$type", "nodeType"},
		{"test/nodeID/$properties", ""},
	})
	assert.Empty(t, values)
}

func TestGetNodeAttributes(t *testing.T) {
	node := newNode(nil, "test", "nodeID", "nodeName", "nodeType")
	node.AddProperty("prop1", "prop1", TypeInteger).Set(10)
	node.AddProperty("prop2", "prop2", TypeInteger).Set(20)

	attributes := node.getAttributes()
	values := node.getValues()
	assert.ElementsMatch(t, attributes, []TopicValuePair{
		{"test/nodeID/$name", "nodeName"},
		{"test/nodeID/$type", "nodeType"},
		{"test/nodeID/$properties", "prop1,prop2"},
		{"test/nodeID/prop1/$name", "prop1"},
		{"test/nodeID/prop1/$datatype", "integer"},
		{"test/nodeID/prop2/$name", "prop2"},
		{"test/nodeID/prop2/$datatype", "integer"},
	})
	assert.ElementsMatch(t, values, []TopicValuePair{
		{"test/nodeID/prop1", "10"},
		{"test/nodeID/prop2", "20"},
	})
}

func TestGetPropertySetters(t *testing.T) {
	node := newNode(nil, "test", "nodeID", "nodeName", "nodeType")
	node.AddProperty("prop1", "prop1", "property1")
	node.AddProperty("prop2", "prop2", "property2").Settable(true)

	setters := node.getSetterProperties()
	assert.NotEmpty(t, setters)
	assert.Len(t, setters, 1)
	assert.NotEmpty(t, setters["test/nodeID/prop2/set"])
}

func TestUndefinedProperty(t *testing.T) {
	node := newNode(nil, "test", "nodeID", "nodeName", "nodeType")
	property := node.Property("propertyID")
	assert.Nil(t, property)
}
