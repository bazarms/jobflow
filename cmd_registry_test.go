package jobflow_test

import (
	//"fmt"
	"testing"

	"github.com/uthng/jobflow"
	log "github.com/uthng/golog"
)

var module = jobflow.Module{
	Name:        "ModTest",
	Version:     "0.1",
	Description: "ModTest",
}

var fn = func(map[string]interface{}) *jobflow.CmdResult {
	log.Debugln("CmdFunc test")
	return &jobflow.CmdResult{Error: nil, Result: nil}
}

var cmds = []jobflow.Cmd{
	{
		Name:   "cmd1",
		Func:   fn,
		Module: module,
	},
	{
		Name:   "cmd2",
		Func:   fn,
		Module: module,
	},
	{
		Name:   "cmd3",
		Func:   fn,
		Module: module,
	},
}

func TestCmdRegister(t *testing.T) {

	for _, cmd := range cmds {
		log.Debugln(cmd)
		jobflow.CmdRegister(cmd)
	}

	registry := jobflow.GetCmdRegistry()
	log.Debugf("Registry %#v\n", registry)

	nb := jobflow.GetNbOfCmds()
	log.Debugf("nb of commands %v\n", nb)

	if nb != 3 {
		t.Fail()
	}
}

func TestCmdUnregister(t *testing.T) {

	jobflow.CmdUnregister(cmds[1])

	registry := jobflow.GetCmdRegistry()
	log.Debugf("Registry %#v\n", registry)

	nb := jobflow.GetNbOfCmds()
	log.Debugf("nb of commands %v\n", nb)

	if nb != 2 {
		t.Fail()
	}

	_, ok := registry.CmdList["cmd1"]
	if ok {
		t.Fail()
	}

	_, ok = registry.CmdList["cmd2"]
	if ok {
		t.Fail()
	}

}

func TestGetCmdByName(t *testing.T) {
	cmd, ok := jobflow.GetCmdByName("ModTest.cmd3")
	if ok {
		cmd.Func(nil)
	} else {
		t.Fail()
	}
}
