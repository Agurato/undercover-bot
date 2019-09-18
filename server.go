package main

import (
	"fmt"
	"net/http"
)

const (
	serverUrl    string = "https://undercover.vmonot.dev/"
	voterIdParam string = "voter_id"
	voteForParam string = "vote_for"
)

func startServer() {
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		voterId := q.Get(voterIdParam)
		voteFor := q.Get(voteForParam)

		if len(voterId) > 0 && len(voteFor) > 0 {
			_, _ = fmt.Fprintf(w, "You voted for @%s\n", game.GetPlayerFromId(voteFor).user.Username)
			game.Vote(voterId, voteFor)
		} else {
			_, _ = fmt.Fprint(w, "Missing parameters\n")
		}
	})

	_ = http.ListenAndServe(":8123", nil)
}
