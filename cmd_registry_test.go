package gojobs_test

import (
	//"fmt"
	"testing"

	"github.com/uthng/gojobs"
	log "github.com/uthng/golog"
)

var module = gojobs.Module{
	Name:        "ModTest",
	Version:     "0.1",
	Description: "ModTest",
}

var fn = func(map[string]interface{}) *gojobs.CmdResult {
	log.Debugln("CmdFunc test")
	return &gojobs.CmdResult{Error: nil, Result: nil}
}

var cmds = []gojobs.Cmd{
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
		gojobs.CmdRegister(cmd)
	}

	registry := gojobs.GetCmdRegistry()
	log.Debugf("Registry %#v\n", registry)

	nb := gojobs.GetNbOfCmds()
	log.Debugf("nb of commands %v\n", nb)

	if nb != 3 {
		t.Fail()
	}
}

func TestCmdUnregister(t *testing.T) {

	gojobs.CmdUnregister(cmds[1])

	registry := gojobs.GetCmdRegistry()
	log.Debugf("Registry %#v\n", registry)

	nb := gojobs.GetNbOfCmds()
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
	cmd, ok := gojobs.GetCmdByName("ModTest.cmd3")
	if ok {
		cmd.Func(nil)
	} else {
		t.Fail()
	}
}
