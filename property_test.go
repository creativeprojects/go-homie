package homie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPropertyAttributes(t *testing.T) {
	prop := newProperty(nil, "test", "id", "name", TypeString)
	attributes := prop.getAttributes()
	assert.ElementsMatch(t, attributes, []TopicValuePair{
		{"test/id/$name", "name"},
		{"test/id/$datatype", "string"},
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
