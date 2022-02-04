package langadaptor

import (
	"bytes"
	"log"
	"strings"
	"text/template"

	"github.com/Reprezent/discord-api-models/syntaxtree"
)

type JavaAdaptor struct{}

func (*JavaAdaptor) Generate(data map[string]*syntaxtree.Object) ([]string, []string, error) {
	var filename_rv []string
	var data_rv []string
	for _, model := range data {
		s := normalizeFileName(model.Name)
		str, err := generateSourceFile(s, model)
		if err != nil {
			return nil, nil, err
		}
		data_rv = append(data_rv, str)
		filename_rv = append(filename_rv, s+".java")
	}
	return filename_rv, data_rv, nil
}

func isReference(f syntaxtree.Field) bool {
	return f.T == syntaxtree.Reference
}

func typeSwitch(f syntaxtree.Field) string {

	switch f.T {
	case syntaxtree.Unknown:
		return "?"
	case syntaxtree.Array:
		return "[]"
	case syntaxtree.Bool:
		return "Boolean"
	case syntaxtree.String:
		return "String"
	case syntaxtree.Int64:
		fallthrough
	case syntaxtree.UInt64:
		return "Long"
	case syntaxtree.Int32:
		fallthrough
	case syntaxtree.UInt32:
		return "Int"
	case syntaxtree.Float:
		return "Float"
	case syntaxtree.Double:
		return "Double"
	case syntaxtree.Snowflake:
		return "String"
	case syntaxtree.Reference:
		if f.Reference != nil {
			return normalizeFileName(f.Reference.Name)
		} else {
			log.Printf("Reference Missed: %+v", f)
			return "?"
		}
	case syntaxtree.BoolArray:
		return "boolean[]"
	case syntaxtree.StringArray:
		return "String[]"
	case syntaxtree.Int64Array:
		fallthrough
	case syntaxtree.UInt64Array:
		return "long[]"
	case syntaxtree.Int32Array:
		fallthrough
	case syntaxtree.UInt32Array:
		return "int[]"
	case syntaxtree.FloatArray:
		return "float[]"
	case syntaxtree.DoubleArray:
		return "double[]"
	case syntaxtree.SnowflakeArray:
		return "String[]"
	case syntaxtree.ReferenceArray:
		if f.Reference != nil {
			return normalizeFileName(f.Reference.Name) + "[]"
		} else {
			return "?[]"
		}
	case syntaxtree.Binary:
		return "byte[]"
	default:
		return "?"
	}

}

func normalizeFileName(s string) string {
	var rv string = ""
	for _, tok := range strings.Split(s, " ") {
		if tok == "Structure" || tok == "Object" {
			continue
		}
		rv += strings.Title(strings.ToLower(tok))
	}

	return rv
}

func normalzieVariableNames(s string) string {
	var rv string = ""
	for _, tok := range strings.Split(s, "_") {
		rv += strings.Title(strings.ToLower(tok))
	}

	return rv
}

const templateString = `
package java.source;

import java.util.Optional;

public class {{normalize .Name}}
{
{{range .Fields}}{{$name := n .Name}}    public {{if .Optional}}Optional<{{end}}{{t .}}{{if .Optional}}>{{end}} {{$name}};
{{end}}
{{range .Fields}}{{$name := n .Name}}    public void Set{{$name}}({{t .}} {{$name}}) { this.{{$name}} = {{$name}}; }
{{end}}
{{range .Fields}}{{$name := n .Name}}    public {{t .}} Get{{$name}}({{t .}} {{$name}}) { return this.{{$name}}; }
{{end}}


    public {{normalize .Name}}(JSONObject data)
    {
{{range .Fields}}        {{$name := n .Name}}parse{{$name}}(data);
{{end}}
    }
{{range .Fields}}{{$name := n .Name}}
    private void parse{{$name}}(JSONObject data)
    {
        final String key = "{{.Name}}";
        if(!data.containsKey(key))
        {
            missingID(key);
        }
        
        this.{{$name}} = data.get{{caps (t .)}}(key).get{{if eq (t .) "String"}}String{{else if (r .)}}Object{{else}}Val{{end}}();
    }
{{end}}
}
`

func generateSourceFile(s string, object *syntaxtree.Object) (string, error) {
	tmpl, err := template.New("JavaSource").Funcs(template.FuncMap{
		"upper":     strings.ToUpper,
		"caps":      func(a string) string { return strings.Title(strings.ToLower(a)) },
		"n":         normalzieVariableNames,
		"normalize": normalizeFileName,
		"t":         typeSwitch,
		"r":         isReference,
	}).Parse(templateString)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, object)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
