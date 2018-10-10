package main

import (
	"time"

	_ "github.com/zenaton/integration/go/client" // initialize zenaton client with credentials
	"github.com/zenaton/integration/go/workflows"
)

func main() {

	workflows.TestW.New(workflows.Test2{
		IDstr:      "testFind",
		TaskReturn: "testTaskReturn",
	}).Dispatch()

	wf, err := workflows.TestW.WhereID("testFind").Find()

	if err != nil {
		panic(err)
	}

	instance := wf.GetData().(*workflows.Test2)

	//to make the logs more predictable between dispatches
	time.Sleep(5 * time.Second)

	// launch a new instance with same return value
	workflows.TestW.New(workflows.Test2{
		IDstr:      "testFind",
		TaskReturn: instance.TaskReturn,
	}).Dispatch()
}
