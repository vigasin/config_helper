package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"os"
)

func main() {
	numArgs := len(os.Args)

	if numArgs < 3 || numArgs > 5 {
		fmt.Println(`Usage:\n\
			config_helper <file> <base64_json_data>
			config_helper <file> <section> <name>
			config_helper <file> <section> <name> <value>`)

		os.Exit(1)
	}

	filename := os.Args[1]
	cfg, err := ini.LoadSources(ini.LoadOptions{PreserveSurroundedQuote: true}, filename)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	if numArgs == 3 {
		b64Json := os.Args[2]

		jsonBytes, err := base64.StdEncoding.DecodeString(b64Json)
		if err != nil {
			fmt.Println("Failed to decode base64:", err)
			return
		}

		var dict map[string]interface{}
		err = json.Unmarshal(jsonBytes, &dict)
		if err != nil {
			fmt.Println("Failed to parse json:", err)
			return
		}

		for section, sectionData := range dict {
			for key, value := range sectionData.(map[string]interface{}) {
				cfg.Section(section).Key(key).SetValue(fmt.Sprintf("\"%v\"", value))
			}
		}

		cfg.SaveTo(filename)
	} else if numArgs == 4 {
		section := os.Args[2]
		key := os.Args[3]

		fmt.Println(cfg.Section(section).Key(key).String())
	} else {
		section := os.Args[2]
		key := os.Args[3]
		value := os.Args[3]

		cfg.Section(section).Key(key).SetValue(value)

		cfg.SaveTo(filename)
	}
}
