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

func (*GoldMarkTranslator) Translate(data []byte) []Object {
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
	)
	nodes := md.Parser().Parse(text.NewReader(data))
	// dfsPrintChildren(nodes, 0)
	tables := getTables(nodes, data)

	return translate(tables, data)
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

func translate(nodes []markdownObject, data []byte) []Object {
	rv := make([]Object, len(nodes), len(nodes))
	for i, item := range nodes {
		rv[i].Name = extractObjectName(item.Header, data)
		rv[i].Fields = extractFields(item.Table, data)
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
		"snowflake": Snowflake,
		"array":     Array,
		"boolean":   Bool,
		"string":    String,
		"integer":   Int64,
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
