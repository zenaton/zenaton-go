package main

import (
	"github.com/zenaton/integration/go/client" // initialize zenaton client with credentials
	"github.com/zenaton/integration/go/workflows"
)

func init() { client.SetEnv("testReturnError.env") }
func main() {

	workflows.TestW.New(workflows.Test2{
		TaskError: "testTaskError",
	}).Dispatch()
}
