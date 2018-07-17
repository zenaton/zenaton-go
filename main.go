package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/subosito/gotenv"
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
	"github.com/zenaton/zenaton-go/v1/zenaton/job"
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

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	//workflows.SequentialWorkflow.Execute()
	//workflows.AsynchronousWorkflow.Execute()
	//workflows.ParallelWorkflow.Execute()
	////
	////workflows.EventWorkflow.Execute()
	////time.Sleep(2000 * time.Millisecond)
	////workflows.EventWorkflow.WhereID("MyId").Send("MyEvent", nil)
	////workflows.WaitWorkflow.Dispatch()
	//
	//workflows.WaitEventWorkflow.Dispatch()
	//
	//time.Sleep(2 * time.Second)
	//workflows.WaitEventWorkflow.WhereID("MyId").Send("MyEvent", nil)

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"job": &job.JobPlugin{Impl: &workflows.SequentialWorkflow},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})

}

//todo: change parallel to work synchronously per conversation with Gilles
