package common

import (
	"context"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/leoppro/go-plugin-demo/pkg/sink"
)

// Here is an implementation that talks over RPC
type SinkRPC struct{ client *rpc.Client }

func (s *SinkRPC) EmitRowChangedEvents(ctx context.Context, rows ...*sink.RowChangedEvent) error {
	var resp string
	err := s.client.Call("Plugin.EmitRowChangedEvents", rows, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *SinkRPC) EmitDDLEvent(ctx context.Context, ddl *sink.DDLEvent) error {
	var resp string
	err := s.client.Call("Plugin.EmitDDLEvent", ddl, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *SinkRPC) FlushRowChangedEvents(ctx context.Context, resolvedTs uint64) error {
	var resp string
	err := s.client.Call("Plugin.FlushRowChangedEvents", resolvedTs, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *SinkRPC) EmitCheckpointTs(ctx context.Context, ts uint64) error {
	var resp string
	err := s.client.Call("Plugin.EmitCheckpointTs", ts, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (s *SinkRPC) Close() error {
	var resp string
	err := s.client.Call("Plugin.Close", nil, &resp)
	if err != nil {
		return err
	}
	return nil
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type SinkRPCServer struct {
	// This is the real implementation
	Impl sink.Sink
}

func (s *SinkRPCServer) EmitRowChangedEvents(args []*sink.RowChangedEvent, _ *string) error {
	return s.Impl.EmitRowChangedEvents(context.Background(), args...)
}

func (s *SinkRPCServer) EmitDDLEvent(args *sink.DDLEvent, _ *string) error {
	return s.Impl.EmitDDLEvent(context.Background(), args)
}

func (s *SinkRPCServer) FlushRowChangedEvents(args uint64, _ *string) error {
	return s.Impl.FlushRowChangedEvents(context.Background(), args)
}

func (s *SinkRPCServer) EmitCheckpointTs(args uint64, _ *string) error {
	return s.Impl.EmitCheckpointTs(context.Background(), args)
}

func (s *SinkRPCServer) Close(args interface{}, _ *string) error {
	return s.Impl.Close()
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type SinkPlugin struct {
	// Impl Injection
	Impl sink.Sink
}

func (p *SinkPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &SinkRPCServer{Impl: p.Impl}, nil
}

func (SinkPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &SinkRPC{client: c}, nil
}
