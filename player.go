package main

import "github.com/bwmarrin/discordgo"

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
	// Team enum
	None       Team = 0
	Citizen    Team = 1
	Undercover Team = 2
	MrWhite    Team = 3
)
