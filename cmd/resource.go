/*
Copyright © 2024 Lawrence McDaniel <lpm0073@gmail.com>
Website: https://lawrencemcdaniel.com>
*/
package cmd

import (
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Article returns "an" for a vowel-leading display name, "a" otherwise. It's
// a small convenience for verb packages templating Short/Long text for
// resource kinds whose display name isn't known until it's declared.
func Article(display string) string {
	if len(display) > 0 && strings.ContainsRune("AEIOUaeiou", rune(display[0])) {
		return "an"
	}
	return "a"
}

// FlagKind identifies the pflag type a FlagSpec should register.
type FlagKind int

const (
	FlagString FlagKind = iota
	FlagBool
)

// FlagSpec declares one flag to attach to a resource leaf command. It is
// bound to viper under Name, so its value is read back with viper.GetString
// or viper.GetBool depending on Kind.
type FlagSpec struct {
	Name        string
	Shorthand   string
	Usage       string
	Kind        FlagKind
	Default     string // used when Kind == FlagString
	DefaultBool bool   // used when Kind == FlagBool
	Required    bool
	Choices     []string // optional: restrict a FlagString value to this set
}

// NameArgMode describes how a leaf command receives the resource's
// identifying parameter (usually "name", sometimes "username" or
// "session_key").
type NameArgMode int

const (
	// NameArgNone means the resource has no identifying parameter (e.g. "account").
	NameArgNone NameArgMode = iota
	// NameArgPositional means the parameter is a required positional argument, args[0].
	NameArgPositional
	// NameArgFlag means the parameter is supplied via a named flag.
	NameArgFlag
)

// NameArgSpec configures the resource's identifying parameter.
type NameArgSpec struct {
	Mode      NameArgMode
	Kwarg     string // kwargs key, e.g. "name", "username", "session_key"
	Shorthand string // flag shorthand, only used when Mode == NameArgFlag
	Usage     string
	Required  bool // MarkFlagRequired, only used when Mode == NameArgFlag
}

// ResourceSpec declares one leaf command: one resource kind under one verb
// (deploy/describe/get/manifest).
type ResourceSpec struct {
	Use     string
	Short   string
	Long    string
	APIKind string // the "kind" path segment sent to the Smarter API
	NameArg NameArgSpec
	Flags   []FlagSpec
}

// RequestFunc issues the verb-specific API request for a resource kind.
type RequestFunc func(kind string, kwargs map[string]string) ([]byte, error)

// OutputFunc renders a successful response body to the console.
type OutputFunc func(bodyJson []byte)

// ErrFunc renders an error to the console.
type ErrFunc func(error)

func mustBindPFlag(key string, flag *pflag.Flag) {
	if err := viper.BindPFlag(key, flag); err != nil {
		log.Fatalf("Error binding flag '%s': %v", key, err)
	}
}

func mustMarkRequired(c *cobra.Command, name string) {
	if err := c.MarkFlagRequired(name); err != nil {
		log.Fatalf("Error marking flag '%s' as required: %v", name, err)
	}
}

func validateChoices(name, value string, choices []string) {
	for _, choice := range choices {
		if value == choice {
			return
		}
	}
	log.Fatalf("Invalid value '%s' for flag '%s'. Allowed values are: %v", value, name, choices)
}

// RegisterResourceCmd builds a *cobra.Command from spec, registers its flags
// (bound to viper) and its Run function, and attaches it to parent. The Run
// function assembles kwargs from spec.NameArg and spec.Flags, calls
// request(spec.APIKind, kwargs), and dispatches the result to output or onErr.
func RegisterResourceCmd(parent *cobra.Command, spec ResourceSpec, request RequestFunc, output OutputFunc, onErr ErrFunc) *cobra.Command {
	leaf := &cobra.Command{
		Use:   spec.Use,
		Short: spec.Short,
		Long:  spec.Long,
		Run: func(c *cobra.Command, args []string) {
			kwargs := map[string]string{}

			switch spec.NameArg.Mode {
			case NameArgPositional:
				kwargs[spec.NameArg.Kwarg] = args[0]
			case NameArgFlag:
				kwargs[spec.NameArg.Kwarg] = viper.GetString(spec.NameArg.Kwarg)
			}

			for _, f := range spec.Flags {
				switch f.Kind {
				case FlagBool:
					kwargs[f.Name] = strconv.FormatBool(viper.GetBool(f.Name))
				default:
					val := viper.GetString(f.Name)
					if len(f.Choices) > 0 && val != "" {
						validateChoices(f.Name, val, f.Choices)
					}
					kwargs[f.Name] = val
				}
			}

			bodyJson, err := request(spec.APIKind, kwargs)
			if err != nil {
				onErr(err)
			} else {
				output(bodyJson)
			}
		},
	}

	if spec.NameArg.Mode == NameArgPositional {
		leaf.Args = cobra.ExactArgs(1)
	}

	if spec.NameArg.Mode == NameArgFlag {
		leaf.Flags().StringP(spec.NameArg.Kwarg, spec.NameArg.Shorthand, "", spec.NameArg.Usage)
		mustBindPFlag(spec.NameArg.Kwarg, leaf.Flags().Lookup(spec.NameArg.Kwarg))
		if spec.NameArg.Required {
			mustMarkRequired(leaf, spec.NameArg.Kwarg)
		}
	}

	for _, f := range spec.Flags {
		switch f.Kind {
		case FlagBool:
			leaf.Flags().BoolP(f.Name, f.Shorthand, f.DefaultBool, f.Usage)
		default:
			leaf.Flags().StringP(f.Name, f.Shorthand, f.Default, f.Usage)
		}
		mustBindPFlag(f.Name, leaf.Flags().Lookup(f.Name))
		if f.Required {
			mustMarkRequired(leaf, f.Name)
		}
	}

	parent.AddCommand(leaf)
	return leaf
}
