package main

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
)

type Game struct {
	Host      string
	Players   []string
	Map, Name string
	ID        string
}

type AvailableGames struct {
	Games []*Game
}

type NewGame struct {
	Map, Name, ID string
}

type HostResponse struct {
	GameID string
}

type JoinGame struct {
	GameID string
}

var hosted_games []*Game
var waiting_players []*Player

func notifyWaiting(game *Game) {
	newGame := NewGame{
		Map:  game.Map,
		Name: game.Name,
		ID:   game.ID,
	}

	for _, player := range waiting_players {
		player.SendMessage("new_game", newGame)
	}
}

func (p *Player) join_game(data []byte) {
	var mess JoinGame
	if err := DecodeJSON(data, &mess); err != nil {
		p.SendError(err.Error())
		return
	}

	for _, game := range hosted_games {
		if game.ID == mess.GameID {
			game.Players = append(game.Players, p.Username)
			p.State = InGame
			p.SendMessage("game_joined", nil)
			return
		}
	}

	p.SendError("Game could not be found")
}

func (p *Player) request_games(data []byte) {
	waiting_players = append(waiting_players, p)
	resp := AvailableGames{
		Games: hosted_games,
	}
	p.SendMessage("available_games", resp)
}

// host_game indicates that a player wants to host a game.
// The player is added to a list of players that are hoping to host a game,
// as well as the name of the map that the game will be played on.
// A response message including the UUID for the game is returned
func (p *Player) host_game(data []byte) {
	if p.State == Unknown {
		p.SendError("User is not logged in")
		return
	}

	var game *Game
	if err := DecodeJSON(data, &game); err != nil {
		fmt.Println(err)
		p.SendError(err.Error())
		return
	}
	game.ID = uuid.New()
	game.Players = append(game.Players, p.Username)
	game.Host = p.Username

	// Add the game to the list of hosted games
	hosted_games = append(hosted_games, game)

	// Indicate that the game was successfuly listed
	resp := HostResponse{
		GameID: game.ID,
	}
	p.SendMessage("host_response", resp)

	// Tell all players that are currently waiting to join a game that
	// one is available
	notifyWaiting(game)
}
