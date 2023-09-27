package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-ini/ini"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
	"strings"
)

var Version string;

func main() {
	numArgs := len(os.Args)

    if numArgs == 2 && os.Args[1] == "--version" {
		fmt.Println(Version)

		os.Exit(1)
    } else if numArgs < 3 || numArgs > 5 {
		fmt.Println(`Usage:
			config_helper <file> <base64_yaml_data>
			config_helper <file> <section> <name>
			config_helper <file> <section> <name> <value>
			config_helper --version`)

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
				rt := reflect.TypeOf(value)

				var strValue string

				switch rt.Kind() {
				case reflect.Slice, reflect.Array:
					value := value.([] interface {})
					valueStrList := make([]string, len(value))
					for i, v := range value {
						valueStrList[i] = fmt.Sprint(v)
					}
					strValue = strings.Join(valueStrList, ",")
				default:
					strValue = fmt.Sprintf("\"%v\"", value)

				}
				cfg.Section(section.(string)).Key(key.(string)).SetValue(strValue)
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
		value := os.Args[4]

		cfg.Section(section).Key(key).SetValue(value)

		cfg.SaveTo(filename)
	}
}
