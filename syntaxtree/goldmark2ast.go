package syntaxtree

import (
	"container/list"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

// | Field           | Type                                                             | Description                                                               |
// | --------------- | ---------------------------------------------------------------- | ------------------------------------------------------------------------- |
// | id              | ?snowflake                                                       | [emoji id](#DOCS_REFERENCE/image-formatting)                              |
// | name            | ?string (can be null only in reaction emoji objects)             | emoji name                                                                |
// | roles?          | array of [role](#DOCS_TOPICS_PERMISSIONS/role-object) object ids | roles this emoji is whitelisted to                                        |
// | user?           | [user](#DOCS_RESOURCES_USER/user-object) object                  | user that created this emoji                                              |
// | require_colons? | boolean                                                          | whether this emoji must be wrapped in colons                              |
// | managed?        | boolean                                                          | whether this emoji is managed                                             |
// | animated?       | boolean                                                          | whether this emoji is animated                                            |
// | available?      | boolean                                                          | whether this emoji can be used, may be false due to loss of Server Boosts |

// md := goldmark.New(
// 	goldmark.WithExtensions(extension.Table),
// )
// nodes := md.Parser().Parse(text.NewReader(fileData))
// // dfsPrintChildren(nodes, 0)
// tables := getTables(nodes)

// _ = syntaxtree.Translate(tables, fileData)

type GoldMarkTranslator struct{}

// Translate acts as our entry point when performing model adaption between
// a Discord markdown API model and our GoldMark Abstract Syntax Tree
// representation. It requires the Table extension as all the models within
// the discord markdown API are within the tables!
func (*GoldMarkTranslator) Translate(data []byte) map[string]*Object {
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
	)
	nodes := md.Parser().Parse(text.NewReader(data))
	// dfsPrintChildren(nodes, 0)
	tables := getTables(nodes, data)

	return translate(tables, data)
}

func getTables(node ast.Node, data []byte) []markdownObject {
	rv := make([]markdownObject, 0, 16)
	if node.HasChildren() {
		for i := node.FirstChild(); i != nil; i = i.NextSibling() {
			if i.Kind().String() == "Table" && ensureObject(i.PreviousSibling(), data) {
				var temp markdownObject
				temp.Header = i.PreviousSibling()
				temp.Table = i
				rv = append(rv, temp)
			}
		}
	}

	return rv
}

func ensureObject(node ast.Node, data []byte) bool {
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		line = line.TrimLeftSpace(data)
		line = line.TrimRightSpace(data)
		s := string(line.Value(data))
		// log.Println(s)
		return strings.Contains(strings.ToLower(s), "structure")
	}

	return false
}

func translate(nodes []markdownObject, data []byte) map[string]*Object {
	rv := make(map[string]*Object)
	for _, item := range nodes {
		obj := new(Object)
		obj.Name = extractObjectName(item.Header, data)
		obj.Fields = extractFields(item.Table, data)
		rv[obj.Name] = obj
	}
	return rv
}

func extractObjectName(item ast.Node, data []byte) string {
	lines := item.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		line = line.TrimLeftSpace(data)
		line = line.TrimRightSpace(data)
		s := string(line.Value(data))
		return s
	}
	return ""
}

func extractFields(item ast.Node, data []byte) []Field {
	dfsPrintChildren(item, 0)
	rv := make([]Field, item.ChildCount()-1, item.ChildCount()-1)
	// Skip the header
	currentChild := item.FirstChild().NextSibling()
	for i := 0; i < item.ChildCount()-1; i++ {
		cell := currentChild.FirstChild()
		name := string(cell.Text(data))
		rv[i].Optional = strings.HasSuffix(name, "?")
		rv[i].Name = strings.TrimSuffix(string(cell.Text(data)), "?")

		cell = cell.NextSibling()
		cellText := string(cell.Text(data))
		rv[i].Nullable = strings.HasPrefix(cellText, "?")
		rv[i].T = extractType(cell, data)
		currentChild = currentChild.NextSibling()
	}

	return rv
}

func extractType(cell ast.Node, data []byte) Type {
	var stringTypeMapping = map[string]Type{
		"array":     Array,
		"snowflake": Snowflake,
		"boolean":   Bool,
		"string":    String,
		"integer":   Int64,
		"object":    Reference,
	}

	cellText := string(cell.Text(data))
	rv := Unknown

	possibleTypes := list.New()

	for k, v := range stringTypeMapping {
		if strings.Contains(cellText, k) {
			possibleTypes.PushBack(v)
		}
	}

	if possibleTypes.Len() == 1 {
		rv = possibleTypes.Front().Value.(Type)
	} else {
		if possibleTypes.Front().Value == Array {
			switch possibleTypes.Front().Next().Value {
			case Bool:
				rv = BoolArray
			case String:
				rv = StringArray
			case Int64:
				rv = Int64Array
			case UInt64:
				rv = UInt64Array
			case Int32:
				rv = Int32Array
			case UInt32:
				rv = UInt32Array
			case Float:
				rv = FloatArray
			case Double:
				rv = DoubleArray
			case Snowflake:
				rv = SnowflakeArray
			case Reference:
				rv = ReferenceArray
			}
		}
	}

	return rv
}

func dfsPrintChildren(node ast.Node, ident int) {
	// log.Printf("%*c%s\n", ident, ' ', node.Kind().String())
	if node.HasChildren() {
		for i := node.FirstChild(); i != nil; i = i.NextSibling() {
			dfsPrintChildren(i, ident+4)
		}
	}
}
