package cmd

import (
	"fmt"

	"github.com/ghodss/yaml"
)

func ConsoleOutput(bodyJson []byte, jsonFlagValue bool, yamlFlagValue bool) {
	switch {
	case jsonFlagValue:
		fmt.Println(string(bodyJson))
	case yamlFlagValue:
		bodyYaml, err := yaml.JSONToYAML(bodyJson)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(string(bodyYaml))
		}
	default:
		fmt.Println(string(bodyJson))
	}
}
