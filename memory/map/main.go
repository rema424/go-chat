package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	debug.FreeOSMemory()
	var (
		m1 runtime.MemStats
		m2 runtime.MemStats
		m3 runtime.MemStats
		m4 runtime.MemStats
		m5 runtime.MemStats
		m6 runtime.MemStats
		m7 runtime.MemStats
	)
	runtime.ReadMemStats(&m1)

	t := struct{}{}
	time.Sleep(1)
	runtime.ReadMemStats(&m2)
	fmt.Printf("%#v\n", t)
	memUsage(&m1, &m2)
	P()
	fmt.Println()

	b := false
	time.Sleep(1)
	runtime.ReadMemStats(&m3)
	fmt.Printf("%#v\n", b)
	memUsage(&m2, &m3)
	P()
	fmt.Println()

	// debug.FreeOSMemory()
	runtime.GC()
	time.Sleep(1)
	fmt.Printf("%#v\n", "GC")
	runtime.ReadMemStats(&m4)
	memUsage(&m3, &m4)
	P()
	fmt.Println()

	m := make(map[int]interface{})
	time.Sleep(1)
	fmt.Printf("%#v\n", m)
	runtime.ReadMemStats(&m5)
	memUsage(&m4, &m5)
	P()
	fmt.Println()

	for i := 0; i < 1000000; i++ {
		m[i] = true
	}
	time.Sleep(1)
	runtime.ReadMemStats(&m6)
	memUsage(&m5, &m6)
	P()
	fmt.Println()

	for i := 0; i < 1000000; i++ {
		delete(m, i)
	}
	time.Sleep(1)
	runtime.ReadMemStats(&m7)
	memUsage(&m6, &m7)
	P()
	fmt.Println()
}

var p = fmt.Println

func memUsage(m1, m2 *runtime.MemStats) {
	p("Alloc:", m2.Alloc-m1.Alloc,
		"TotalAlloc:", m2.TotalAlloc-m1.TotalAlloc,
		"HeapAlloc:", m2.HeapAlloc-m1.HeapAlloc)
}

var mem runtime.MemStats

func P() {
	runtime.ReadMemStats(&mem)
	fmt.Println(mem.HeapSys, mem.HeapAlloc, mem.Alloc, mem.Frees)
}
