package main

import (
	"log"

	"go.universe.tf/garden/gogarden/reload"
)

func main() {
	r, err := reload.NewRunner("./garden.bin", "-dev")
	if err != nil {
		log.Fatal(err)
	}

	tasks := []*reload.Task{
		{
			Name:    "generate sql",
			Match:   []string{"*.sql", "sqlc.yaml"},
			Command: []string{"sqlc", "generate"},
		},
		{
			Name:    "generate templ",
			Match:   []string{"*.templ"},
			Command: []string{"templ", "generate"},
		},
		{
			Name:    "generate tailwind",
			Match:   []string{"style.css", "*.templ"},
			Command: []string{"tailwindcss", "-i", "style.css", "-o", "static/app.css"},
		},
		{
			Name:    "compile Go",
			Match:   []string{"*.go"},
			Ignore:  []string{"*_test.go"},
			Command: []string{"make", "bin"},
			Notify:  r.Restart,
		},
	}

	if _, err = reload.New(tasks); err != nil {
		log.Fatal(err)
	}

	select {}
}
