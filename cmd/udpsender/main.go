package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn, err := net.DialUDP("udp", nil, listener)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	reading := bufio.NewReader(os.Stdin)
	fmt.Println("===== Loop will start ======")
	for {
		fmt.Print("> ")
		data, err := reading.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Message sent: %s", data)
	}
}
