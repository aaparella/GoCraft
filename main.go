package main

import "net"

func main() {
	listener, _ := net.Listen("tcp", ":8007")
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		player := Player{Conn: conn, State: 0}
		go player.listen()
	}
}
