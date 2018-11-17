package gojobs

import (
	//    "fmt"
	log "github.com/uthng/golog"
)

/////// DECLARATION OF ALL TYPES /////////////

// CmdResult represents the result of a module's command executed.
type CmdResult struct {
	// Error is command error
	Error error
	// Result is a map containing output of each task
	Result map[string]interface{}
}

// CmdFunc is a command function
type CmdFunc func(map[string]interface{}) *CmdResult

// Module contains les informations of a module
type Module struct {
	// Name is module name
	Name string
	// Version is module version
	Version string
	// Description describes shortly what the module does
	Description string
}

// Cmd is a structure for a command of a given module
type Cmd struct {
	// Name is command name
	Name string
	// Func is command function
	Func CmdFunc
	// Module is the module to which the command belongs to
	Module Module
}

// CmdRegistry is a registry for commands
//
// This is a map with Key: string composed of module name and command name
// and Cmd: command structure
type CmdRegistry struct {
	CmdList map[string]Cmd
}

///////// DECLARATION OF ALL GLOBAL VARIABLES ///////////

var cmdRegistry *CmdRegistry

///////// DECLARATION OF ALL FUNCTIONS /////////////////

// NewCmdRegistry initialize a unique instance of command registry
func NewCmdRegistry() {
	log.Debugln("Instanciate new command registry")
	if cmdRegistry == nil {
		cmdRegistry = &CmdRegistry{}
		cmdRegistry.CmdList = make(map[string]Cmd)
	}
}

// NewCmdResult instanciates a new command result
func NewCmdResult() *CmdResult {
	log.Debugln("Instanciate a new command result")
	c := &CmdResult{
		Error:  nil,
		Result: make(map[string]interface{}),
	}

	return c
}

// GetCmdRegistry returns the command registry initialized
func GetCmdRegistry() *CmdRegistry {
	return cmdRegistry
}

// CmdRegister registers a new command
// cmd: Command to register
func CmdRegister(cmd Cmd) error {
	// Name in commande registry = <module name>.<cmd name>
	var name = cmd.Module.Name + "." + cmd.Name

	// Verify if command already exists in the registry
	_, ok := cmdRegistry.CmdList[name]
	if ok == false {
		cmdRegistry.CmdList[name] = cmd
		log.Debugln("Command registered:", name)
	}

	return nil
}

// CmdUnregister unregister a command in command registry
// cmd: command to unregister
func CmdUnregister(cmd Cmd) error {
	var name = cmd.Module.Name + "." + cmd.Name

	// Remove if cmd exists
	_, ok := cmdRegistry.CmdList[name]
	if ok == true {
		delete(cmdRegistry.CmdList, name)
		log.Debugln("Unregister command:", name)

	}

	return nil
}

// GetNbOfCmds returns the number of commands in the registry
func GetNbOfCmds() int {
	return len(cmdRegistry.CmdList)
}

// GetCmdByName returns a command by its name
func GetCmdByName(name string) (Cmd, bool) {
	cmd, ok := cmdRegistry.CmdList[name]
	return cmd, ok
}
