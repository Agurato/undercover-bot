package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"testing"
)

func TestGame_SetState(t *testing.T) {
	var g Game
	if g.state != Off {
		t.Errorf("Game state should be %s, but is %s", Off.String(), g.state.String())
	}
	g.SetState(Waiting)
	if g.state != Waiting {
		t.Errorf("Game state should be %s, but is %s", Waiting.String(), g.state.String())
	}
	g.SetState(Running)
	if g.state != Running {
		t.Errorf("Game state should be %s, but is %s", Running.String(), g.state.String())
	}
}

func TestGame_AddPlayer(t *testing.T) {
	g := Game{channel: &discordgo.Channel{ID: "test"}}

	if len(g.players) != 0 {
		t.Errorf("There should be 0 player, there are %d", len(g.players))
	}

	dummy := &discordgo.User{
		ID: "dummy0",
	}
	g.AddPlayer(dummy)
	if len(g.players) != 1 {
		t.Errorf("There should be 1 player, there are %d", len(g.players))
	}

	for i := 0; i < 6; i++ {
		dummy = &discordgo.User{
			ID: fmt.Sprintf("dummy%d", i),
		}
		g.AddPlayer(dummy)
	}
	if len(g.players) != 6 {
		t.Errorf("There should be 6 players, there are %d", len(g.players))
	}
}

func TestGame_SetRandomTeams(t *testing.T) {
	var g Game
	playerNumber := 6
	var undercoverNumber int64 = 2
	var mrWhiteNumber int64 = 1
	var citizenNumber int64 = int64(playerNumber) - undercoverNumber - mrWhiteNumber

	// Add 6 players
	for i := 0; i < playerNumber; i++ {
		g.players = append(g.players, Player{user: nil, team: None})
	}
	// Set Random Teams
	g.SetRandomTeams(undercoverNumber, mrWhiteNumber)

	// Check if number of players hasn't changed
	if len(g.players) != playerNumber {
		t.Errorf("There should be %d players, there are %d", playerNumber, len(g.players))
	}

	// Check number of players per team
	var undercoverCount, mrWhiteCount, citizenCount int64 = 0, 0, 0
	for i := 0; i < playerNumber; i++ {
		if g.players[i].team == Undercover {
			undercoverCount++
		} else if g.players[i].team == MrWhite {
			mrWhiteCount++
		} else if g.players[i].team == Citizen {
			citizenCount++
		}
	}
	if undercoverCount != undercoverNumber {
		t.Errorf("There should be %d in Team %s, there are %d", undercoverNumber, Undercover.String(), undercoverCount)
	}
	if mrWhiteCount != mrWhiteNumber {
		t.Errorf("There should be %d in Team %s, there are %d", mrWhiteNumber, MrWhite.String(), mrWhiteCount)
	}
	if citizenCount != citizenNumber {
		t.Errorf("There should be %d in Team %s, there are %d", citizenNumber, Citizen.String(), citizenCount)
	}

}

func TestGame_ResetVotes(t *testing.T) {
	g := Game{channel: &discordgo.Channel{ID: "test"}}
	playerNumber := 4
	// Add 4 players
	for i := 0; i < playerNumber; i++ {
		dummy := &discordgo.User{
			ID: fmt.Sprintf("dummy%d", i),
		}
		g.AddPlayer(dummy)
	}

	if len(g.votes) != 0 {
		t.Errorf("There should be no votes, there are %d", len(g.votes))
	}

	for _, p := range g.players {
		if p.canVote {
			t.Errorf("%s shouldn't be able to vote!", p.user.ID)
		}
	}

	g.ResetVotes()

	if len(g.votes) != 4 {
		t.Errorf("There should be 4 votes, there are %d", len(g.votes))
	}

	for _, p := range g.players {
		if g.votes[p.user.ID] != 0 {
			t.Errorf("%s should have 0 vote, he has %d", p.user.ID, g.votes[p.user.ID])
		}
		if !p.canVote {
			t.Errorf("%s should be able to vote!", p.user.ID)
		}
	}
}

func TestGame_Vote(t *testing.T) {
	g := Game{channel: &discordgo.Channel{ID: "test"}}
	playerNumber := 4
	stillVoting := false
	// Add 4 players
	for i := 0; i < playerNumber; i++ {
		dummy := &discordgo.User{
			ID: fmt.Sprintf("dummy%d", i),
		}
		g.AddPlayer(dummy)
	}

	if len(g.votes) != 0 {
		t.Errorf("There should be no votes, there are %d", len(g.votes))
	}

	g.ResetVotes()

	stillVoting = g.Vote("dummy0", "dummy1")
	if !stillVoting {
		t.Errorf("There should still be players voting")
	}
	stillVoting = g.Vote("dummy1", "dummy2")
	if !stillVoting {
		t.Errorf("There should still be players voting")
	}
	stillVoting = g.Vote("dummy2", "dummy1")
	if !stillVoting {
		t.Errorf("There should still be players voting")
	}
	stillVoting = g.Vote("dummy3", "dummy0")
	if stillVoting {
		t.Errorf("All players should have voted")
	}

	if len(g.votes) != playerNumber {
		t.Errorf("There should be %d votes, there are %d", playerNumber, len(g.votes))
	}

	if g.votes["dummy0"] != 1 {
		t.Errorf("dummy0 should have 1 vote, there are %d", g.votes["dummy0"])
	}
	if g.votes["dummy1"] != 2 {
		t.Errorf("dummy1 should have 2 votes, there are %d", g.votes["dummy1"])
	}
	if g.votes["dummy2"] != 1 {
		t.Errorf("dummy1 should have 1 vote, there are %d", g.votes["dummy2"])
	}
	if g.votes["dummy3"] != 0 {
		t.Errorf("dummy1 should have 0 vote, there are %d", g.votes["dummy3"])
	}
}
