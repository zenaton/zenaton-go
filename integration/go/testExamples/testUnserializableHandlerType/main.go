package main

import (
	"github.com/zenaton/integration/go/client" // initialize zenaton client with credentials
	"github.com/zenaton/integration/go/tasks"
	"github.com/zenaton/integration/go/workflows"
)

func init() { client.SetEnv("testUnserializableHandlerType.env") }
func main() {
	//workflows.TestW.WhereID("MyID").Kill()
	t := tasks.TestTask.New()

	workflows.WithTaskWorkflow.New(t).Dispatch()
}
