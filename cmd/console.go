package cmd

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
)

var jsonFlagValue bool
var yamlFlagValue bool

func init() {
	jsonFlagValue = viper.GetBool("json")
	yamlFlagValue = viper.GetBool("yaml")
}

func ConsoleOutput(bodyJson []byte) {
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
