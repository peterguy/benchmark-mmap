package main

import (
	"fmt"
	"os"
)

func MapFileUsingFS(f *os.File) (MmappedFile, error) {
	return &mmappedFileFromFS{f}, nil
}

type mmappedFileFromFS struct {
	f *os.File
}

func (f *mmappedFileFromFS) Read(off, sz uint32) ([]byte, error) {
	r := make([]byte, sz)
	_, err := f.f.ReadAt(r, int64(off))
	return r, err
}

func (f mmappedFileFromFS) Size() (uint32, error) {
	fi, err := f.f.Stat()
	if err != nil {
		return 0, err
	}

	sz := fi.Size()

	if sz >= maxUInt32 {
		return 0, fmt.Errorf("overflow")
	}

	return uint32(sz), nil
}

func (f mmappedFileFromFS) Close() {
	f.f.Close()
}

func (f mmappedFileFromFS) Name() string {
	return f.f.Name()
}
