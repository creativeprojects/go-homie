package homie

import (
	"fmt"
	"path"
)

type Setter func(topic, value string, dataType PropertyType)

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

func (p *Property) DataType() PropertyType {
	return p.dataType
}

func (p *Property) Set(value interface{}) *Property {
	p.value = fmt.Sprintf("%v", value)
	if p.setter != nil {
		p.setter(p.prefix, p.value, p.dataType)
	}
	return p
}

func (p *Property) Settable(settable bool) *Property {
	p.settable = settable
	return p
}

func (p *Property) SetUnit(unit string) *Property {
	p.unit = unit
	return p
}

func (p *Property) SetFormat(format string) *Property {
	p.format = format
	return p
}

func (p *Property) GetValue() TopicValuePair {
	return TopicValuePair{
		p.prefix,
		p.value,
	}
}

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
