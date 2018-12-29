package jobflow_test

import (
	//"fmt"
	"testing"

	"github.com/uthng/jobflow"
	_ "github.com/uthng/jobflow/plugins/all"
	log "github.com/uthng/golog"
)

func TestLoadModules(t *testing.T) {
	log.SetVerbosity(log.DEBUG)
	//jobflow.NewCmdRegistry()
	//jobflow.NewModuleRegistry()

	//jobflow.LoadModules("../app/modules")

	//registry := jobflow.GetModuleRegistry()
	//fmt.Printf("Registry %#v\n", registry)

	cmd, ok := jobflow.GetCmdByName("shell.exec")
	if ok {
		result := cmd.Func(map[string]interface{}{"cmd": "ls -la"})
		log.Infoln(result)
	} else {
		t.Fail()
	}
}
