package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

type Response struct {
	Type string
	Data interface{}
}

type Message struct {
	Type string
	Data json.RawMessage
}

type Error struct {
	Error string
}

func (p *Player) SendMessage(typ string, data interface{}) {
	response := Response{
		Type: typ,
		Data: data,
	}
	json, _ := json.Marshal(response)
	fmt.Fprintf(p.Conn, "%s\n", json)
}

func (p *Player) SendError(e string) {
	err := Error{
		Error: e,
	}
	p.SendMessage("error", err)
}

func DecodeJSON(data []byte, structure interface{}) error {
	reader := json.NewDecoder(bytes.NewBuffer(data))
	return reader.Decode(structure)
}

type Player struct {
	Conn     net.Conn
	Username string
	State    int // Can change to enum later
	Room     *Room
}

func (p *Player) Handle(data []byte) {
	reader := json.NewDecoder(bytes.NewBuffer(data))
	var mess Message
	if err := reader.Decode(&mess); err != nil {
		p.SendError(err.Error())
	} else {
		// Valid message or so we think
		switch mess.Type {
		case "auth_hello":
			p.auth_hello(mess.Data)
		case "chat_message":
			p.chat_message(mess.Data)
		default:
			p.SendError("Unsupported command")
		}
	}
}

// Runs in own goroutine
func (p *Player) Listen() {
	reader := bufio.NewReader(p.Conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			global_room.RemovePlayer(p)
			p.Room = nil
			return
		}
		p.Handle(line)
	}
}
