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

/// SendMessage sends passed message to client.
/// Take type string, and the data object as arguments.
func (p *Player) SendMessage(typ string, data interface{}) {
	response := Response{
		Type: typ,
		Data: data,
	}
	json, _ := json.Marshal(response)
	fmt.Fprintf(p.Conn, "%s\n", json)
}

/// SendError sends an error string to client in proper format.
/// Takes string instead of error object for when we send errors that
/// are not embodied in actual error objects.
func (p *Player) SendError(e string) {
	err := Error{
		Error: e,
	}
	p.SendMessage("error", err)
}

/// DecodeJSON unmarshals data into structure.
/// Argument structure needs to be a pointer to the desired structure.
func DecodeJSON(data []byte, structure interface{}) error {
	reader := json.NewDecoder(bytes.NewBuffer(data))
	return reader.Decode(structure)
}

/// Player model.
/// Includes connection for reading / writing.
/// Necessary state information included as well.
type Player struct {
	Conn     net.Conn
	Username string
	State    int // Can change to enum later
	Room     *Room
}

/// Handle.
/// Get the type of the message and switch based on that type.
/// Default case is sending an unsupported command error.
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

/// Listen listens for each message received, and calls handler.
/// Also does cleanup and teardown for Player struct.
/// Runs in own goroutine.
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
