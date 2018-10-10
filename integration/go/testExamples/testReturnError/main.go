package main

import (
	"github.com/zenaton/zenaton-go/integration/go/client" // initialize zenaton client with credentials
	"github.com/zenaton/zenaton-go/integration/go/workflows"
)

func init() { client.SetEnv("testReturnError.env") }
func main() {

	workflows.TestW.New(workflows.Test2{
		TaskError: "testTaskError",
	}).Dispatch()
}
