package shared

import (
	"net/rpc"

	"github.com/leoppro/go-plugin-demo/pkg/sink"
)

// RPCClient is an implementation of KV that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) Put(key string, value []byte) error {
	// We don't expect a response, so we can just use interface{}
	var resp interface{}

	// The args are just going to be a map. A struct could be better.
	return m.client.Call("Plugin.Put", map[string]interface{}{
		"key":   key,
		"value": value,
	}, &resp)
}

func (m *RPCClient) Get(key string) ([]byte, error) {
	var resp []byte
	err := m.client.Call("Plugin.Get", key, &resp)
	return resp, err
}

func (m *RPCClient) EmitRow(row *sink.RowChangedEvent) error {
	var resp interface{}
	err := m.client.Call("Plugin.EmitRow", row, &resp)
	return err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl KV
}

func (m *RPCServer) Put(args map[string]interface{}, resp *interface{}) error {
	return m.Impl.Put(args["key"].(string), args["value"].([]byte))
}

func (m *RPCServer) Get(key string, resp *[]byte) error {
	v, err := m.Impl.Get(key)
	*resp = v
	return err
}

func (m *RPCServer) EmitRow(row *sink.RowChangedEvent) error {
	return m.Impl.EmitRow(row)
}
