package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/dustin/go-humanize"
	"github.com/pkg/profile"
)

var one sync.Once

func P(title string) {
	// time.Sleep(1 * time.Second)
	runtime.ReadMemStats(&mem)
	one.Do(func() {
		fmt.Printf(
			"%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s\n",
			"title",
			"HeapSys",
			"HeapAlloc",
			"HeapInuse",
			"HeapIdle",
			"HeapObjects",
			"HeapReleased",
			"StackSys",
			"StackInuse",
			"Sys",
		)
	})
	fmt.Printf(
		"%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s%-13s\n",
		title,
		humanize.Bytes(mem.HeapSys),
		humanize.Bytes(mem.HeapAlloc),
		humanize.Bytes(mem.HeapInuse),
		humanize.Bytes(mem.HeapIdle),
		humanize.Comma(int64(mem.HeapObjects)),
		humanize.Bytes(mem.HeapReleased),
		humanize.Bytes(mem.StackSys),
		humanize.Bytes(mem.StackInuse),
		humanize.Bytes(mem.Sys),
	)
}

type A struct {
	A string
	B string
	C string
	D string
	E string
	F string
}

func main() {
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	// debug.FreeOSMemory()
	exec()
}

var p = fmt.Println

func memUsage(m1, m2 *runtime.MemStats) {
	p("Alloc:", m2.Alloc-m1.Alloc,
		"TotalAlloc:", m2.TotalAlloc-m1.TotalAlloc,
		"HeapAlloc:", m2.HeapAlloc-m1.HeapAlloc)
}

var mem runtime.MemStats

func exec() {
	m := make(map[int]interface{})
	P("init")

	for i := 0; i < 1000000; i++ {
		m[i] = true
		m[i] = &A{"A", "A", "A", "A", "A", "A"}
	}
	P("add1")

	runtime.GC()
	P("gc")

	debug.FreeOSMemory()
	P("free")

	for i := 1000000; i < 1000000+1000000; i++ {
		m[i] = false
	}
	P("add2")

	for i := 0; i < 1000000; i++ {
		delete(m, i)
	}
	P("rm1")

	for i := 1000000; i < 1000000+1000000; i++ {
		delete(m, i)
	}
	P("rm2")

	runtime.GC()
	P("gc")

	debug.FreeOSMemory()
	P("free")

	runtime.GC()
	P("gc")

	debug.FreeOSMemory()
	P("free")

	m = nil
	P("nil")

	runtime.GC()
	P("gc")

	debug.FreeOSMemory()
	P("free")
}
