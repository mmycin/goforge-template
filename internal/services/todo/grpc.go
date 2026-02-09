// Generated at 2026-02-09T18:08:46+06:00
package todo

import (
	"context"

	"github.com/mmycin/goforge/internal/server"
	pb "github.com/mmycin/goforge/proto/todo/gen"
	"google.golang.org/grpc"
)

type TodoGRPC struct {
	pb.UnimplementedTodoServiceServer
}

func init() {
	server.RegisterGRPC(func(s *grpc.Server) {
		pb.RegisterTodoServiceServer(s, &TodoGRPC{})
	})
}


func (t *TodoGRPC) CreateTodo(ctx context.Context, in *pb.CreateTodoRequest) (*pb.Todo, error) {
	return nil, nil
}

func (t *TodoGRPC) GetTodo(ctx context.Context, in *pb.GetTodoRequest) (*pb.Todo, error) {
	return nil, nil
}

func (t *TodoGRPC) ListTodos(ctx context.Context, in *pb.ListTodoRequest) (*pb.ListTodoResponse, error) {
	return nil, nil
}

func (t *TodoGRPC) UpdateTodo(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.Todo, error) {
	return nil, nil
}

func (t *TodoGRPC) DeleteTodo(ctx context.Context, in *pb.DeleteTodoRequest) (*pb.DeleteTodoResponse, error) {
	return nil, nil
}

