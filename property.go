package homie

import (
	"fmt"
	"path"
)

// PropertyType is the type of the property
type PropertyType string

// PropertyType
const (
	TypeInteger PropertyType = "integer"
	TypeFloat   PropertyType = "float"
	TypeBoolean PropertyType = "boolean"
	TypeString  PropertyType = "string"
	TypeEnum    PropertyType = "enum"
	TypeColor   PropertyType = "color"
)

// Property definition
type Property struct {
	node     *Node
	prefix   string
	id       string
	name     string
	dataType PropertyType
	format   string
	unit     string
	settable bool
	retained bool
	value    string
	setter   Setter
}

func newProperty(node *Node, prefix, id, name string, dataType PropertyType) *Property {
	if !IsValidID(id) {
		panic(fmt.Sprintf("invalid property ID: '%s'", id))
	}
	return &Property{
		node:     node,
		prefix:   path.Join(prefix, id),
		id:       id,
		name:     name,
		dataType: dataType,
		retained: true,
	}
}

// Node returns the node this property is attached to.
// This is handy for chaining declaration
func (p *Property) Node() *Node {
	return p.node
}

// DataType returns the current type of the property
func (p *Property) DataType() PropertyType {
	return p.dataType
}

// Set a new property value
func (p *Property) Set(value interface{}) *Property {
	p.value = fmt.Sprintf("%v", value)
	if p.setter != nil {
		p.setter(p.prefix, p.value, p.dataType)
	} else if p.node != nil && p.node.device != nil && p.node.device.setter != nil { // protection for unit tests creating lose properties
		p.node.device.setter(p.prefix, p.value, p.dataType)
	}
	return p
}

// Settable tells the property if it can be set via a Homie set command.
// for more information, https://homieiot.github.io/specification/#property-command-topic
func (p *Property) Settable(settable bool) *Property {
	p.settable = settable
	return p
}

// SetUnit defines a unit on the property
func (p *Property) SetUnit(unit string) *Property {
	p.unit = unit
	return p
}

// SetFormat defines a property format.
// for more information on property format, see https://homieiot.github.io/specification/#properties
func (p *Property) SetFormat(format string) *Property {
	p.format = format
	return p
}

// SetRetained changes the retained flag as described:
// https://homieiot.github.io/specification/#property-attributes
func (p *Property) SetRetained(retained bool) *Property {
	p.retained = retained
	return p
}

// GetValue returns the Topic/Value pair of the property
func (p *Property) GetValue() TopicValuePair {
	return TopicValuePair{
		p.prefix,
		p.value,
	}
}

// OnSet defines a callback for when you Set() a new property value
// if there's already a OnSet callback defined on the device, this callback will be used instead.
// You pass nil to the method to remove the callback
func (p *Property) OnSet(setter Setter) *Property {
	p.setter = setter
	return p
}

func (p *Property) getSetterTopic() string {
	if !p.settable {
		return ""
	}
	return path.Join(p.prefix, "set")
}

func (p *Property) getAttributes() []TopicValuePair {
	attributes := make([]TopicValuePair, 0, 6)
	attributes = append(attributes, TopicValuePair{path.Join(p.prefix, attributeName), p.name})
	attributes = append(attributes, TopicValuePair{path.Join(p.prefix, attributeDatatype), string(p.dataType)})
	if p.format != "" {
		attributes = append(attributes, TopicValuePair{path.Join(p.prefix, attributeFormat), p.format})
	}
	if p.unit != "" {
		attributes = append(attributes, TopicValuePair{path.Join(p.prefix, attributeUnit), p.unit})
	}
	if p.settable {
		attributes = append(attributes, TopicValuePair{path.Join(p.prefix, attributeSettable), fmt.Sprintf("%v", p.settable)})
	}
	if !p.retained {
		attributes = append(attributes, TopicValuePair{path.Join(p.prefix, attributeRetained), fmt.Sprintf("%v", p.retained)})
	}
	return attributes
}
