package syntaxtree

import (
	"github.com/yuin/goldmark/ast"
)

// Type is a Pog
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
	Name     string
	T        Type
	Optional bool
	Nullable bool
}

type markdownObject struct {
	Header ast.Node
	Table  ast.Node
}
