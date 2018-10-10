package main

import (
	_ "github.com/zenaton/examples-go" // initialize zenaton client with credentials
	"github.com/zenaton/zenaton-go/integration/go/tasks"
	"github.com/zenaton/zenaton-go/integration/go/workflows"
)

func main() {
	//workflows.TestW.WhereID("MyID").Kill()

	workflows.TestW.New(&workflows.Test2{
		TaskReturn: tasks.Test{Print: "field value"},
	}).Dispatch()
}
