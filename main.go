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
	workflows.SequentialWorkflow.Execute()
	workflows.AsynchronousWorkflow.Execute()
	//workflows.ParallelWorkflow.Dispatch()
	//
	//workflows.EventWorkflow.Dispatch()
	//time.Sleep(2000 * time.Millisecond)
	//workflows.EventWorkflow.WhereID("MyId").Send("MyEvent", nil);
	time.Sleep(30 * time.Millisecond)
}
