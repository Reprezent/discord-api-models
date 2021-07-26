package langadaptor

import (
	"log"
	"os"
	"strings"
	"text/template"

	"../syntaxtree"
)

type JavaAdaptor struct{}

func (*JavaAdaptor) Generate(data []syntaxtree.Object) ([]string, error) {
	for _, model := range data {
		s := normalizeFileName(model.Name)
		err := generateSourceFile(s, model)
		if err != nil {
			log.Println("We ded")
		}

	}
	return nil, nil
}

func typeSwitch(t syntaxtree.Type) string {
	switch t {
	case syntaxtree.Unknown:
		return "?"
	case syntaxtree.Array:
		return "[]"
	case syntaxtree.Bool:
		return "boolean"
	case syntaxtree.String:
		return "String"
	case syntaxtree.Int64:
		fallthrough
	case syntaxtree.UInt64:
		return "long"
	case syntaxtree.Int32:
		fallthrough
	case syntaxtree.UInt32:
		return "int"
	case syntaxtree.Float:
		return "float"
	case syntaxtree.Double:
		return "double"
	case syntaxtree.Snowflake:
		return "String"
	case syntaxtree.Reference:
		return "?"
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
		return "?"
	default:
		return "?"
	}

}

func normalizeFileName(s string) string {
	var rv string = ""
	for _, tok := range strings.Split(s, " ") {
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
{{range .Fields}}{{$name := n .Name}}    public {{if .Optional}}Optional<{{end}}{{t .T}}{{if .Optional}}>{{end}} {{$name}};
{{end}}
{{range .Fields}}{{$name := n .Name}}    public void Set{{$name}}({{t .T}} {{$name}}) { this.{{$name}} = {{$name}}; }
{{end}}
{{range .Fields}}{{$name := n .Name}}    public {{t .T}} Get{{$name}}({{t .T}} {{$name}}) { return this.{{$name}}; }
{{end}}



    public {{normalize .Name}}(JSONObject data)
    {
{{range .Fields}}        {{$name := n .Name}}        parse{{$name}}(data);
{{end}}
    }

{{range .Fields}}{{$name := n .Name}}    private void parse{{$name}}(JSONObject data)
    {
        final String key = "{{.Name}}";
        if(!data.containsKey(key))
        {
            missingID(key);
        }
        
        this.{{$name}} = data.get{{caps (t .T)}}(key).get{{if eq (t .T) "String"}}String{{else}}Val{{end}}();
    }
{{end}}
}
`

func generateSourceFile(s string, object syntaxtree.Object) error {
	tmpl, err := template.New("JavaSource").Funcs(template.FuncMap{
		"upper":     strings.ToUpper,
		"caps":      func(a string) string { return strings.Title(strings.ToLower(a)) },
		"n":         normalzieVariableNames,
		"normalize": normalizeFileName,
		"t":         typeSwitch,
	}).Parse(templateString)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, object)
	if err != nil {
		panic(err)
	}

	return nil
}
