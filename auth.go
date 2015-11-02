package main

import (
	"fmt"
	"net/http"
)

/// AuthHello message struct.
/// Decoded from JSON, has username and password fields.
type AuthHello struct {
	User string
	Pass string
}

type AuthResponse struct {
	LobbyID int `json:"lobby_id"`
	Success bool
}

/// auth_hello handles authentication check.
/// Hits API to check if the credentials are valid.
/// Assumes any response other than a 200 is an indication.
/// that the crdentials provided were invalid.
func (p *Player) auth_hello(data []byte) {
	var message AuthHello
	if err := DecodeJSON(data, &message); err != nil {
		fmt.Println(err)
		p.SendError(err.Error())
		return
	}

	// Perform request
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8080/login", nil)
	req.Header.Set("username", message.User)
	req.Header.Set("password", message.Pass)
	resp, err := client.Do(req)

	// If there was an error, or the status code is not 200, auth failed
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(err)
		p.SendError("Could not authenticate")
		return
	}

	// Add player to global room, update name, etc.
	global_room.AddPlayer(p)
	p.Room = &global_room
	p.State = InLobby
	p.Username = message.User

	// Send auth_success message
	response := AuthResponse{
		LobbyID: 1,
		Success: true,
	}
	p.SendMessage("auth_response", response)
	resp.Body.Close()
}
