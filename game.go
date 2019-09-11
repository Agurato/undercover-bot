package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Team uint8

func (team Team) String() string {
	names := map[Team]string{
		Citizen:    "Citizen",
		Undercover: "Undercover",
		MrWhite:    "Mr. White",
	}

	return names[team]
}

type GameState uint8

type Player struct {
	user *discordgo.User
	team Team
}

type Game struct {
	session *discordgo.Session
	channel *discordgo.Channel
	state   GameState
	players []Player
}

const (
	// Team enum
	None       Team = 0
	Citizen    Team = 1
	Undercover Team = 2
	MrWhite    Team = 3

	// Game state enum
	Off     GameState = 0
	Waiting GameState = 1
	Running GameState = 2

	// Max number of players
	playerMin int = 4
)

func (g *Game) AddPlayer(user *discordgo.User) (msg string) {
	g.players = append(g.players, Player{user: user, team: None})
	msg = fmt.Sprintf("%s joined the game. %d player(s) total have joined.\n", user.Mention(), len(g.players))
	if g.IsReady() {
		msg += fmt.Sprintf("%s can start the game by typing `%s` in channel %s\n", g.players[0].user.Mention(), cmdStart, g.channel.Mention())
	}
	return msg
}

func (g *Game) SetState(state GameState) {
	g.state = state
}

func (g Game) IsReady() bool {
	return len(g.players) >= playerMin
}

func (g Game) IsOnSameChannel(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	channel, _ := s.State.GuildChannel(m.GuildID, m.ChannelID)
	return g.channel.ID == channel.ID
}

func (g Game) GetCreator() *discordgo.User {
	return g.players[0].user
}

func (g Game) SendMessage(msg string) {
	_, _ = g.session.ChannelMessageSend(g.channel.ID, msg)
}

func (g *Game) Start() {
	g.SetState(Running)
	g.SendMessage("The game has started. You all have received your word via private message")
}
