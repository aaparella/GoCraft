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

func (p *Player) handle_auth(data []byte) {
	var message AuthHello
	if err := decodeJSON(data, &message); err != nil {
		fmt.Println(err)
	}

	// Log in the user
	response := AuthResponse{
		LobbyID: 1,
		Success: true,
	}
	p.sendMessage("auth_response", response)
}
