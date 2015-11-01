package main

import "fmt"

type Room struct {
	Players []*Player
}

type ChatMessage struct {
	Message string
}

var global_room Room

func (r *Room) AddPlayer(p *Player) {
	r.Players = append(r.Players, p)
}

func (r *Room) RemovePlayer(p *Player) {
	for i, player := range r.Players {
		if player == p {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			return
		}
	}
}

func (r *Room) Broadcast(sender *Player, message string) {
	mess := ChatMessage{Message: message}
	for _, player := range r.Players {
		if player != sender {
			player.SendMessage("chat_message", mess)
		}
	}
}

func (p *Player) chat_message(data []byte) {
	var message ChatMessage
	if err := DecodeJSON(data, &message); err != nil {
		fmt.Println(err)
		p.SendError(err.Error())
	}

	p.Room.Broadcast(p, message.Message)
}
