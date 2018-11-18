package shell

import "C"
import (
	"fmt"
	"os/exec"
	//"strings"

	"github.com/uthng/gojobs"
)

var module = gojobs.Module{
	Name:        "shell",
	Version:     "0.1",
	Description: "Everything concern shell",
}

// List of available commands for this module
var commands = []gojobs.Cmd{
	{
		Name:   "ExecCmd",
		Func:   ExecCmd,
		Module: module,
	},
}

// Init initializes module by registering all its commands
// to command registry
func init() {
	for _, cmd := range commands {
		gojobs.CmdRegister(cmd)
	}
}

// ExecCmd executes a command shell (bash).
// It takes a map of params
func ExecCmd(params map[string]interface{}) *gojobs.CmdResult {
	//var command []string
	var res = gojobs.NewCmdResult()

	value, ok := params["cmd"]
	if ok == false {
		res.Error = fmt.Errorf("param pkgname missing")
		return res
	}
	//command = strings.Fields(value.(string))
	command := value.(string)

	// Execute kubectl command
	//cmd := exec.Command(command[0], command[1:len(command)]...)
	cmd := exec.Command("bash", "-c", command)

	// Check if error
	output, err := cmd.Output()
	if err != nil {
		res.Error = err
		return res
	}

	res.Result["result"] = string(output)
	return res
}
