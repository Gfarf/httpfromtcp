package main

import (
	"fmt"
	"net"
	"os"
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
		lineChannels, err := RequestFromReader(connection)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for line := range lineChannels {
			fmt.Println(line)
		}
		fmt.Println("===== Connection And Channel Closed ======")
	}
}
