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

type Unauthorized struct {
}

func (u *Unauthorized) Handle(data Message, p *Player) {
	if data.Type != "auth_hello" {
		p.SendError("User not yet authenticated")
		return
	}

	var message AuthHello
	if err := DecodeJSON(data.Data, &message); err != nil {
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
	p.Mutex.Lock()
	p.Room = &global_room
	p.State = new(GameLobby)
	p.Username = message.User
	p.Mutex.Unlock()

	// Send auth_success message
	response := AuthResponse{
		LobbyID: 1,
		Success: true,
	}
	p.SendMessage("auth_response", response)
}
