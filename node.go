package homie

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

type Node struct {
	device     *Device
	prefix     string
	id         string
	name       string
	nodeType   string
	properties map[string]*Property
}

func newNode(device *Device, prefix, id, name, nodeType string) *Node {
	if !IsValidID(id) {
		panic(fmt.Sprintf("invalid node ID: '%s'", id))
	}
	return &Node{
		device:     device,
		prefix:     path.Join(prefix, id),
		id:         id,
		name:       name,
		nodeType:   nodeType,
		properties: make(map[string]*Property, 1),
	}
}

// AddProperty creates and add the property to the node
//
// ID is used to create topics: homie/device/node/<ID>/...
// Name is the fullname of the property
//
// It will panic if ID cannot be used in a topic. You should check with IsValidID
func (n *Node) AddProperty(id, name string, propertyType PropertyType) *Property {
	prop := newProperty(n, n.prefix, id, name, propertyType)
	n.properties[id] = prop
	return prop
}

// Device returns the device the node is attached to.
// This can be handy for chaining declaration
func (n *Node) Device() *Device {
	return n.device
}

func (n *Node) Property(id string) *Property {
	return n.properties[id]
}

func (n *Node) getAttributes() []TopicValuePair {
	attributes := make([]TopicValuePair, 0, 6*len(n.properties))
	attributes = append(attributes, TopicValuePair{path.Join(n.prefix, attributeName), n.name})
	attributes = append(attributes, TopicValuePair{path.Join(n.prefix, attributeType), n.nodeType})

	properties := ""
	if n.properties != nil && len(n.properties) > 0 {
		keys := make([]string, len(n.properties))

		i := 0
		for key := range n.properties {
			keys[i] = key
			i++
		}
		// this is not strictly necessary but it helps with the unit tests
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		// add list of keys
		properties = strings.Join(keys, ",")
	}
	attributes = append(attributes, TopicValuePair{path.Join(n.prefix, attributeProperties), properties})

	for _, prop := range n.properties {
		attributes = append(attributes, prop.getAttributes()...)
	}
	return attributes
}

func (n *Node) getValues() []TopicValuePair {
	attributes := make([]TopicValuePair, 0, len(n.properties))
	if len(n.properties) == 0 {
		return attributes
	}
	for _, prop := range n.properties {
		attributes = append(attributes, prop.GetValue())
	}
	return attributes
}

func (n *Node) getSetterProperties() map[string]*Property {
	properties := make(map[string]*Property, len(n.properties))
	for _, prop := range n.properties {
		topic := prop.getSetterTopic()
		if topic == "" {
			continue
		}
		properties[topic] = prop
	}
	return properties
}
