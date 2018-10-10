package workflow

import (
	"fmt"
	"sync"

	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
)

type versionOrWorkflowDef struct {
	versionDef  *VersionDefinition
	workflowDef *Definition
}

// UnsafeManager is used by the agent, and thus must be exported. But a normal user of the library shouldn't use this
// directly.
var UnsafeManager = &Store{
	workflows: make(map[string]*versionOrWorkflowDef),
	mu:        &sync.RWMutex{},
}

// Store is a thread-safe store of workflow Definitions. This is used to insure that no two workflows can have the same name.
// It also will be used by the agent to be able take a workflow name (as well as any workflow data if it exists) and produce an
// Instance of that workflow.
type Store struct {
	workflows map[string]*versionOrWorkflowDef
	mu        *sync.RWMutex
}

// UnsafeGetDefinition is used by the agent, and thus must be exported. But a normal user of the library shouldn't use this
// directly.
func (wfm *Store) UnsafeGetDefinition(name string) *versionOrWorkflowDef {

	wfm.mu.RLock()
	def := wfm.workflows[name]
	wfm.mu.RUnlock()

	return def
}

// UnsafeGetInstance is used by the agent, and thus must be exported. But a normal user of the library shouldn't use this
// directly.
func (wfm *Store) UnsafeGetInstance(name, encodedData string) (*Instance, error) {

	def := wfm.UnsafeGetDefinition(name)

	if def == nil {
		panic(fmt.Sprint("unknown workflow: ", name))
	}

	if encodedData == `""` {
		encodedData = "{}"
	}

	var wfDef *Definition
	if def.versionDef != nil {
		// in this case the workflow was versioned while running.
		// so we get the initial workflow from the list of versions in the version definition
		wfDef = def.versionDef.getInitialDefinition()
	} else {
		wfDef = def.workflowDef
	}

	err := serializer.Decode(encodedData, wfDef.defaultInstance.Handler)

	return wfDef.defaultInstance, err
}

func (wfm *Store) setDefinition(name string, workflow *Definition) {
	if wfm.UnsafeGetDefinition(name) != nil {
		panic(fmt.Sprint("workflowDef definition with name '", name, "' already exists"))
	}
	wfm.mu.Lock()
	wfm.workflows[name] = &versionOrWorkflowDef{
		workflowDef: workflow,
	}
	wfm.mu.Unlock()
}

func (wfm *Store) setVersionDef(name string, versionDef *VersionDefinition) {
	if wfm.UnsafeGetDefinition(name) != nil {
		panic(fmt.Sprint("workflowDef definition with name '", name, "' already exists"))
	}
	wfm.mu.Lock()
	wfm.workflows[name] = &versionOrWorkflowDef{
		versionDef: versionDef,
	}
	wfm.mu.Unlock()
}
