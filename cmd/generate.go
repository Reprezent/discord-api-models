/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Reprezent/discord-api-models/langadaptor"
	"github.com/Reprezent/discord-api-models/syntaxtree"
	"github.com/spf13/cobra"
)

// generateCommand represents the base command when called without any subcommands
var generateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generates source code",
	// Long: `Discord api models `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		i, _ := cmd.Flags().GetStringSlice("input")
		fmt.Println(i)
		o, _ := cmd.Flags().GetString("output")
		fmt.Println(o)

		entry(i, o)
	},
}

func init() {
	rootCmd.AddCommand(generateCommand)
	generateCommand.Flags().StringSliceP("input", "i", nil, "Input files to use for generating source code.")
	generateCommand.Flags().StringP("output", "o", "", "Output directory to place source code.")
	generateCommand.MarkFlagFilename("input")
	generateCommand.MarkFlagRequired("input")
	generateCommand.MarkFlagDirname("output")
	generateCommand.MarkFlagRequired("output")
}

type LangAdaptor interface {
	Generate(map[string]*syntaxtree.Object) ([]string, error)
}

type ASTTranslator interface {
	Translate(data []byte) map[string]*syntaxtree.Object
}

func dfsPrintObjects(node *syntaxtree.Object, ident int) {
	log.Printf("%*c%s - %p\n", ident, ' ', node.Name, node)
	ident += 4
	for _, val := range node.Fields {
		log.Printf("%*cName:          %s\n", ident, ' ', val.Name)
		log.Printf("%*cReferenceName: %s\n", ident, ' ', val.ReferenceName)
		log.Printf("%*cRef:           %p\n", ident, ' ', val.Reference)
	}
}
func getAllFiles(inputs []string) []string {
	var rv []string
	for _, i := range inputs {
		err := filepath.Walk(i,
			func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					rv = append(rv, path)
				}
				return err
			})
		if err != nil {
			log.Fatal(err)
		}
	}

	return rv
}
func entry(inputs []string, outputDirectory string) {

	var totalData []byte
	files := getAllFiles(inputs)
	log.Printf("Files writing %v", files)
	for _, file := range files {
		fileData, err := ioutil.ReadFile(file)
		totalData = append(totalData, fileData...)
		if err != nil {
			log.Fatal(err)
		}
	}

	var xlator syntaxtree.GoldMarkTranslator
	stuff := xlator.Translate(totalData)

	for _, val := range stuff {
		dfsPrintObjects(val, 0)
	}

	// json_str, err := json.MarshalIndent(stuff, "", "  ")

	var gene langadaptor.JavaAdaptor

	files, data, err := gene.Generate(stuff)
	if err != nil {
		panic(err)
	}

	for i := range files {
		p := path.Join(outputDirectory, files[i])
		log.Printf("Writing to %v file", p)
		err := os.MkdirAll(path.Dir(p), 0755)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(p, []byte(data[i]), 0644)
		if err != nil {
			panic(err)
		}
	}

	// log.Printf("%v", string(json_str))
}
