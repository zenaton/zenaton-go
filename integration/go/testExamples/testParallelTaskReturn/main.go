package main

import (
	_ "github.com/zenaton/examples-go" // initialize zenaton client with credentials
	"github.com/zenaton/integration/go/workflows"
)

func main() {
	workflows.TestW.WhereID("MyID").Kill()

	//workflows.TestW.NewInstance(&workflows.Test{
	//	Parallel: true,
	//	Return:   "testReturn",
	//}).Dispatch()
}
