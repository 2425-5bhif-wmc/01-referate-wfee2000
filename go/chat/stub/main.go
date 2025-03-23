package main

import (
	"bufio"
	pb "chat/proto"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var message string

type state struct {
	termios unix.Termios
}

func main() {
	conn, err := grpc.NewClient("winnie.at:5555", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to server\n")
	defer conn.Close()

	client := pb.NewChatClient(conn)
	response, err := client.ClaimName(context.Background(), &pb.ClaimNameRequest{Name: "Winnie"})

	if err != nil {
		panic(err)
	}

	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", response.Token))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := client.Connect(ctx)

	if err != nil {
		panic(err)
	}

	go displayMessages(stream)

	sendMessage(stream)
}

func displayMessages(stream pb.Chat_ConnectClient) {
	for {
		in, err := stream.Recv()
		if err != nil {
			panic(err)
		}

		fmt.Printf("\033[2K\r%s: %v\n\rWrite message: %s", in.Name, in.Response, message)
	}
}

func sendMessage(stream pb.Chat_ConnectClient) {
	state, err := disableCanonical(os.Stdin.Fd())

	if err != nil {
		fmt.Println(err)
		return
	}

	defer setState(state)

	reader := bufio.NewReader(os.Stdin)

	for {
		print("Write message: ")

		for {
			rune, err := reader.ReadByte()
			if err != nil {
				fmt.Println(err)
				return
			}

			if rune == '\r' || rune == '\n' {
				break
			}

			message += string(rune)
		}

		err := stream.Send(&pb.OutgoingMessage{Message: strings.TrimSpace(message)})

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("\033[A\033[2K\rYou wrote: %s\n\r", message)

		message = ""
	}
}

func disableCanonical(fd uintptr) (*unix.Termios, error) {
	attributes := unix.Termios{}

	if err := termios.Tcgetattr(fd, &attributes); err != nil {
		return nil, err
	}

	oldState := attributes

	attributes.Lflag &^= unix.ICANON
	if err := termios.Tcsetattr(fd, termios.TCSANOW, &attributes); err != nil {
		return nil, err
	}

	return &oldState, nil
}

func setState(state *unix.Termios) error {
	return termios.Tcsetattr(os.Stdin.Fd(), termios.TCSANOW, state)
}
