package all

import (
	// Enable plugin shell
	_ "github.com/uthng/jobflow/plugins/shell"
	// Enable plugin gox
	_ "github.com/uthng/jobflow/plugins/gox"
	// Enable plugin github
	_ "github.com/uthng/jobflow/plugins/github"
)
