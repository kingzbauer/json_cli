package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"bytes"
	"encoding/json"
	"github.com/kingzbauer/json_cli/utils"
)

var (
	key  = flag.String("k", "", "The key value for the field you want to access. Separate nested keys using `.`")
	file = flag.String("f", "", "Json file to search")
)

func main() {
	flag.Parse()
	if len(*key) == 0 {
		fmt.Printf("Field `%s` is required\n", "-k")
		flag.Usage()
		os.Exit(1)
	}

	var (
		content []byte
		err     error
	)
	if len(*file) > 0 {
		content, err = ioutil.ReadFile(*file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		// try reading from Stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			content, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	// check if any content is present
	if len(content) == 0 {
		fmt.Println("No file content could be read. Either pass a file name via the `-f` flag or through stdin")
		os.Exit(0)
	}

	v, err := utils.Parse(content)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	result := utils.Get(*key, v)

	// if result is among any of the concrete types, print as is
	switch t := result.(type) {
	case float64, bool, string, nil:
		fmt.Println(t)
	default:
		// format output for non concrete types
		buffer := new(bytes.Buffer)
		resultBytes, _ := json.Marshal(result)
		json.Indent(buffer, resultBytes, "", "  ")
		fmt.Println(buffer.String())
	}
}
