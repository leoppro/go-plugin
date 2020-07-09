package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/leoppro/go-plugin-demo/hashicorp_plugin_grpc/shared"
	"github.com/leoppro/go-plugin-demo/pkg/sink"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

func main() {
	// We don't want to see the plugin logs.

	pluginPath := "./plugin"
	if len(os.Args) >= 2 {
		pluginPath = os.Args[1]
	}
	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command(pluginPath),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("kv_grpc")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	kv := raw.(shared.KV)
	kv.Put("aa", []byte("bb"))
	v, err := kv.Get("aa")
	fmt.Printf("%s %#v", v, err)

	kv.EmitRow(&sink.RowChangedEvent{StartTs: 1111, Table: &sink.TableName{Table: "123"}})

	log.Info("run benckmark =================")
	benckmark(kv)
	log.Info("finished benckmark ============")
}

func benckmark(s shared.KV) {
	startTime := time.Now()
	for i := int64(0); i < 50_000; i++ {
		s.EmitRow(newRow(i))
	}
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
