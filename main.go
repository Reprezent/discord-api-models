package main

import (
	"io/ioutil"
	"log"
	"os"

	"./langadaptor"
	"./syntaxtree"
)

type LangAdaptor interface {
	Generate([]*syntaxtree.Object) ([]string, error)
}

type ASTTranslator interface {
	Translate(data []byte) []*syntaxtree.Object
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
