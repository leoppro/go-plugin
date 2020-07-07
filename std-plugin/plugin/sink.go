package main

import (
	"context"

	"github.com/leoppro/go-plugin-demo/pkg/sink"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

type SimpleSink struct{}

const PLUGIN_NAME = "STD_PLUGIN"

func NewSimpleSink() sink.Sink {
	return &SimpleSink{}
}

func (s *SimpleSink) EmitRowChangedEvents(ctx context.Context, rows ...*sink.RowChangedEvent) error {
	log.Info("EmitRowChangedEvents", zap.String("plugin", PLUGIN_NAME), zap.Reflect("rows", rows))
	return nil
}

func (s *SimpleSink) EmitDDLEvent(ctx context.Context, ddl *sink.DDLEvent) error {
	log.Info("EmitDDLEvent", zap.String("plugin", PLUGIN_NAME), zap.Reflect("ddl", ddl))
	return nil
}

func (s *SimpleSink) FlushRowChangedEvents(ctx context.Context, resolvedTs uint64) error {
	log.Info("FlushRowChangedEvents", zap.String("plugin", PLUGIN_NAME), zap.Reflect("resolvedTs", resolvedTs))
	return nil
}

func (s *SimpleSink) EmitCheckpointTs(ctx context.Context, ts uint64) error {
	log.Info("EmitCheckpointTs", zap.String("plugin", PLUGIN_NAME), zap.Reflect("ts", ts))
	return nil
}

func (s *SimpleSink) Close() error {
	log.Info("Close", zap.String("plugin", PLUGIN_NAME))
	return nil
}


type BenchmarkSink struct{
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
	log.Info("Close", zap.String("plugin", PLUGIN_NAME),zap.Uint64("count", s.count))
	return nil
}