package main

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/leoppro/go-plugin-demo/hashicorp_plugin/common"
	"github.com/leoppro/go-plugin-demo/pkg/sink"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	pluginPath := "./plugin"
	if len(os.Args) >= 2 {
		pluginPath = os.Args[1]
	}

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(pluginPath),
		Logger:          logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal("-", zap.Error(err))
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("simple_sink")
	if err != nil {
		log.Fatal("-", zap.Error(err))
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	s := raw.(sink.Sink)
	err = s.EmitCheckpointTs(context.Background(), 123)
	log.Error("-", zap.Error(err))
	err = s.EmitRowChangedEvents(context.Background(), &sink.RowChangedEvent{RowID: 1})
	log.Error("-", zap.Error(err))

	log.Info("run benckmark =================")
	raw, err = rpcClient.Dispense("benchmark_sink")
	if err != nil {
		log.Fatal("-", zap.Error(err))
	}
	s = raw.(sink.Sink)
	benckmark(s)
	log.Info("finished benckmark ============")
}

func benckmark(s sink.Sink) {
	ctx := context.Background()
	startTime := time.Now()
	for i := int64(0); i < 50_000; i++ {
		s.EmitRowChangedEvents(ctx, newRow(i))
	}
	s.Close()
	log.Info("-", zap.Duration("cost", time.Since(startTime)), zap.Duration("op", time.Since(startTime)/50_000))
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
	"simple_sink":    &common.SinkPlugin{},
	"benchmark_sink": &common.SinkPlugin{},
}
