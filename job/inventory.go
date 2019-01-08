package job

import (
//"fmt"

//log "github.com/uthng/golog"
//utils "github.com/uthng/goutils"
)

/////////// DECLARATION OF ALL TYPES //////////////////////////

// Host describes structure of a host
type Host struct {
	Name string

	Groups []string
	Vars   map[string]interface{}
}

// Group describes attributes of a group
type Group struct {
	Name string

	Hosts []string
	Vars  map[string]interface{}
}

// Inventory describes attributes of a host inventory
type Inventory struct {
	Global map[string]interface{}
	Hosts  map[string]Host
	Groups map[string]Group
}

////////// DEFINITION OF ALL FUNCTIONS ///////////////////////////

// NewInventory instancies a new Inventory
func NewInventory() *Inventory {
	inventory := &Inventory{
		Global: make(map[string]interface{}),
		Hosts:  make(map[string]Host),
		Groups: make(map[string]Group),
	}

	return inventory
}
