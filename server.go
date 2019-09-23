package main

import (
	"fmt"
	"net/http"
)

const (
	serverURL    string = "https://undercover.vmonot.dev/"
	voterIDParam string = "voter_id"
	voteForParam string = "vote_for"
)

// StartServer starts the HTTP server for people to vote.
// Entry point is :8123/vote with GET parameters voter_id and vote_for as IDs of discord's users
func StartServer() {
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		voterID := q.Get(voterIDParam)
		voteFor := q.Get(voteForParam)

		if len(voterID) > 0 && len(voteFor) > 0 {
			_, _ = fmt.Fprintf(w, "You voted for @%s\n", game.GetPlayerFromID(voteFor).user.Username)
			game.Vote(voterID, voteFor)
		} else {
			_, _ = fmt.Fprint(w, "Missing parameters\n")
		}
	})

	_ = http.ListenAndServe(":8123", nil)
}
