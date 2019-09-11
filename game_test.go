package main

import "testing"

func TestTeam_String(t *testing.T) {
	team := Citizen

	if team.String() != "Citizen" {
		t.Error("Expected \"Citizen\", got "+team.String())
	}
}