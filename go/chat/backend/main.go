package main

import (
	pb "chat/proto"
	"context"
	"fmt"
	"io"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
	connections []*connection
}

type connection struct {
	stream pb.ChatService_ConnectServer
	name   string
	index  int
}

func (s *server) ClaimName(ctx context.Context, in *pb.ClaimNameRequest) (*pb.ClaimNameResponse, error) {
	return &pb.ClaimNameResponse{Token: "string"}, nil
}

func (s *server) Connect(stream pb.ChatService_ConnectServer) error {
	// TODO: check header for token

	conn := connection{stream: stream, name: "TODO", index: len(s.connections)}
	s.connections = append(s.connections, &conn)
	return s.broadcastMessages(conn)
}

func (s *server) broadcastMessages(conn connection) error {
	for {
		in, err := conn.stream.Recv()
		if err == io.EOF {
			s.connections = append(s.connections[:conn.index], s.connections[conn.index+1:]...)
			break
		}
		if err != nil {
			fmt.Println(err)
			return err
		}

		s.broadcastMessage(in.Message, conn)

		fmt.Printf("Received message: %v\n", in)
	}

	return nil
}

func (s *server) broadcastMessage(message string, senderConnection connection) {
	for _, conn := range s.connections {
		if conn.stream != senderConnection.stream {
			_ = conn.stream.Send(&pb.IncomingMessage{Name: conn.name, Response: message})
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":5555")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	grpcServer := &server{}

	pb.RegisterChatServiceServer(s, grpcServer)

	fmt.Printf("Server listening on %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
