package rpc

import (
	"context"
	"database/sql"
	"net"

	pb "github.com/nilesh/shopcore-contracts/shopcorepb/shopcorepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}

func (s *Server) GetUserDetails(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	var name, email string
	err := s.db.QueryRowContext(ctx, `SELECT name,email FROM users WHERE id=$1 AND deleted_at IS NULL`, req.GetUserId()).Scan(&name, &email)
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "could not load user")
	}
	return &pb.UserResponse{UserName: name, Email: email}, nil
}

func Listen(address string, db *sql.DB) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, NewServer(db))
	return server.Serve(listener)
}
