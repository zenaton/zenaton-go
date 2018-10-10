package main

import (
	"github.com/zenaton/zenaton-go/integration/go/client" // initialize zenaton client with credentials
	"github.com/zenaton/zenaton-go/integration/go/tasks"
	"github.com/zenaton/zenaton-go/integration/go/workflows"
)

func init() { client.SetEnv("testUnserializableHandlerType.env") }
func main() {
	//workflows.TestW.WhereID("MyID").Kill()
	t := tasks.TestTask.New()

	workflows.WithTaskWorkflow.New(t).Dispatch()
}
