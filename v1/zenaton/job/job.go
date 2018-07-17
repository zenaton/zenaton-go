package job

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// todo: maybe I don't want this to be exported, so only tasks and workflows can implement this interface

// Job is the interface that we're exposing as a plugin.
type Job interface {
	Handle() interface{}
	AsyncHandle(chan interface{})
}

// Here is an implementation that talks over RPC
type JobRPC struct{ client *rpc.Client }

func (g *JobRPC) Handle() interface{} {
	var resp interface{}
	err := g.client.Call("Plugin.Handle", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

func (g *JobRPC) AsyncHandle(channel chan interface{}) {
	var resp interface{}
	err := g.client.Call("Plugin.AsyncHandle", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}
	channel <- resp
}

// Here is the RPC server that JobRPC talks to, conforming to
// the requirements of net/rpc
type JobRPCServer struct {
	// This is the real implementation
	Impl Job
}

func (s *JobRPCServer) Handle(args interface{}, resp *interface{}) error {
	*resp = s.Impl.Handle()
	return nil
}

func (s *JobRPCServer) AsyncHandle(args interface{}, resp *interface{}) error {
	*resp = s.Impl.Handle()
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a JobRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return JobRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type JobPlugin struct {
	// Impl Injection
	Impl Job
}

func (p *JobPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &JobRPCServer{Impl: p.Impl}, nil
}

func (JobPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &JobRPC{client: c}, nil
}
