package syntaxtree

import (
	"github.com/yuin/goldmark/ast"
)

// Type is an int
type Type int

const (
	Unknown Type = iota
	Array
	Bool
	String
	Int64
	UInt64
	Int32
	UInt32
	Float
	Double
	Snowflake
	Reference
	Timestamp
	Date
	Binary
	Byte
	BoolArray
	StringArray
	Int64Array
	UInt64Array
	Int32Array
	UInt32Array
	FloatArray
	DoubleArray
	SnowflakeArray
	ReferenceArray
)

// Object is a Pog
type Object struct {
	Name   string
	Fields []Field
}

// Field is a Pog
type Field struct {
	Name          string
	T             Type
	Optional      bool
	Nullable      bool
	Reference     *Object
	ReferenceName string
}

type markdownObject struct {
	Header ast.Node
	Table  ast.Node
}

type typeSorter struct {
	types []Type
}

func (s *typeSorter) Len() int {
	return len(s.types)
}

// Swap is part of sort.Interface.
func (s *typeSorter) Swap(i, j int) {
	s.types[i], s.types[j] = s.types[j], s.types[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *typeSorter) Less(i, j int) bool {
	return s.types[i] < s.types[j]
}
