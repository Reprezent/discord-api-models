package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/Reprezent/discord-api-models/langadaptor"
	"github.com/Reprezent/discord-api-models/syntaxtree"
)

type LangAdaptor interface {
	Generate(map[string]*syntaxtree.Object) ([]string, error)
}

type ASTTranslator interface {
	Translate(data []byte) map[string]*syntaxtree.Object
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need a file argument")
	}

	fileData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var xlator syntaxtree.GoldMarkTranslator
	stuff := xlator.Translate(fileData)

	// json_str, err := json.MarshalIndent(stuff, "", "  ")

	var gene langadaptor.JavaAdaptor

	_, _ = gene.Generate(stuff)

	if err != nil {
		log.Fatalf(err.Error())
	}

	// log.Printf("%v", string(json_str))
}
