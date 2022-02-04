package syntaxtree

import (
	"fmt"
	"log"
	"sort"
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
			if i.Kind().String() == "Table" {
				for current := i.PreviousSibling(); current != nil; current = current.PreviousSibling() {
					if current.Kind().String() == "Heading" {
						if ensureLinkableType(current, data) {
							var temp markdownObject
							temp.Table = i
							temp.Header = current
							rv = append(rv, temp)
						} else {
							log.Printf("Found table without header for %v", extractName(current, data))
						}
						break
					}
				}
			}
		}
	}

	return rv
}

func extractName(node ast.Node, data []byte) string {
	lines := node.Lines()
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		line = line.TrimLeftSpace(data)
		line = line.TrimRightSpace(data)
		s := string(line.Value(data))
		return s
	}

	return ""
}
func ensureLinkableType(node ast.Node, data []byte) bool {
	if node.Kind().String() != "Heading" {
		return false
	}
	s := strings.ToLower(extractName(node, data))
	return strings.Contains(s, "structure") ||
		strings.Contains(s, "struct") ||
		strings.Contains(s, "object")
}

func translate(nodes []markdownObject, data []byte) map[string]*Object {
	rv := make(map[string]*Object)
	for _, item := range nodes {
		obj := new(Object)
		// log.Printf("%s is at %p", obj.Name, obj)
		obj.Name = extractName(item.Header, data)
		obj.Fields = extractFields(item.Table, data)
		rv[cleanReferenceName(obj.Name)] = obj
		// log.Printf("%s is at %p", obj.Name, obj)
	}

	var builder strings.Builder
	builder.WriteString("References:\n")
	for k, v := range rv {
		builder.WriteString("    - ")
		builder.WriteString(k)
		builder.WriteString(" - ")
		builder.WriteString(fmt.Sprintf("%p", v))
		builder.WriteString("\n")
	}
	// keys := make([]string, 0, len(rv))
	// for k := range rv {
	// 	keys = append(keys, k)
	// }

	log.Printf("References %s\n", builder.String())
	for _, v := range rv {
		for d, i := range v.Fields {
			if i.T == Reference || i.T == ReferenceArray {
				temp, ok := rv[cleanReferenceName(i.ReferenceName)]
				if ok {
					i.Reference = temp
					// log.Printf("Pointer Reference '%v' found for Reference Name '%v' at %p-%p for structure '%v' at %p", i.ReferenceName, cleanReferenceName(i.ReferenceName), i.Reference, v.Fields[d].Reference, v.Name, v)
					v.Fields[d] = i
					// log.Printf("Pointer Reference AFTER '%v' found for Reference Name '%v' at %p-%p for structure '%v' at %p", i.ReferenceName, cleanReferenceName(i.ReferenceName), i.Reference, v.Fields[d].Reference, v.Name, v)
				} else {
					log.Printf("Object '%s' Could not find reference for %v -- %v\n", v.Name, i.ReferenceName, cleanReferenceName(i.ReferenceName))
				}
			}
		}
	}
	return rv
}

func removeTextComments(text string) string {
	index := strings.Index(text, ",")
	if index > 0 {
		text = text[:index]
	}

	index = strings.Index(text, "(")
	if index > 0 {
		text = text[:index]
	}

	return text
}
func cleanReferenceName(reference string) string {
	strings_to_remove := [...]string{"Partial", "Array Of", "?"}
	for _, v := range strings_to_remove {
		reference = strings.ReplaceAll(reference, v, "")
	}
	reference = strings.ReplaceAll(reference, "A ", "")
	reference = strings.ReplaceAll(reference, "Objects", "Object")
	reference = strings.ReplaceAll(reference, "Structure", "Object")
	reference = strings.ReplaceAll(reference, "Struct", "Object")
	reference = strings.ReplaceAll(reference, "Struct", "Object")
	reference = removeTextComments(reference)

	return strings.Trim(reference, " ")
}

func extractFields(item ast.Node, data []byte) []Field {
	// dfsPrintChildren(item, 0)
	rv := make([]Field, item.ChildCount()-1, item.ChildCount()-1)
	// Skip the header
	currentChild := item.FirstChild().NextSibling()
	for i := 0; i < item.ChildCount()-1; i++ {
		fieldCell := currentChild.FirstChild()
		fieldName := strings.ReplaceAll(string(fieldCell.Text(data)), "\\*", "")
		rv[i].Optional = strings.HasSuffix(fieldName, "?")
		rv[i].Name = strings.TrimSuffix(string(fieldName), "?")

		typeCell := fieldCell.NextSibling()
		typeCellText := string(typeCell.Text(data))
		rv[i].Nullable = strings.HasPrefix(typeCellText, "?")
		rv[i].T = extractType(typeCell, data)
		if rv[i].T == Reference || rv[i].T == ReferenceArray {
			rv[i].ReferenceName = normalizeObjectReference(typeCellText)
			log.Printf("Reference Name '%v' found for Name '%v'", rv[i].ReferenceName, rv[i].Name)
		}
		currentChild = currentChild.NextSibling()
	}

	return rv
}

func normalizeObjectReference(s string) string {
	var rv string = ""
	for _, tok := range strings.Split(s, " ") {
		rv += strings.Title(strings.ToLower(tok)) + " "
	}

	rv = strings.Trim(rv, " ")

	return rv
}

func extractType(cell ast.Node, data []byte) Type {
	var stringTypeMapping = map[string]Type{
		"array":       Array,
		"snowflake":   Snowflake,
		"id":          Snowflake,
		"ids":         Snowflake,
		"object ids":  Snowflake,
		"bool ":       Bool,
		"boolean":     Bool,
		"string":      String,
		"int":         Int64,
		"integer":     Int64,
		"float":       Float,
		"object":      Reference,
		"timestamp":   Timestamp,
		"binary":      Binary,
		"single byte": Byte,
	}

	cellText := string(cell.Text(data))
	rv := Unknown

	possibleTypes := make(map[Type]bool)

	for k, v := range stringTypeMapping {
		if strings.Contains(removeTextComments(strings.ToLower(cellText)), k) {
			possibleTypes[v] = true
		}
	}

	keys := make([]Type, 0, len(possibleTypes))
	for k := range possibleTypes {
		keys = append(keys, k)
	}

	sort.Sort(&typeSorter{
		types: keys,
	})

	if len(keys) == 1 {
		for _, e := range keys {
			rv = e
		}
	} else if len(keys) > 1 {
		for _, e := range keys {
			if e == Array {
				continue
			}
			switch e {
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
			break
		}

		if rv == Unknown {
			log.Printf("Unknown type for %v. Have %#v\n", cellText, keys)
		}
	} else {
		log.Printf("Unknown type for %v\n", cellText)
	}

	return rv
}

func dfsPrintChildren(node ast.Node, ident int) {
	log.Printf("%*c%s\n", ident, ' ', node.Kind().String())
	if node.HasChildren() {
		for i := node.FirstChild(); i != nil; i = i.NextSibling() {
			dfsPrintChildren(i, ident+4)
		}
	}
}
