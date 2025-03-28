package main

import (
	pb "chat/proto"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var secret string

type server struct {
	pb.UnimplementedChatServer
	connections []*connection
	names       map[string]bool
}

type connection struct {
	stream pb.Chat_ConnectServer
	name   string
	index  int
}

func (s *server) ClaimName(ctx context.Context, in *pb.ClaimNameRequest) (*pb.ClaimNameResponse, error) {
	if occupied, ok := s.names[in.Name]; ok && occupied {
		return nil, status.Error(codes.AlreadyExists, "name already claimed")
	}

	s.names[in.Name] = false

	token, err := createTokenForName(in.Name)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid name")
	}

	return &pb.ClaimNameResponse{Token: token}, nil
}

func createTokenForName(name string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": name,
		"iss": "winnie.at",
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func (s *server) Connect(stream pb.Chat_ConnectServer) error {
	name, err := getName(stream)

	if err != nil {
		return err
	}

	if s.names[name] {
		return status.Error(codes.AlreadyExists, "This Name is already claimed and used please choose another name or wait for the name to be released!")
	}

	s.names[name] = true

	conn := connection{stream: stream, name: name, index: len(s.connections)}
	s.connections = append(s.connections, &conn)
	return s.broadcastMessages(&conn)
}

func getName(stream pb.Chat_ConnectServer) (string, error) {
	raw_token, err := getTokenString(stream)
	if err != nil {
		return "", err
	}

	split_string := strings.Split(raw_token, " ")

	if len(split_string) != 2 || split_string[0] != "Bearer" {
		return "", status.Error(codes.Unauthenticated, "invalid token format")
	}

	token, err := jwt.Parse(split_string[1], func(token *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil && token.Valid {
		return "", err
	}

	err = validateToken(token)

	if err != nil {
		return "", err
	}

	return token.Claims.GetSubject()
}

func validateToken(token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.Claims)
	if !ok {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	err := jwt.NewValidator().Validate(claims)

	if err != nil {
		return err
	}

	iss, err := claims.GetIssuer()

	if err != nil || iss != "winnie.at" {
		return err
	}

	return nil
}

func getTokenString(stream pb.Chat_ConnectServer) (string, error) {
	metadata, _ := metadata.FromIncomingContext(stream.Context())
	raw_token := metadata.Get("authorization")

	if raw_token == nil || len(raw_token) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization token")
	}

	return raw_token[0], nil
}

func (s *server) broadcastMessages(conn *connection) error {
	for {
		in, err := conn.stream.Recv()
		if err != nil {
			s.reactToError(err, conn)
			return err
		}

		s.broadcastMessage(in.Message, conn)

		fmt.Printf("%s sent message: %s\n", conn.name, in.Message)
	}
}

func (s *server) reactToError(err error, conn *connection) {
	s.connections = slices.Delete(s.connections, conn.index, conn.index+1)

	for _, c := range s.connections[conn.index:] {
		c.index--
	}

	s.names[conn.name] = false

	if err, ok := status.FromError(err); ok {
		if err.Code() == codes.Canceled {
			return
		}
	}

	if err == io.EOF {
		return
	}

	fmt.Printf("%s encountered error: %v\n", conn.name, err)
}

func (s *server) broadcastMessage(message string, senderConnection *connection) {
	for _, conn := range s.connections {
		if conn.stream != senderConnection.stream {
			_ = conn.stream.Send(&pb.IncomingMessage{Name: senderConnection.name, Response: message})
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	secret = os.Getenv("SECRET")

	lis, err := net.Listen("tcp4", ":5555")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	grpcServer := &server{names: map[string]bool{}}

	pb.RegisterChatServer(s, grpcServer)

	fmt.Printf("Server listening on %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
