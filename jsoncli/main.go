package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kingzbauer/json_cli/jsongear"
)

var (
	key         = flag.String("k", "", "The key value for the field you want to access. Separate nested keys using `.`")
	file        = flag.String("f", "", "Json file to search")
	indent      = flag.Int("indent", 4, "Indent level for the json output")
	indentStr   = flag.String("indentStr", " ", "String used to indent")
	listKeys    = flag.Bool("l", false, "List all keys under -k")
	search      = flag.Bool("s", false, "Search for the key")
	searchDepth = flag.Int("d", 3, "Search depth")
)

func main() {
	flag.Parse()
	if len(*key) == 0 && !*listKeys {
		fmt.Printf("Field `%s` is required\n", "-k")
		flag.Usage()
		os.Exit(1)
	}

	if *listKeys && *search {
		fmt.Println("You can't provide both -s and -l")
		os.Exit(1)
	}

	content, err := readContent(*file)
	must(err)

	// check if any content is present
	if len(content) == 0 {
		fmt.Println("No file content could be read. Either pass a file name via the `-f` flag or through stdin")
		os.Exit(0)
	}

	v, err := jsongear.Parse(content)
	must(err)

	if *listKeys {
		listkeys(*key, v)
	} else if *search {
		result := jsongear.Search(*key, v, *searchDepth)
		printJSON(result)
	} else {
		result := jsongear.Get(*key, v)
		printJSON(result)
	}
}

func readContent(file string) ([]byte, error) {
	if len(file) > 0 {
		return ioutil.ReadFile(file)
	}
	// Try reading from Stdin
	stat, _ := os.Stdin.Stat()
	// check if the input is being piped in
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return ioutil.ReadAll(os.Stdin)
	}

	return make([]byte, 0), nil
}

func printJSON(v interface{}) {
	// if result is among any of the concrete types, print as is
	switch t := v.(type) {
	case float64, bool, string, nil:
		fmt.Println(t)
	default:
		// format output for non concrete types
		buffer := new(bytes.Buffer)
		resultBytes, _ := json.Marshal(v)
		json.Indent(buffer, resultBytes, "", getIndentString(*indent, *indentStr))
		fmt.Println(buffer.String())
	}
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getIndentString(indentLevel int, indentStr string) string {
	return strings.Repeat(indentStr, indentLevel)
}

func listkeys(root string, content interface{}) {
	keys := jsongear.ListKeys(root, content)
	if key == nil {
		fmt.Println("No keys")
		return
	}
	// print keys, one by one
	for _, k := range keys {
		fmt.Println("-", k)
	}
}
