package shared

import (
	"github.com/leoppro/go-plugin-demo/hashicorp_plugin_grpc/proto"
	"github.com/leoppro/go-plugin-demo/pkg/sink"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.KVClient }

func (m *GRPCClient) Put(key string, value []byte) error {
	_, err := m.client.Put(context.Background(), &proto.PutRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (m *GRPCClient) Get(key string) ([]byte, error) {
	resp, err := m.client.Get(context.Background(), &proto.GetRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

func (m *GRPCClient) EmitRow(row *sink.RowChangedEvent) error {
	_, err := m.client.EmitRow(context.Background(), &proto.RowChangedEvent{
		StartTs:   row.StartTs,
		CommitTs:  row.CommitTs,
		TableName: row.Table.Table,
	})
	if err != nil {
		return err
	}
	return nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl KV
	proto.UnimplementedKVServer
}

func (m *GRPCServer) Put(
	ctx context.Context,
	req *proto.PutRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Put(req.Key, req.Value)
}

func (m *GRPCServer) Get(
	ctx context.Context,
	req *proto.GetRequest) (*proto.GetResponse, error) {
	v, err := m.Impl.Get(req.Key)
	return &proto.GetResponse{Value: v}, err
}

func (m *GRPCServer) EmitRow(ctx context.Context, in *proto.RowChangedEvent) (*proto.Empty, error) {
	err := m.Impl.EmitRow(&sink.RowChangedEvent{StartTs: in.StartTs})
	return &proto.Empty{}, err
}
