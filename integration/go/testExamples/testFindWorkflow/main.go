package main

import (
	"github.com/zenaton/zenaton-go/integration/go/client" // initialize zenaton client with credentials
	"github.com/zenaton/zenaton-go/integration/go/workflows"
	"time"
)

func init() { client.SetEnv("dev-testFindWorkflow") }
func main() {

	wf, err := workflows.TestW.WhereID("testFind").Find()
	if err != nil {
		panic(err)
	}

	if wf != nil {
		_, err := workflows.TestW.WhereID("testFind").Kill()
		if err != nil {
			panic(err)
		}
	}

	workflows.TestW.New(workflows.Test2{
		IDstr:      "testFind",
		TaskReturn: "testTaskReturn",
	}).Dispatch()

	wf, err = workflows.TestW.WhereID("testFind").Find()

	if err != nil {
		panic(err)
	}

	instance := wf.GetData().(*workflows.Test2)

	//to make the logs more predictable between dispatches
	time.Sleep(12 * time.Second)

	// launch a new instance with same return value
	workflows.TestW.New(workflows.Test2{
		IDstr:      "testFind",
		TaskReturn: instance.TaskReturn,
	}).Dispatch()
}
