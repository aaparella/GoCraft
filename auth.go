package main

import "fmt"

type AuthHello struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type AuthResponse struct {
	LobbyID int  `json:"lobby_id"`
	Success bool `json:"success"`
}

func (p *Player) auth_hello(data []byte) {
	var message AuthHello
	if err := DecodeJSON(data, &message); err != nil {
		fmt.Println(err)
		p.SendError(err.Error())
		return
	}
	// TODO(rweichler) : Hit the API lololol
	global_room.AddPlayer(p)
	p.Room = &global_room

	response := AuthResponse{
		LobbyID: 1,
		Success: true,
	}
	p.SendMessage("auth_response", response)
}
