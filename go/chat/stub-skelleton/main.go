package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var message string

type state struct {
	termios unix.Termios
}

func main() {
	// TODO: generate protobuf files

	// Open unsafe gRPC connection to server
	conn, err := grpc.NewClient("winnie.at:80", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connected to server\n")

	// Close connection when done
	defer conn.Close()

	// TODO: Create Stub on connection
	// TODO: Claim Name from server

	// Append jwt to metadata
	// md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", response.Token))
	// ctx := metadata.NewOutgoingContext(context.Background(), md)

	// TODO: open bidirectional stream to server

	// TODO: Asynchronously listen for messages from server

	// TODO: Send messages to server
}

// Listens to stream and prints received messages
// without overriding terminal.
func displayMessages() { //stream pb.Chat_ConnectClient) {
	for {
		// TODO: Wait for message from server
		// Print message
		// fmt.Printf("\033[2K\r%s: %v\n\rWrite message: %s", in.Name, in.Response, message)
	}
}

// Reads input from terminal and sends read
// Messages to server.
func sendMessages() { // stream pb.Chat_ConnectClient) {
	// Remove canonical flag from terminal
	state, err := disableCanonical(os.Stdin.Fd())

	if err != nil {
		fmt.Println(err)
		return
	}

	// Set terminal state back to original state
	defer setState(state)

	// Open input stream to terminal
	_ = bufio.NewReader(os.Stdin)

	for {
		// TODO: Get message from terminal

		// TODO: Send message to server

		// Clear message
		message = ""
	}
}

// Reads input from terminal and manipulates message string accordingly.
// Special characters like \n \r and 127 should be handled.
// Returns message string and error if any.
func readMessage(reader *bufio.Reader) (string, error) {
	fmt.Print("Write message: ")

	// message := ""

	for {
		// TODO: read first character

		// TODO: return if enter or newline

		// remove character from message
		/* if rune == 127 {
		fmt.Print("\b\b  \b\b")
		if len(message) > 1 {
			message = message[:len(message)-1]
			fmt.Print("\b \b")
		} else if len(message) == 1 {
			fmt.Print("\r\033[2K\rWrite message: ")
		}

		continue
		}*/

		// TODO: append character to message
	}
}

// Prints message and sends it to server.
func sendMessage() { // stream pb.Chat_ConnectClient, message string) {
	// TODO: send message

	// Print new message
	// fmt.Printf("\033[A\033[2K\rYou wrote: %s\n\r", message)
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
