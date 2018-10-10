package main

import (
	"time"

	_ "github.com/zenaton/examples-go" // initialize zenaton client with credentials
	"github.com/zenaton/zenaton-go/integration/go/workflows"
)

func main() {
	workflows.TestEventValueWorkflow.New().Dispatch()

	time.Sleep(2 * time.Second)
	workflows.TestEventValueWorkflow.WhereID("TestEventValueID").Send("MyOtherEvent", nil)

	time.Sleep(2 * time.Second)

	workflows.TestEventValueWorkflow.WhereID("TestEventValueID").Send("MyEvent", "test data")
}
