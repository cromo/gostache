package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"code.google.com/p/go.crypto/ssh/terminal"
	"github.com/hoisie/mustache"
	"launchpad.net/goyaml"
)

func main() {
	if passedHelpArgument() || isTerminal(os.Stdin) && noFilesPassed() {
		printHelp()
		return
	}
	if passedVersionArgument() {
		printVersion()
		return
	}

	files := getFileList()
	input := concatContents(files)
	template, data := splitTemplateAndData(input)
	renderTemplateWithData(template, data)
}

func passedHelpArgument() bool {
	return argumentsContain("-h") || argumentsContain("--help")
}

func argumentsContain(value string) bool {
	for _, argument := range os.Args[1:] {
		if argument == value {
			return true
		}
	}
	return false
}

func isTerminal(file *os.File) bool {
	return terminal.IsTerminal(int(file.Fd()))
}

func noFilesPassed() bool {
	return len(os.Args) == 1
}

func printHelp() {
	fmt.Printf(`Usage: gostache FILE ...

Examples:
  $ gostache data.yml template.mustache
  $ cat data.yml | gostache - template.mustache

  This operates similarly to mustache(1)
  (available online at http://mustache.github.com/mustache.1.html), but does
  not support the -c, -t, or -r flags.

Common Options:
    -v, --version                    Print the version
    -h, --help                       Show this message
`)
}

func passedVersionArgument() bool {
	return argumentsContain("-v") || argumentsContain("--version")
}

func printVersion() {
	fmt.Printf("Gostache v0.1.0\n")
}

func getFileList() []string {
	if noFilesPassed() {
		return []string{"-"}
	}
	return os.Args[1:]
}

func concatContents(files []string) string {
	var fileContents []string
	for _, file := range files {
		contents, err := readFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read file %s: %s\n", file, err.Error())
			os.Exit(1)
		}
		fileContents = append(fileContents, contents)
	}

	return strings.Join(fileContents, "")
}

func readFile(filename string) (string, error) {
	if filename == "-" {
		contents, err := ioutil.ReadAll(os.Stdin)
		return string(contents), err
	}
	contents, err := ioutil.ReadFile(filename)
	return string(contents), err
}

func splitTemplateAndData(input string) (string, []string) {
	yamlDocuments := splitYamlDocuments(input)
	if len(yamlDocuments) <= 2 {
		return input, []string{""}
	}
	return yamlDocuments[0] + yamlDocuments[len(yamlDocuments)-1], yamlDocuments[1 : len(yamlDocuments)-1]
}

func splitYamlDocuments(yaml string) []string {
	r := regexp.MustCompile("(?m:^---\n)")
	return r.Split(yaml, 5)
}

func renderTemplateWithData(template string, data []string) {
	t, err := mustache.ParseString(template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing mustache template: %s\n", err.Error())
		os.Exit(1)
	}
	for i, yaml := range data {
		deserialized, err := deserializeYaml(yaml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing YAML blob %d: %s\n", i, err.Error())
			os.Exit(1)
		}
		fmt.Printf(t.Render(deserialized))
	}
}

func deserializeYaml(yaml string) (map[string]interface{}, error) {
	deserialized := make(map[string]interface{})
	err := goyaml.Unmarshal([]byte(yaml), &deserialized)
	return deserialized, err
}
