package jobflow

import (
	"fmt"
	//"io/ioutil"
	"path/filepath"
	"plugin"

	log "github.com/uthng/golog"
)

////////// DECLARATION ALL TYPES ////////////////////

// ModuleRegistry is a registry for modules
type ModuleRegistry struct {
	ModuleList map[string]string
}

/////////// DECLARATION OF ALL GLOBAL VARIABLES //////////////

var moduleRegistry *ModuleRegistry

////////// DEFINITION OF ALL FUNCTIONS /////////////////

// init initializes a new registry for loaded modules
func init() {
	if moduleRegistry == nil {
		log.Debugln("Instancie a new module registry")
		moduleRegistry = &ModuleRegistry{}
		moduleRegistry.ModuleList = make(map[string]string)
	}

}

// GetModuleRegistry returns the module registry initialized
func GetModuleRegistry() *ModuleRegistry {
	return moduleRegistry
}

// LoadModules loads all modules (plugins) present
// in the module directory
//
// dir: string specifing the location of modules to load
func LoadModules(dir string) error {
	if dir == "" {
		return fmt.Errorf("module directory is nil")
	}

	log.Infoln("Loading modules from...:", dir)

	// Search only the file with extension .so
	pattern := dir + "/" + "*.so"

	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	log.Debugln("Searching modules...:", files)

	// Loop to open each plugin found
	for _, file := range files {
		log.Debugln("Openning module:", file)
		_, err := plugin.Open(file)
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		// Register to register opened module
		// Verify if module already exists in the registry
		name := filepath.Base(file)
		log.Debugln("Registering module:", name, file)

		_, ok := moduleRegistry.ModuleList[name]
		if ok == false {
			moduleRegistry.ModuleList[name] = file
			log.Debugln("Module registered:", name, file)
		}

	}

	return nil
}
