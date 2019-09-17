package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
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
	// Check that the player isn't already in the game
	for _, p := range g.players {
		if user.ID == p.user.ID {
			return fmt.Sprintf("%s, you already joined the game!\n", user.Mention())
		}
	}
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

func (g *Game) Reset() {
	g.session = nil
	g.channel = nil
	g.players = nil
	g.SetState(Off)
}

func (g *Game) SetRandomTeams(undercoverNumber, mrWhiteNumber int64) {
	playersTeamNone := make([]int, len(g.players))
	for i := 0; i < len(playersTeamNone); i++ {
		playersTeamNone[i] = i
	}
	for i := 0; i < int(undercoverNumber); i++ {
		index := rand.Intn(len(playersTeamNone))
		g.players[playersTeamNone[index]].team = Undercover
		playersTeamNone = append(playersTeamNone[:index], playersTeamNone[index+1:]...)
	}
	for i := 0; i < int(mrWhiteNumber); i++ {
		index := rand.Intn(len(playersTeamNone))
		g.players[playersTeamNone[index]].team = MrWhite
		playersTeamNone = append(playersTeamNone[:index], playersTeamNone[index+1:]...)
	}
	for i := 0; i < len(playersTeamNone); i++ {
		g.players[playersTeamNone[i]].team = Citizen
	}
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

func (g *Game) Start(undercoverNumber, mrWhiteNumber int64) {
	g.SetState(Running)

	g.SetRandomTeams(undercoverNumber, mrWhiteNumber)

	// TODO: remove the following lines
	g.SendWords(GenerateWords())
	g.SendMessage(fmt.Sprintf("The game has started. You all have received your word via private message.\nYou can start voting for a player with `%s @player` in private message", cmdVote))

	// TODO: uncomment the following lines to close the game
	//if g.SendWords(GenerateWords()) {
	//	g.SendMessage(fmt.Sprintf("The game has started. You all have received your word via private message.\nYou can start voting for a player with `%s @player` in private message", cmdVote))
	//} else {
	//	g.SendMessage("Some players didn't receive their word. Exiting the game...")
	//	g.Reset()
	//}

}

func (g *Game) SendWords(word1, word2 string) bool {
	success := true
	for _, p := range g.players {
		userChannel, err := g.session.UserChannelCreate(p.user.ID)
		if err != nil {
			g.SendMessage(fmt.Sprintf("The bot couldn't send the word to %s", p.user.Mention()))
			success = false
			continue
		}
		var wordMsg string
		switch p.team {
		case Citizen:
			wordMsg = fmt.Sprintf("Your word is **%s**.\n", word1)
		case Undercover:
			wordMsg = fmt.Sprintf("Your word is **%s**.\n", word2)
		case MrWhite:
			wordMsg = fmt.Sprintf("You are a %s, you don't have a word.\n", MrWhite.String())
		}
		_, _ = g.session.ChannelMessageSend(userChannel.ID, wordMsg)
	}
	return success
}

func GenerateWords() (word1, word2 string) {
	word1 = "pomme"
	word2 = "poire"

	return
}
