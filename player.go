package main

import "github.com/bwmarrin/discordgo"

// Team represents a player's team (enum)
type Team uint8

// String returns Team name as a string
func (team Team) String() string {
	names := map[Team]string{
		None:       "None",
		Citizen:    "Citizen",
		Undercover: "Undercover",
		MrWhite:    "Mr. White",
	}

	return names[team]
}

// A Player represents a player who joined a game.
type Player struct {
	user    *discordgo.User // discord's library user object
	team    Team            // team of the player
	canVote bool            // bool set to true when the player is able to vote
}

const (
	// None is the initial team for players
	None Team = 0
	// Citizen is the team of Citizens
	Citizen Team = 1
	// Undercover is the team of Undercovers
	Undercover Team = 2
	// MrWhite is the team of Mr. Whites
	MrWhite Team = 3
)
