package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-ini/ini"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	numArgs := len(os.Args)

	if numArgs < 3 || numArgs > 5 {
		fmt.Println(`Usage:\n\
			config_helper <file> <base64_yaml_data>
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

		bytes, err := base64.StdEncoding.DecodeString(b64Json)
		if err != nil {
			fmt.Println("Failed to decode base64:", err)
			return
		}

		var dict map[interface{}]interface{}
		err = yaml.Unmarshal(bytes, &dict)
		if err != nil {
			fmt.Println("Failed to parse yaml:", err)
			return
		}

		for section, sectionData := range dict {
			for key, value := range sectionData.(map[interface{}]interface{}) {
				cfg.Section(section.(string)).Key(key.(string)).SetValue(fmt.Sprintf("\"%v\"", value))
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
