package main

import (
	"github.com/zenaton/zenaton-go/integration/go/client" // initialize zenaton client with credentials
	"github.com/zenaton/zenaton-go/integration/go/workflows"
)

func init() { client.SetEnv("testReturnFromTaskInsideTask.env") }
func main() {
	workflows.TestRelaunchTaskWorkflow.New().Dispatch()
}
