package main

import (
	"fmt"
	"sync"
)

type Room struct {
	Mutex   sync.Mutex
	Players []*Player
}

type ChatMessage struct {
	Message string
}

var global_room Room

func (r *Room) AddPlayer(p *Player) {
	r.Mutex.Lock()
	r.Players = append(r.Players, p)
	r.Mutex.Unlock()
}

func (r *Room) RemovePlayer(p *Player) {
	for i, player := range r.Players {
		if player == p {
			r.Mutex.Lock()
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			r.Mutex.Unlock()
			return
		}
	}
}

/// Broadcast sends passed message to each player in the room other than
/// the sender.
func (r *Room) Broadcast(sender *Player, message string) {
	mess := ChatMessage{Message: message}
	r.Mutex.Lock()
	for _, player := range r.Players {
		if player != sender {
			player.SendMessage("chat_message", mess)
		}
	}
	r.Mutex.Unlock()
}

/// chat_message handles receiving a chat message from a client.
/// Broadcasts the message to everyone in the chatroom that that user is in.
/// Does not send message back to the user who sent it.
func (p *Player) chat_message(data []byte) {
	var message ChatMessage
	if err := DecodeJSON(data, &message); err != nil {
		fmt.Println(err)
		p.SendError(err.Error())
	}

	p.Room.Broadcast(p, message.Message)
}
