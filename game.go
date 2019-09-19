package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
)

type GameState uint8

// A Game represents a game, whichever state it is in
type Game struct {
	session *discordgo.Session // discord's library session object
	channel *discordgo.Channel // discord's library channel object
	state   GameState          // game's current state
	players []Player           // slice of players who joined the game
	votes   map[string]uint    // current votes for eliminating players (votes[userID] = <number of votes against userID>)
}

const (
	// Game state enum
	Off     GameState = 0
	Waiting GameState = 1
	Running GameState = 2

	// Max number of players
	playerMin int = 4
)

// String returns a GameState as a string
func (state GameState) String() string {
	states := map[GameState]string{
		Off:     "Off",
		Waiting: "Waiting",
		Running: "Running",
	}

	return states[state]
}

// AddPlayer adds a player to the game. A player can't register twice (the variable checked is user.ID).
// Returns a message to be sent back to players
func (g *Game) AddPlayer(user *discordgo.User) (msg string) {
	// Check that the player isn't already in the game
	for _, p := range g.players {
		if user.ID == p.user.ID {
			return fmt.Sprintf("%s, you already joined the game!\n", user.Mention())
		}
	}

	// Add the player to the slice
	g.players = append(g.players, Player{user: user, team: None})
	msg = fmt.Sprintf("%s joined the game. %d player(s) total have joined.\n", user.Mention(), len(g.players))
	if g.IsReady() {
		msg += fmt.Sprintf("%s can start the game by typing `%s` in channel %s\n", g.players[0].user.Mention(), cmdStart, g.channel.Mention())
	}
	return msg
}

// GetPlayerFromId returns the player with id given as parameter
func (g Game) GetPlayerFromId(id string) (player Player) {
	for _, p := range g.players {
		if id == p.user.ID {
			player = p
			return
		}
	}

	return
}

// SetState changes the current state of the game.
func (g *Game) SetState(state GameState) {
	g.state = state
}

// Reset resets the game, set all vars to nil and set state to Off.
func (g *Game) Reset() {
	g.session = nil
	g.channel = nil
	g.players = nil
	g.votes = nil
	g.SetState(Off)
}

// SetRandomTeams set random teams to all players, according to the number of Undercovers and Mr. Whites
func (g *Game) SetRandomTeams(undercoverNumber, mrWhiteNumber int64) {
	// Remaining players without a team
	playersTeamNone := make([]int, len(g.players))
	// Add everyone in it at the beginning
	for i := 0; i < len(playersTeamNone); i++ {
		playersTeamNone[i] = i
	}

	// Put random players in the Undercover team
	for i := 0; i < int(undercoverNumber); i++ {
		index := rand.Intn(len(playersTeamNone))
		g.players[playersTeamNone[index]].team = Undercover
		playersTeamNone = append(playersTeamNone[:index], playersTeamNone[index+1:]...)
	}

	// Put random players in the Mr. White team
	for i := 0; i < int(mrWhiteNumber); i++ {
		index := rand.Intn(len(playersTeamNone))
		g.players[playersTeamNone[index]].team = MrWhite
		playersTeamNone = append(playersTeamNone[:index], playersTeamNone[index+1:]...)
	}

	// Put all remaining players in the Citizen team
	for i := 0; i < len(playersTeamNone); i++ {
		g.players[playersTeamNone[i]].team = Citizen
	}
}

// IsReady checks if the minimum number of players is reached.
func (g Game) IsReady() bool {
	return len(g.players) >= playerMin
}

// IsOnSameChannel checks if a message received is on the same channel as the game.
func (g Game) IsOnSameChannel(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	channel, _ := s.State.GuildChannel(m.GuildID, m.ChannelID)
	return g.channel.ID == channel.ID
}

// GetCreator returns the discord's library user object of the creator of the game
func (g Game) GetCreator() *discordgo.User {
	return g.players[0].user
}

// SendMessage sends a message to the channel where the game was created.
func (g Game) SendMessage(msg string) {
	_, _ = g.session.ChannelMessageSend(g.channel.ID, msg)
}

// Start starts the game, after all players have been added
func (g *Game) Start(undercoverNumber, mrWhiteNumber int64) {
	// Change state of the game to Running
	g.SetState(Running)

	// Put players in teams
	g.SetRandomTeams(undercoverNumber, mrWhiteNumber)

	// TODO: remove the following lines
	g.SendWords(GenerateWords())
	g.SendVotes()
	g.SendMessage(fmt.Sprintf("The game has started. You all have received your word and the list of players via private message.\nYou can vote for a player by clicking on his name"))

	// TODO: uncomment the following lines to close the game
	//if g.SendWords(GenerateWords()) {
	//	g.SendMessage(fmt.Sprintf("The game has started. You all have received your word via private message.\nYou can start voting for a player with `%s @player` in private message", cmdVote))
	//} else {
	//	g.SendMessage("Some players didn't receive their word. Exiting the game...")
	//	g.Reset()
	//}

}

// SendWords sends each player their word according to their respective teams.
// word1 is sent to Citizens, word2 is sent to Undercovers
func (g *Game) SendWords(word1, word2 string) bool {
	success := true
	// For each player
	for _, p := range g.players {
		// Establish private message communication
		userChannel, err := g.session.UserChannelCreate(p.user.ID)
		if err != nil {
			// TODO: uncomment
			//g.SendMessage(fmt.Sprintf("The bot couldn't send the word to %s", p.user.Mention()))
			success = false
			continue
		}
		var wordMsg string
		// Send a different for each team
		switch p.team {
		case Citizen:
			wordMsg = fmt.Sprintf("Your word is **%s**.\n", word1)
		case Undercover:
			wordMsg = fmt.Sprintf("Your word is **%s**.\n", word2)
		case MrWhite:
			wordMsg = fmt.Sprintf("You are a %s, you don't have a word.\n", MrWhite.String())
		}
		// Sends PM
		_, _ = g.session.ChannelMessageSend(userChannel.ID, wordMsg)
	}
	return success
}

// SendVotes sends a PM to each player for them to be able to vote.
// The message contains a list of other players where they can click to vote for someone.
// This uses the HTTP server started at the beginning when users click on of the names.
func (g *Game) SendVotes() bool {
	success := true
	// Reset the votes map
	g.ResetVotes()

	for i, p := range g.players {
		userChannel, err := g.session.UserChannelCreate(p.user.ID)
		if err != nil {
			// TODO: uncomment
			//g.SendMessage(fmt.Sprintf("The bot couldn't send the votes to %s", p.user.Mention()))
			success = false
			continue
		}

		var playerList string
		for j, otherPlayer := range g.players {
			if i != j {
				playerList += fmt.Sprintf("- [%s](%s?%s=%s&%s=%s)\n", otherPlayer.user.Username, serverUrl, voterIdParam, p.user.ID, voteForParam, otherPlayer.user.ID)
			}
		}
		// Create Embed message for hyperlinks
		_, _ = g.session.ChannelMessageSendEmbed(userChannel.ID, &discordgo.MessageEmbed{
			Title:       "Let's vote!",
			Description: fmt.Sprintf("You can vote for one of these players:\n%s", playerList),
			Color:       2148295,
		})
	}
	return success
}

// ResetVotes resets the votes map
func (g *Game) ResetVotes() {
	g.votes = make(map[string]uint)
	// Set an entry for each player with a value of 0
	for i, p := range g.players {
		g.players[i].canVote = true
		g.votes[p.user.ID] = 0
	}
}

// Vote adds a vote from player to another player.
// Returns true if there are still some players who have not voted
func (g *Game) Vote(voter, voteFor string) bool {
	stillVoting := false
	// For each player
	for i, p := range g.players {
		// if the user voting is able to vote
		if p.user.ID == voter {
			if p.canVote {
				// Add a vote to the total of the voteFor user
				g.votes[voteFor]++
				// Remove the right to vote
				g.players[i].canVote = false
			}
		}
		// Checks if players are still voting
		stillVoting = stillVoting || g.players[i].canVote
	}

	return stillVoting
}

// GenerateWords generate a pair of words with some ressemblance
// TODO: the logic of this function (this is the hardest part of the game)
func GenerateWords() (word1, word2 string) {

	word1 = "pomme"
	word2 = "poire"

	return
}
