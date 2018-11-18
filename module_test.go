package gojobs_test

import (
	//"fmt"
	"testing"

	"github.com/uthng/gojobs"
	_ "github.com/uthng/gojobs/modules/all"
	log "github.com/uthng/golog"
)

func TestLoadModules(t *testing.T) {
	log.SetVerbosity(log.DEBUG)
	//gojobs.NewCmdRegistry()
	//gojobs.NewModuleRegistry()

	//gojobs.LoadModules("../app/modules")

	//registry := gojobs.GetModuleRegistry()
	//fmt.Printf("Registry %#v\n", registry)

	cmd, ok := gojobs.GetCmdByName("shell.ExecCmd")
	if ok {
		result := cmd.Func(map[string]interface{}{"cmd": "ls -la"})
		log.Infoln(result)
	} else {
		t.Fail()
	}
}
