package main

import (
	"fmt"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	fmt.Println("start")
	new()
	add()
	remove()
	fmt.Println("end")
}

func exec() {
	new()
	add()
	remove()
}

var m map[int64]struct{}

func new() {
	m = make(map[int64]struct{})
}

func add() {
	for i := 0; i < 1000; i++ {
		m[int64(i)] = struct{}{}
	}
}

func remove() {
	for i := 0; i < 1000; i++ {
		delete(m, int64(i))
	}
}
