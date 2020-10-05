package homie

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

// Configuration variables
var (
	DefaultVersion = "4.0.0"
	DefaultRoot    = "homie"
)

// DeviceState represents the current state of the device
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

// Device is the definition of your Homie device
type Device struct {
	prefix  string
	version string
	id      string
	name    string
	state   DeviceState
	setter  Setter
	nodes   map[string]*Node
}

// NewDevice creates a homie device.
//
// ID is used to create topics: homie/<ID>/...
// Name is the fullname of the device
//
// It will panic if ID cannot be used in a topic. You can check with IsValidID before calling NewDevice.
//
// see documentation: https://homieiot.github.io/specification/#topic-ids
func NewDevice(id, name string) *Device {
	if !IsValidID(id) {
		panic(fmt.Sprintf("invalid device ID: '%s'", id))
	}
	return &Device{
		prefix:  path.Join(DefaultRoot, id),
		version: DefaultVersion,
		id:      id,
		name:    name,
		state:   StateInit,
		nodes:   make(map[string]*Node, 0),
	}
}

// SetRoot changes the MQTT root topic: the default root is "homie".
//
// for more information: https://homieiot.github.io/specification/#base-topic
func (d *Device) SetRoot(root string) *Device {
	d.prefix = path.Join(root, d.id)
	return d
}

// SetState sets the state of the device
//
// for more information about the device states: https://homieiot.github.io/specification/#device-lifecycle
func (d *Device) SetState(state DeviceState) *Device {
	d.state = state
	if d.setter != nil {
		d.setter(d.GetStateTopic(), string(state), TypeString)
	}
	return d
}

// AddNode creates and add the node to the device
//
// ID is used to create topics: homie/device/<ID>/...
// Name is the fullname of the node
//
// It will panic if ID cannot be used in a topic. You can check with IsValidID before calling the method.
func (d *Device) AddNode(id, name, nodeType string) *Node {
	node := newNode(d, d.prefix, id, name, nodeType)
	d.nodes[id] = node
	return node
}

// Node returns the node of that id
func (d *Device) Node(id string) *Node {
	return d.nodes[id]
}

// GetHomieAttributes returns all attributes as a Topic/Value pair
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

// GetState returns the device state as a Topic/Value pair
func (d *Device) GetState() TopicValuePair {
	return TopicValuePair{d.GetStateTopic(), string(d.state)}
}

// GetStateTopic returns the topic of the device state
func (d *Device) GetStateTopic() string {
	return path.Join(d.prefix, attributeState)
}

// GetValues return the values of all properties
func (d *Device) GetValues() []TopicValuePair {
	attributes := make([]TopicValuePair, 0, len(d.nodes)*3)
	for _, node := range d.nodes {
		attributes = append(attributes, node.getValues()...)
	}
	return attributes
}

// GetPropertySetters returns the topics of all the property setters.
// a property will have a setter if it was defined with Settable(true)
//
// a Homie property setter is a topic like: homie/<deviceID>/<nodeID>/<propertyID>/set
//
// see documentation: https://homieiot.github.io/specification/#property-command-topic
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

// OnSet adds a global callback when a property value is changed (via the Set method)
func (d *Device) OnSet(setter Setter) *Device {
	d.setter = setter
	return d
}
