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
	// Open unsafe gRPC connection to server
	conn, err := grpc.NewClient("winnie.at:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to server\n")

	// Close connection when done
	defer conn.Close()

	// Create Stub on connection
	client := pb.NewChatClient(conn)
	// Claim Name from server
	response, err := client.ClaimName(context.Background(), &pb.ClaimNameRequest{Name: strings.Join(os.Args[1:], " ")})

	if err != nil {
		panic(err)
	}

	// Append jwt to metadata
	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", response.Token))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	// open bidirectional stream to server
	stream, err := client.Connect(ctx)

	if err != nil {
		panic(err)
	}

	// Asynchronously listen for messages from server
	go displayMessages(stream)

	// Send messages to server
	sendMessages(stream)
}

// Listens to stream and prints received messages
// without overriding terminal.
func displayMessages(stream pb.Chat_ConnectClient) {
	for {
		// Wait for message from server
		in, err := stream.Recv()
		if err != nil {
			panic(err)
		}

		// Print message
		fmt.Printf("\033[2K\r%s: %v\n\rWrite message: %s", in.Name, in.Response, message)
	}
}

// Reads input from terminal and sends read
// Messages to server.
func sendMessages(stream pb.Chat_ConnectClient) {
	// Remove canonical flag from terminal
	state, err := disableCanonical(os.Stdin.Fd())

	if err != nil {
		fmt.Println(err)
		return
	}

	// Set terminal state back to original state
	defer setState(state)

	// Open input stream to terminal
	reader := bufio.NewReader(os.Stdin)

	for {
		// Get message from terminal
		err := readMessage(reader)

		if err != nil {
			fmt.Println(err)
			return
		}
		// Send message to server
		sendMessage(stream, strings.TrimSpace(message))

		// Clear message
		message = ""
	}
}

// Reads input from terminal and manipulates message string accordingly.
// Special characters like \n \r and 127 should be handled.
// Returns message string and error if any.
func readMessage(reader *bufio.Reader) error {
	fmt.Print("Write message: ")

	for {
		// read first character
		rune, err := reader.ReadByte()

		if err != nil {
			return err
		}

		// return if enter or newline
		if rune == '\r' || rune == '\n' {
			return nil
		}

		// remove character from message
		if rune == 127 {
			fmt.Print("\b\b  \b\b")
			if len(message) > 1 {
				message = message[:len(message)-1]
				fmt.Print("\b \b")
			} else if len(message) == 1 {
				fmt.Print("\r\033[2K\rWrite message: ")
			}

			continue
		}

		// append character to message
		message += string(rune)
	}
}

// Prints message and sends it to server.
func sendMessage(stream pb.Chat_ConnectClient, message string) {
	// send message
	err := stream.Send(&pb.OutgoingMessage{Message: strings.TrimSpace(message)})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("\033[A\033[2K\rYou wrote: %s\n\r", message)
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
