package schema

//go:generate stringer -type=ValueType valuetype.go

// ValueType is an enum of the type that can be represented by a schema.
type ValueType int

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeList
	TypeMap
	TypeSet
	TypeJSON // arbitrary JSON-serializable data structure
	typeObject
)

// NOTE: ValueType has more functions defined on it in schema.go. We can't
// put them here because we reference other files.
