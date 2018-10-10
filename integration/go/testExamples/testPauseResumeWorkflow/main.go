package main

import (
	_ "github.com/zenaton/examples-go" // initialize zenaton client with credentials
	"github.com/zenaton/integration/go/workflows"
)

func main() {

	workflows.TestEventValueWorkflow.New().Dispatch()

	queryBuilder, err := workflows.TestEventValueWorkflow.WhereID("TestEventValueID").Pause()

	queryBuilder.Send("EventName", "EventData")

	queryBuilder.Resume()

	//queryBuilder.Send("EventName", "EventData")

	if err != nil {
		panic(err)
	}

}
