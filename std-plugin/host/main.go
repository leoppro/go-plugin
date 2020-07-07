package main

import (
	"context"
	"os"
	"plugin"
	"time"

	"github.com/leoppro/go-plugin-demo/pkg/sink"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

func main() {
	pluginPath := "plugin.so"
	if len(os.Args) >= 2 {
		pluginPath = os.Args[1]
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		panic(err)
	}
	newSink, err := p.Lookup("NewSimpleSink")
	if err != nil {
		panic(err)
	}
	newSinkFunc := newSink.(func() sink.Sink)
	s := newSinkFunc()
	ctx := context.Background()
	s.EmitCheckpointTs(ctx, 123)
	s.EmitRowChangedEvents(ctx, &sink.RowChangedEvent{
		CommitTs: 1,
		Table:    &sink.TableName{Table: "123"},
		Columns: map[string]*sink.Column{"col1": &sink.Column{
			Type: 1, Value: "132",
		}},
	})

	log.Info("run benckmark =================")
	benckmark(p)
	log.Info("finished benckmark ============")
}

func benckmark(p *plugin.Plugin) {
	newSink, err := p.Lookup("NewBenchmarkSink")
	if err != nil {
		panic(err)
	}
	newSinkFunc := newSink.(func() sink.Sink)
	s := newSinkFunc()
	ctx := context.Background()
	startTime := time.Now()
	for i := int64(0); i < 50_000_000; i++ {
		s.EmitRowChangedEvents(ctx, newRow(i))
	}
	s.Close()
	log.Info("-", zap.Duration("cost", time.Since(startTime)), zap.Duration("op", time.Since(startTime)/50_000_000))
}

func newRow(mark int64) *sink.RowChangedEvent {
	return &sink.RowChangedEvent{
		RowID:   mark,
		StartTs: uint64(mark),
		Table: &sink.TableName{
			Table: "123",
		},
	}
}
