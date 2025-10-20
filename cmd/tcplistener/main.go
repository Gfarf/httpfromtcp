package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Gfarf/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("===== Loop will start ======")
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("===== Connection Accepted ======")
		lineChannels, err := request.RequestFromReader(connection)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", lineChannels.RequestLine.Method, lineChannels.RequestLine.RequestTarget, lineChannels.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range lineChannels.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
		fmt.Println("===== Connection And Channel Closed ======")
	}
}

/*func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("===== Loop will start ======")
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("===== Connection Accepted ======")
		lineChannels := getLinesChannel(connection)
		for line := range lineChannels {
			fmt.Println(line)
		}
		fmt.Println("===== Connection And Channel Closed ======")
	}

}*/
