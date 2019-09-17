package main

import (
	"fmt"
	"net/http"
)

func startServer() {
	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		voterId := q.Get("voter_id")
		voteFor := q.Get("vote_for")

		if len(voterId) > 0 && len(voteFor) > 0 {
			_, _ = fmt.Fprintf(w, "Your voter id is %s\nYou vote for %s\n", voterId, voteFor)
		} else {
			_, _ = fmt.Fprint(w, "Missing parameters\n")
		}
	})

	_ = http.ListenAndServe(":8123", nil)
}
