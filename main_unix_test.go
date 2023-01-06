//go:build !windows

package main

import (
	"testing"
)

func BenchmarkUnix(b *testing.B) {
	if err := doBenchmark("BenchmarkUnix", MapFileUsingUnixMmap); err != nil {
		b.Errorf("%s", err)
	}
}
