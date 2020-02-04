package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

func main() {
	P("init")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")

	for i := 0; ; i++ {
		if i != 0 && i%5 == 0 {
			m = nil
			m = make(map[int]*A)
			runtime.GC()
			debug.FreeOSMemory()
			P("reset")
			fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")
			time.Sleep(15000 * time.Millisecond)
		}

		if i%2 == 0 {
			add(i)
			P("add")
			time.Sleep(300 * time.Millisecond)

			if i%3 == 0 {
				add(i + 1)
				P("add(i+1)")
				time.Sleep(300 * time.Millisecond)
			}

			remove(i)
			P("remove")
			time.Sleep(300 * time.Millisecond)

			if i%4 == 0 {
				runtime.GC()
				P("gc")
				time.Sleep(300 * time.Millisecond)
			}

			fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")
		} else {
			add(i)
			P("add")
			time.Sleep(300 * time.Millisecond)

			if i%5 == 0 {
				add(i + 1)
				P("add(i+1)")
				time.Sleep(300 * time.Millisecond)
			}

			remove(i)
			P("remove")
			time.Sleep(300 * time.Millisecond)

			if i%4 == 0 {
				debug.FreeOSMemory()
				P("free")
				time.Sleep(300 * time.Millisecond)
			}

			fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")
		}
	}
}

// func exec

type A struct {
	A string
	B string
	C string
	D string
	E string
	F string
}

var m = make(map[int]*A)
var once sync.Once
var mem runtime.MemStats

func add(base int) {
	base = base * 5000000
	for i := base; i < base+5000000; i++ {
		m[i] = &A{"A", "A", "A", "A", "A", "A"}
	}
}

func remove(base int) {
	base = base * 5000000
	for i := base; i < base+5000000; i++ {
		delete(m, i)
	}
}

func P(title string) {
	// time.Sleep(300 * time.Millisecond)
	runtime.ReadMemStats(&mem)
	once.Do(func() {
		fmt.Printf(
			"%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s\n",
			"title",
			"Sys",
			"HeapSys",
			"HeapIdle",
			"HeapInuse",
			"HeapReleased",
			"HeapAlloc",
			"HeapObjects",
			"StackSys",
			"StackInuse",
			"len(map)",
		)
	})
	fmt.Printf(
		"%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s\n",
		title,
		humanize.Bytes(mem.Sys),
		humanize.Bytes(mem.HeapSys),
		humanize.Bytes(mem.HeapIdle),
		humanize.Bytes(mem.HeapInuse),
		humanize.Bytes(mem.HeapReleased),
		humanize.Bytes(mem.HeapAlloc),
		humanize.Comma(int64(mem.HeapObjects)),
		humanize.Bytes(mem.StackSys),
		humanize.Bytes(mem.StackInuse),
		humanize.Comma(int64(len(m))),
	)
}
