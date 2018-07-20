package main

import (
	"os"
	"time"

	"github.com/subosito/gotenv"
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/workflows"
)

func init() {
	gotenv.Load()
	var appID = os.Getenv("ZENATON_APP_ID")
	if appID == "" {
		panic("Please add your Zenaton application id on '.env' file (https://zenaton.com/app/api)")
	}

	var apiToken = os.Getenv("ZENATON_API_TOKEN")
	if apiToken == "" {
		panic("Please add your Zenaton api token on '.env' file (https://zenaton.com/app/api)")
	}

	var appEnv = os.Getenv("ZENATON_APP_ENV")
	if appEnv == "" {
		panic("Please add your Zenaton environment on '.env' file (https://zenaton.com/app/api)")
	}

	// init Zenaton client
	client.InitClient(appID, apiToken, appEnv)
}

func main() {
	workflows.SequentialWorkflow.Dispatch()
	//workflows.AsynchronousWorkflow.Dispatch()
	//workflows.ParallelWorkflow.Dispatch()

	//workflows.EventWorkflow.Dispatch()
	//time.Sleep(2 * time.Second)
	//workflows.EventWorkflow.WhereID("MyId").Send("MyEvent", nil)

	//workflows.WaitWorkflow.Dispatch()

	//workflows.WaitEventWorkflow.Dispatch()
	//time.Sleep(2 * time.Second)
	//workflows.WaitEventWorkflow.WhereID("MyId").Send("MyEvent", nil)

	//recursive.NewRecursiveWorkflow(0, 4).Dispatch()

	//workflows.VersionWorkflow.Dispatch()

	time.Sleep(4 * time.Second)
}

//todo: change parallel to work synchronously per conversation with Gilles
//todo: make sure there are no race conditions in the case of running these things concurrently
//todo: figure out how to make x.New() mandatory (so that you can't just instantiate the workflow/task yourself with a struct literal)
