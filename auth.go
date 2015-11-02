package main

import (
	"fmt"
	"net/http"
)

type AuthHello struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type AuthResponse struct {
	LobbyID int  `json:"lobby_id"`
	Success bool `json:"success"`
}

/// Perform the authentication check
/// Hits API to check if the credentials are valid
/// Assumes any response other than a 200 is an indication
/// that the crdentials provided were invalid
func (p *Player) auth_hello(data []byte) {
	var message AuthHello
	if err := DecodeJSON(data, &message); err != nil {
		fmt.Println(err)
		p.SendError(err.Error())
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8080/login", nil)
	req.Header.Set("username", message.User)
	req.Header.Set("password", message.Pass)
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
		p.SendError("Could not authenticate")
		return
	}

	global_room.AddPlayer(p)
	p.Room = &global_room
	response := AuthResponse{
		LobbyID: 1,
		Success: true,
	}
	p.SendMessage("auth_response", response)
	resp.Body.Close()
}
