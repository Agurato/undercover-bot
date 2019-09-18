package main

import "github.com/bwmarrin/discordgo"

type Team uint8

func (team Team) String() string {
	names := map[Team]string{
		Citizen:    "Citizen",
		Undercover: "Undercover",
		MrWhite:    "Mr. White",
	}

	return names[team]
}

type Player struct {
	user    *discordgo.User
	team    Team
	canVote bool
}

const (
	// Team enum
	None       Team = 0
	Citizen    Team = 1
	Undercover Team = 2
	MrWhite    Team = 3
)
