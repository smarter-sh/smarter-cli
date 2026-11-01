package cmd

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// fakeRequest records the kind/kwargs it was called with and returns a
// canned response, standing in for a real APIRequest wrapper in tests.
type fakeRequest struct {
	calledKind   string
	calledKwargs map[string]string
	response     []byte
	err          error
}

func (f *fakeRequest) Request(kind string, kwargs map[string]string) ([]byte, error) {
	f.calledKind = kind
	f.calledKwargs = kwargs
	return f.response, f.err
}

func TestRegisterResourceCmdPositionalNameArgAndFlags(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	spec := ResourceSpec{
		Use:     "widget <name>",
		APIKind: "Widget",
		NameArg: NameArgSpec{Mode: NameArgPositional, Kwarg: "name"},
		Flags: []FlagSpec{
			{Name: "class", Kind: FlagString, Default: "static", Choices: []string{"static", "sql"}},
			{Name: "verbose-flag", Kind: FlagBool, DefaultBool: false},
		},
	}

	fr := &fakeRequest{response: []byte(`{"ok":true}`)}
	var gotOutput []byte
	RegisterResourceCmd(parent, spec, fr.Request, func(b []byte) { gotOutput = b }, func(err error) { t.Fatalf("onErr called unexpectedly: %v", err) })

	parent.SetArgs([]string{"widget", "my-widget", "--class", "sql", "--verbose-flag"})
	if err := parent.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if fr.calledKind != "Widget" {
		t.Errorf("request called with kind=%q, want Widget", fr.calledKind)
	}
	want := map[string]string{"name": "my-widget", "class": "sql", "verbose-flag": "true"}
	if !reflect.DeepEqual(fr.calledKwargs, want) {
		t.Errorf("kwargs = %v, want %v", fr.calledKwargs, want)
	}
	if string(gotOutput) != `{"ok":true}` {
		t.Errorf("output = %q", gotOutput)
	}
}

func TestRegisterResourceCmdFlagNameArg(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	spec := ResourceSpec{
		Use:     "widgets",
		APIKind: "Widget",
		NameArg: NameArgSpec{Mode: NameArgFlag, Kwarg: "name", Shorthand: "n", Usage: "widget name"},
	}

	fr := &fakeRequest{response: []byte(`{}`)}
	RegisterResourceCmd(parent, spec, fr.Request, func(b []byte) {}, func(err error) { t.Fatalf("onErr: %v", err) })

	parent.SetArgs([]string{"widgets", "-n", "my-widget"})
	if err := parent.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if fr.calledKwargs["name"] != "my-widget" {
		t.Errorf("kwargs[name] = %q, want my-widget", fr.calledKwargs["name"])
	}
}

func TestRegisterResourceCmdNoNameArg(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	spec := ResourceSpec{Use: "account", APIKind: "Account"}

	fr := &fakeRequest{response: []byte(`{}`)}
	RegisterResourceCmd(parent, spec, fr.Request, func(b []byte) {}, func(err error) { t.Fatalf("onErr: %v", err) })

	parent.SetArgs([]string{"account"})
	if err := parent.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(fr.calledKwargs) != 0 {
		t.Errorf("kwargs = %v, want empty", fr.calledKwargs)
	}
}

func TestRegisterResourceCmdRequestErrorInvokesOnErr(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	spec := ResourceSpec{Use: "widget <name>", APIKind: "Widget", NameArg: NameArgSpec{Mode: NameArgPositional, Kwarg: "name"}}

	fr := &fakeRequest{err: errors.New("boom")}
	outputCalled := false
	var gotErr error
	RegisterResourceCmd(parent, spec, fr.Request,
		func(b []byte) { outputCalled = true },
		func(err error) { gotErr = err },
	)

	parent.SetArgs([]string{"widget", "my-widget"})
	if err := parent.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if outputCalled {
		t.Error("output was called despite a request error")
	}
	if gotErr == nil || gotErr.Error() != "boom" {
		t.Errorf("onErr got %v, want boom", gotErr)
	}
}

func TestRegisterResourceCmdPositionalArgIsRequired(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	spec := ResourceSpec{Use: "widget <name>", APIKind: "Widget", NameArg: NameArgSpec{Mode: NameArgPositional, Kwarg: "name"}}

	fr := &fakeRequest{}
	RegisterResourceCmd(parent, spec, fr.Request, func(b []byte) {}, func(err error) {})

	parent.SetArgs([]string{"widget"})
	if err := parent.Execute(); err == nil {
		t.Error("Execute() error = nil, want an error for a missing required positional arg")
	}
}

func TestRegisterResources(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	specs := []ResourceSpec{
		{Use: "widget <name>", APIKind: "Widget", NameArg: NameArgSpec{Mode: NameArgPositional, Kwarg: "name"}},
		{Use: "gadget <name>", APIKind: "Gadget", NameArg: NameArgSpec{Mode: NameArgPositional, Kwarg: "name"}},
	}
	fr := &fakeRequest{response: []byte(`{}`)}
	RegisterResources(parent, specs, fr.Request, func(b []byte) {}, func(err error) {})

	names := map[string]bool{}
	for _, c := range parent.Commands() {
		names[strings.Fields(c.Use)[0]] = true
	}
	if !names["widget"] || !names["gadget"] {
		t.Errorf("registered commands = %v, want widget and gadget", names)
	}
}

func TestAdaptOutput(t *testing.T) {
	called := false
	adapted := AdaptOutput(func() { called = true })
	adapted([]byte("ignored"))
	if !called {
		t.Error("AdaptOutput's wrapped function was not called")
	}
}

func TestValidateChoicesAcceptsAllowedValue(t *testing.T) {
	if err := validateChoices("class", "sql", []string{"static", "sql"}); err != nil {
		t.Errorf("validateChoices() = %v, want nil", err)
	}
}

func TestValidateChoicesRejectsDisallowedValue(t *testing.T) {
	err := validateChoices("class", "bogus", []string{"static", "sql"})
	if err == nil || !strings.Contains(err.Error(), "invalid value") {
		t.Errorf("validateChoices() = %v, want an error explaining the invalid value", err)
	}
}

// TestRegisterResourceCmdReportsADisallowedChoice checks that a leaf command
// passes an invalid flag choice to its error function, without a request.
func TestRegisterResourceCmdReportsADisallowedChoice(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	var gotErr error
	requested := false
	leaf := RegisterResourceCmd(parent, ResourceSpec{
		Use:     "things",
		APIKind: "Thing",
		Flags:   []FlagSpec{{Name: "choicetest", Kind: FlagString, Choices: []string{"a", "b"}}},
	}, func(string, map[string]string) ([]byte, error) {
		requested = true
		return nil, nil
	}, func([]byte) {}, func(err error) { gotErr = err })

	withViperValue(t, "choicetest", "c")
	leaf.Run(leaf, nil)

	if gotErr == nil || requested {
		t.Errorf("err = %v, requested = %v; want an error and no request", gotErr, requested)
	}
}

// TestRegisterResourceCmdReadsItsOwnFlags checks that two leaves whose flags
// share a viper key each send their own value: the later leaf's binding used
// to win, so the earlier one sent its flag empty.
func TestRegisterResourceCmdReadsItsOwnFlags(t *testing.T) {
	parent := &cobra.Command{Use: "parent"}
	fr := &fakeRequest{response: []byte(`{}`)}
	for _, use := range []string{"widgets", "gadgets"} {
		RegisterResourceCmd(parent, ResourceSpec{
			Use:     use,
			APIKind: use,
			NameArg: NameArgSpec{Mode: NameArgFlag, Kwarg: "shared_name", Param: "name", Shorthand: "n"},
			Flags:   []FlagSpec{{Name: "shared_flag", Kind: FlagString}},
		}, fr.Request, func(b []byte) {}, func(err error) { t.Fatalf("onErr: %v", err) })
	}

	parent.SetArgs([]string{"widgets", "-n", "my-widget", "--shared_flag", "x"})
	if err := parent.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	want := map[string]string{"name": "my-widget", "shared_flag": "x"}
	if !reflect.DeepEqual(fr.calledKwargs, want) {
		t.Errorf("kwargs = %v, want %v", fr.calledKwargs, want)
	}
}
