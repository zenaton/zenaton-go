package boot

import (
	// for this boot file to work, all workflows you wish to use must be exported, package level variables
	// inside of github.com/zenaton/zenaton-go/workflows/workflows
	_ "github.com/zenaton/zenaton-go/workflows"
)
