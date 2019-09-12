package main

import "testing"

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
