package main

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkMmapGo(b *testing.B) {
	if err := doBenchmark("BenchmarkMmap", MapFileUsingMmapGo); err != nil {
		b.Errorf("%s", err)
	}
}

func BenchmarkFS(b *testing.B) {
	if err := doBenchmark("BenchmarkFS", MapFileUsingFS); err != nil {
		b.Errorf("%s", err)
	}
}

func TestMain(m *testing.M) {
	if err := generateFile(); err != nil {
		fmt.Printf("setup failed to generate test file: %s\n", err)
		os.Exit(1)
	}
	code := m.Run()
	os.Remove("testfile")
	os.Exit(code)
}
