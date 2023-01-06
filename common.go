// common things used by other files

package main

import (
	"fmt"
	"math/rand"
	"os"
)

const maxUInt32 = 0xffffffff

const count = 20_000_000

const bite = 50_000

// MmappedFile is a file suitable for concurrent read access. For performance
// reasons, it allows a mmap'd implementation.
type MmappedFile interface {
	Read(off uint32, sz uint32) ([]byte, error)
	Size() (uint32, error)
	Close()
	Name() string
}

func doBenchmark(function string, mapFunction func(f *os.File) (MmappedFile, error)) error {
	file, err := openFile()
	if err != nil {
		return fmt.Errorf("%s: failed opening test file: %w", function, err)
	}
	mapped, err := mapFunction(file)
	if err != nil {
		return fmt.Errorf("%s: failed creating memory-mapped file: %w", function, err)
	}
	if err := process(mapped); err != nil {
		return fmt.Errorf("%s: process failed: %w", function, err)
	}
	return nil
}

func process(mapped MmappedFile) error {
	size, _ := mapped.Size()

	for i := 0; i < count; i++ {
		offset := rand.Intn(int(size - bite))
		stuff, err := mapped.Read(uint32(offset), bite)
		if err != nil {
			return fmt.Errorf("process: failed reading %d bytes from offset %d: %w", bite, offset, err)
		}
		if len(stuff) != bite {
			return fmt.Errorf("process: FAIL: read %d bytes instead of %d", len(stuff), bite)
		}
	}
	mapped.Close()
	return nil
}

func generateFile() error {
	if _, err := os.Stat("testfile"); err != nil {
		data := make([]byte, 1024*1024*1024)
		amt, err := rand.Read(data)
		if err != nil {
			return fmt.Errorf("Failed to generate a testing file: %w", err)
		}
		if amt != len(data) {
			return fmt.Errorf("Failed to generate a testing file: %w", err)
		}
		os.WriteFile("testfile", data, os.ModePerm)
	}
	return nil
}

func openFile() (*os.File, error) {
	file, err := os.Open("testfile")
	if err != nil {
		return nil, fmt.Errorf("openFile: failed opening %s: %w", "testfile", err)
	}
	return file, nil
}
