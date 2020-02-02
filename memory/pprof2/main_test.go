package main

import "testing"

func BenchmarkMap_A(b *testing.B) {
	m := make(map[int64]struct{})

	for i := 0; i < b.N; i++ {
		m[int64(i)] = struct{}{}
	}
}

func BenchmarkMap_B(b *testing.B) {
	m := make(map[int64]struct{})

	for i := 0; i < b.N; i++ {
		m[int64(i)] = struct{}{}
		delete(m, int64(i))
	}
}

func BenchmarkMap_C(b *testing.B) {
	m := make(map[int64]struct{})

	for i := 0; i < b.N; i++ {
		m[int64(i)] = struct{}{}
	}

	for i := 0; i < b.N; i++ {
		delete(m, int64(i))
	}
}

func BenchmarkMap_D(b *testing.B) {
	m := make(map[int]interface{})

	for i := 0; i < b.N; i++ {
		add(m, struct{}{}, i)
	}
}

func BenchmarkMap_E(b *testing.B) {
	m := make(map[int]interface{})

	for i := 0; i < b.N; i++ {
		add(m, true, i)
	}
}

func BenchmarkMap_F(b *testing.B) {
	m := make(map[int]interface{})

	for i := 0; i < b.N; i++ {
		add(m, true, i)
		remove(m, i)
	}
}
