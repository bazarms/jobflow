package gojobs_test

import (
	"fmt"
	"testing"

	"github.com/uthng/gojobs"
)

func TestLoadModules(t *testing.T) {
	gojobs.NewCmdRegistry()
	gojobs.NewModuleRegistry()

	gojobs.LoadModules("../app/modules")

	registry := gojobs.GetModuleRegistry()
	fmt.Printf("Registry %#v\n", registry)

	cmd, ok := gojobs.GetCmdByName("shell.ExecCmd")
	if ok {
		cmd.Func(nil)
	} else {
		t.Fail()
	}
}
