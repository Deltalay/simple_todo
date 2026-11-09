package main

import (
	authpb "auth_service/proto"
	"context"
	"errors"
	"log"
	"net"

	"google.golang.org/grpc"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	return nil, errors.New("invalid credentials")
}
func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	username := req.GetUsername()
	return nil, errors.New("invalid credentials")
}
func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, &AuthServer{})
	log.Println("Auth server running on port 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
