package homie

const (
	attributeHomieVersion = "$homie"
	attributeName         = "$name"
	attributeState        = "$state"
	attributeNodes        = "$nodes"
	attributeExtensions   = "$extensions"
	attributeProperties   = "$properties"
	attributeType         = "$type"
	attributeDatatype     = "$datatype"
	attributeFormat       = "$format"
	attributeUnit         = "$unit"
	attributeSettable     = "$settable"
	attributeRetained     = "$retained"
)

// TopicValuePair represents a MQTT topic, value and data type
type TopicValuePair struct {
	Topic string
	Value string
	// DataType PropertyType
}
