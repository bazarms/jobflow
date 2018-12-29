package jobflow

import (
	//    "fmt"
	log "github.com/uthng/golog"
)

/////// DECLARATION OF ALL TYPES /////////////

// ValueRegistry is a registry of values returned by each task
//
// It is map containing: Key: string for context or step name and
// Value: interface{} storing step result or context values
type ValueRegistry struct {
	ValueList map[string]interface{}
}

///////// DECLARATION OF ALL GLOBAL VARIABLES ///////////

///////// DECLARATION OF ALL FUNCTIONS /////////////////

// NewValueRegistry initialize a unique instance of value registry
func NewValueRegistry() *ValueRegistry {
	log.Debugln("Initializing a new value registry")
	v := &ValueRegistry{
		ValueList: make(map[string]interface{}),
	}

	return v
}

// AddValue adds a new value to registry
func (v *ValueRegistry) AddValue(key string, value interface{}) error {

	// Verify if the key already exists in the registry
	_, ok := v.ValueList[key]
	if ok == false {
		v.ValueList[key] = value
		log.Debugw("Value added", "key", key, "value", value)
	}

	return nil
}

// DeleteValue removes a existing value in the registry
func (v *ValueRegistry) DeleteValue(key string) error {

	// Remove if value exists
	_, ok := v.ValueList[key]
	if ok == true {
		delete(v.ValueList, key)
		log.Debugw("Value deleted", "key", key)
	}

	return nil
}

// GetNbOfValues returns the number of commands in the registry
func (v *ValueRegistry) GetNbOfValues() int {
	return len(v.ValueList)
}

// GetValueByKey return value of a key
func (v *ValueRegistry) GetValueByKey(key string) (interface{}, bool) {
	value, ok := v.ValueList[key]
	return value, ok
}
