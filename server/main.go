package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	msgch := make(chan string, 10)
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, msgch)
		go readConsole(msgch)
	}
}

func handleConn(c net.Conn, msgch chan string) {
	defer c.Close()

	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n\r"))
		if err != nil {
			return
		}

		select {
		case msg := <-msgch:
			_, err := io.WriteString(c, msg)
			if err != nil {
				fmt.Println("Error on sending: ", err)
			}
			fmt.Println("Sent: ", msg)
		default:
		}

		time.Sleep(1 * time.Second)
	}
}

func readConsole(msgch chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error on reading from console: ", err)
			continue
		}

		if len(text) > 0 {
			msgch <- text
		}
	}
}
