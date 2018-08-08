package service

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

import (
	"errors"
	"fmt"
	"log"
	"net"

	pb "cuttleserver/common/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("say hello")

	sess, err := SessionMgr.GetSession(in.AccountID)
	if err != nil {
		sess, err = NewSession(in.AccountID)
		if err != nil {
			fmt.Println("NewSession() error")
			return nil, errors.New("say hello error")
		}
		SessionMgr.AddSession(sess)
	}

	respID, byteMessage := sess.Handler(uint16(in.MessageID), in.Content)

	return &pb.HelloReply{AccountID: sess.AccountID, MessageID: uint32(respID), Content: byteMessage}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("say hello again")
	return &pb.HelloReply{AccountID: in.AccountID, MessageID: 2, Content: []byte("Hello again")}, nil
}

func StartGRPC() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
