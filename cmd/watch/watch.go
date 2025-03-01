package main

import (
	"log"

	"go.universe.tf/garden/reload"
)

func main() {
	r, err := reload.NewRunner("./garden.tmp", "-dev")
	if err != nil {
		log.Fatal(err)
	}

	tasks := []*reload.Task{
		{
			Name:    "generate sql",
			Match:   []string{"*.sql", "sqlc.yaml"},
			Command: []string{"go", "tool", "sqlc", "generate"},
		},
		{
			Name:    "generate templ",
			Match:   []string{"*.templ"},
			Command: []string{"go", "tool", "templ", "generate"},
		},
		{
			Name:   "update styles",
			Match:  []string{"*.css"},
			Notify: r.Restart,
		},
		{
			Name:    "compile Go",
			Match:   []string{"*.go"},
			Ignore:  []string{"*_test.go"},
			Command: []string{"go", "build", "-v", "-tags", "osusergo,netgo", "-o", "garden.tmp", "."},
			Notify:  r.Restart,
		},
	}

	if _, err = reload.New(tasks); err != nil {
		log.Fatal(err)
	}

	select {}
}
