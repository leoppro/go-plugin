package main

import (
	"context"

	"log"

	"github.com/hashicorp/go-plugin"
	"github.com/leoppro/go-plugin-demo/hashicorp_plugin/common"
	"github.com/leoppro/go-plugin-demo/pkg/sink"
	"go.uber.org/zap"
)

type SimpleSink struct{}

const PLUGIN_NAME = "HASHICORP_PLUGIN"

func NewSimpleSink() sink.Sink {
	return &SimpleSink{}
}

func (s *SimpleSink) EmitRowChangedEvents(ctx context.Context, rows ...*sink.RowChangedEvent) error {
	log.Print("EmitRowChangedEvents", zap.String("plugin", PLUGIN_NAME), zap.Reflect("rows", rows))
	return nil
}

func (s *SimpleSink) EmitDDLEvent(ctx context.Context, ddl *sink.DDLEvent) error {
	log.Print("EmitDDLEvent", zap.String("plugin", PLUGIN_NAME), zap.Reflect("ddl", ddl))
	return nil
}

func (s *SimpleSink) FlushRowChangedEvents(ctx context.Context, resolvedTs uint64) error {
	log.Print("FlushRowChangedEvents", zap.String("plugin", PLUGIN_NAME), zap.Reflect("resolvedTs", resolvedTs))
	return nil
}

func (s *SimpleSink) EmitCheckpointTs(ctx context.Context, ts uint64) error {
	log.Print("EmitCheckpointTs", zap.String("plugin", PLUGIN_NAME), zap.Reflect("ts", ts))
	return nil
}

func (s *SimpleSink) Close() error {
	log.Print("Close", zap.String("plugin", PLUGIN_NAME))
	return nil
}

type BenchmarkSink struct {
	count uint64
}

func NewBenchmarkSink() sink.Sink {
	return &BenchmarkSink{}
}

func (s *BenchmarkSink) EmitRowChangedEvents(ctx context.Context, rows ...*sink.RowChangedEvent) error {
	for range rows {
		s.count++
	}
	return nil
}

func (s *BenchmarkSink) EmitDDLEvent(ctx context.Context, ddl *sink.DDLEvent) error {
	s.count++
	return nil
}

func (s *BenchmarkSink) FlushRowChangedEvents(ctx context.Context, resolvedTs uint64) error {
	return nil
}

func (s *BenchmarkSink) EmitCheckpointTs(ctx context.Context, ts uint64) error {
	s.count++
	return nil
}

func (s *BenchmarkSink) Close() error {
	log.Print("Close", zap.String("plugin", PLUGIN_NAME), zap.Uint64("count", s.count))
	return nil
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
	simpleSink := &SimpleSink{}
	benchmarkSink := &BenchmarkSink{}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"simple_sink":    &common.SinkPlugin{Impl: simpleSink},
		"benchmark_sink": &common.SinkPlugin{Impl: benchmarkSink},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
