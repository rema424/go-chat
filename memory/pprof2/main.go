package main

import (
	"fmt"
	"strconv"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	fmt.Println("start")
	exec()
	fmt.Println("end")
}

func exec() {
	var someMap = make(map[int64]struct{})

	doNothing("foo", 100000)
	addMap(someMap, "buzz", 1000000)
	removeMap(someMap)
	doNothing("foo", 100000)
}

func addMap(s map[int64]struct{}, prefix string, count int) {
	for i := 0; i < count; i++ {
		s[int64(i)] = struct{}{}
	}
}

func removeMap(s map[int64]struct{}) {
	for key := range s {
		delete(s, key)
	}
}

func doNothing(prefix string, count int) {
	for i := 0; i < count; i++ {
		key := prefix + "key" + strconv.Itoa(i)
		val := "value" + strconv.Itoa(i)
		key = key + val
	}
}

func add(m map[int]interface{}, elm interface{}, base int) {
	base *= 1000
	for i := base; i < base+1000; i++ {
		m[i] = elm
	}
}

func remove(m map[int]interface{}, base int) {
	base *= 1000
	for i := base; i < base+1000; i++ {
		delete(m, i)
	}
}
