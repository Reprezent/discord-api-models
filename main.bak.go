// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/Reprezent/discord-api-models/syntaxtree"
// 	"github.com/urfave/cli/v2"
// )

// type LangAdaptor interface {
// 	Generate(map[string]*syntaxtree.Object) ([]string, error)
// }

// type ASTTranslator interface {
// 	Translate(data []byte) map[string]*syntaxtree.Object
// }

// func dfsPrintObjects(node *syntaxtree.Object, ident int) {
// 	log.Printf("%*c%s - %p\n", ident, ' ', node.Name, node)
// 	ident += 4
// 	for _, val := range node.Fields {
// 		log.Printf("%*cName:          %s\n", ident, ' ', val.Name)
// 		log.Printf("%*cReferenceName: %s\n", ident, ' ', val.ReferenceName)
// 		log.Printf("%*cRef:           %p\n", ident, ' ', val.Reference)
// 	}
// }

// func main() {
// 	// if len(os.Args) < 2 {
// 	// 	log.Fatal("Need a file argument")
// 	// }

// 	// var totalData []byte
// 	// files := os.Args[1:]
// 	// for _, file := range files {
// 	// 	fileData, err := ioutil.ReadFile(file)
// 	// 	totalData = append(totalData, fileData...)
// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}
// 	// }

// 	// var xlator syntaxtree.GoldMarkTranslator
// 	// stuff := xlator.Translate(totalData)

// 	// for _, val := range stuff {
// 	// 	dfsPrintObjects(val, 0)
// 	// }

// 	// // json_str, err := json.MarshalIndent(stuff, "", "  ")

// 	// var gene langadaptor.JavaAdaptor

// 	// files, data, err := gene.Generate(stuff)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// for i := range files {
// 	// 	p := path.Join("java", files[i])
// 	// 	log.Printf("Writing to %v file", p)
// 	// 	err := os.MkdirAll(path.Dir(p), 0755)
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	err = ioutil.WriteFile(p, []byte(data[i]), 0644)
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// }

// 	app := &cli.App{
// 		Name:  "discord-api-models",
// 		Usage: "generates object models for Discord API SDKs",
// 	}

// 	app.Commands = []*cli.Command{
// 		{
// 			Name:    "generate",
// 			Aliases: []string{"g"},
// 			Flags: []cli.Flag{
// 				&cli.StringSliceFlag{
// 					Name:     "input-files",
// 					Aliases:  []string{"i", "input"},
// 					Usage:    "Input markdown files from the discord api docs (https://github.com/discord/discord-api-docs)",
// 					Required: true,
// 				},
// 			},
// 			Action: func(c *cli.Context) error {
// 				fmt.Printf("Hello %q", c.StringSlice("input-files"))
// 				return nil
// 			},
// 		},
// 	}
// 	err := app.Run(os.Args)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// log.Printf("%v", string(json_str))
// }
