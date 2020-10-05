package homie

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

const (
	homiePrefix = "homie"
)

type DeviceState string

// DeviceState
const (
	StateInit         DeviceState = "init"
	StateReady        DeviceState = "ready"
	StateDisconnected DeviceState = "disconnected"
	StateSleeping     DeviceState = "sleeping"
	StateLost         DeviceState = "lost"
	StateAlert        DeviceState = "alert"
)

type Device struct {
	prefix  string
	version string
	id      string
	name    string
	state   DeviceState
	nodes   map[string]*Node
}

// NewDevice creates a homie device.
//
// ID is used to create topics: homie/<ID>/...
// Name is the fullname of the device
//
// It will panic if ID cannot be used in a topic. You should check with IsValidID
func NewDevice(id, name string) *Device {
	if !IsValidID(id) {
		panic(fmt.Sprintf("invalid device ID: '%s'", id))
	}
	return &Device{
		prefix:  path.Join(homiePrefix, id),
		version: "4.0.0",
		id:      id,
		name:    name,
		state:   StateInit,
		nodes:   make(map[string]*Node, 0),
	}
}

func (d *Device) SetState(state DeviceState) *Device {
	d.state = state
	return d
}

// AddNode creates and add the node to the device
//
// ID is used to create topics: homie/device/<ID>/...
// Name is the fullname of the node
//
// It will panic if ID cannot be used in a topic. You should check with IsValidID
func (d *Device) AddNode(id, name, nodeType string) *Node {
	node := newNode(d, d.prefix, id, name, nodeType)
	d.nodes[id] = node
	return node
}

func (d *Device) Node(id string) *Node {
	return d.nodes[id]
}

func (d *Device) GetHomieAttributes() []TopicValuePair {
	attributes := make([]TopicValuePair, 0, len(d.nodes)*20)
	attributes = append(attributes, TopicValuePair{path.Join(d.prefix, attributeHomieVersion), d.version})
	attributes = append(attributes, TopicValuePair{path.Join(d.prefix, attributeName), d.name})
	attributes = append(attributes, TopicValuePair{path.Join(d.prefix, attributeState), string(d.state)})

	nodes := ""
	if d.nodes != nil && len(d.nodes) > 0 {
		keys := make([]string, len(d.nodes))

		i := 0
		for key := range d.nodes {
			keys[i] = key
			i++
		}
		// this is not strictly necessary but it helps with the unit tests
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		// add list of keys
		nodes = strings.Join(keys, ",")
	}
	attributes = append(attributes, TopicValuePair{path.Join(d.prefix, attributeNodes), nodes})

	// empty extensions for now
	attributes = append(attributes, TopicValuePair{path.Join(d.prefix, attributeExtensions), ""})

	// now get properties from children
	for _, node := range d.nodes {
		attributes = append(attributes, node.getAttributes()...)
	}
	return attributes
}

func (d *Device) GetState() TopicValuePair {
	return TopicValuePair{d.GetStateTopic(), string(d.state)}
}

func (d *Device) GetStateTopic() string {
	return path.Join(homiePrefix, d.id, attributeState)
}

func (d *Device) GetValues() []TopicValuePair {
	attributes := make([]TopicValuePair, 0, len(d.nodes)*3)
	for _, node := range d.nodes {
		attributes = append(attributes, node.getValues()...)
	}
	return attributes
}

func (d *Device) GetPropertySetters() map[string]*Property {
	properties := make(map[string]*Property, len(d.nodes)*3)
	for _, node := range d.nodes {
		props := node.getSetterProperties()
		if len(props) == 0 {
			continue
		}
		for topic, prop := range props {
			properties[topic] = prop
		}
	}
	return properties
}
