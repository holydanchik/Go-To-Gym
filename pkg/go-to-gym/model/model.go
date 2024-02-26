package model

import "fmt"

type Exercises struct {
	name string
	reps int
}

func Output() {
	rdl := Exercises {
		name: "Romanian Deadlift",
		reps: 12,
	}
	fmt.Println(rdl.name, rdl.reps)
}