package main

import (
	"fmt"
	"log"
	"os"

	"net"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/zenaton/zenaton-go/v1/zenaton/job"
)

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	//addr, err := net.ResolveUnixAddr("/var/folders/vm/x61t9b351ps2b94ff7x5g6k40000gn/T/plugin914058232", "bob")

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		//Cmd:             exec.Command("../zenaton-go"),
		Logger:  logger,
		Managed: true,
		Reattach: &plugin.ReattachConfig{
			Protocol: plugin.ProtocolNetRPC,
			Addr:     &net.UnixAddr{Name: "127.0.0.1:10000", Net: "tcp"},
		},
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("job")
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	job := raw.(job.Job)
	fmt.Println(job.Handle())

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

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"job": &job.JobPlugin{},
}
