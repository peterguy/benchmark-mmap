//go:build !windows

package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

type mmappedFileFromUnixMmap struct {
	name string
	size uint32
	data []byte
}

func (f *mmappedFileFromUnixMmap) Read(off, sz uint32) ([]byte, error) {
	if off > off+sz || off+sz > uint32(len(f.data)) {
		return nil, fmt.Errorf("out of bounds: %d, len %d, name %s", off+sz, len(f.data), f.name)
	}
	return f.data[off : off+sz], nil
}

func (f *mmappedFileFromUnixMmap) Name() string {
	return f.name
}

func (f *mmappedFileFromUnixMmap) Size() (uint32, error) {
	return f.size, nil
}

func (f *mmappedFileFromUnixMmap) Close() {
	if err := unix.Munmap(f.data); err != nil {
		log.Printf("WARN failed to Munmap %s: %v", f.name, err)
	}
}

func MapFileUsingUnixMmap(f *os.File) (MmappedFile, error) {
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	sz := fi.Size()
	if sz >= maxUInt32 {
		return nil, fmt.Errorf("file %s too large: %d", f.Name(), sz)
	}
	r := &mmappedFileFromUnixMmap{
		name: f.Name(),
		size: uint32(sz),
	}

	// round up to the OS page size because mmap likes to align on pages
	// mmap will zero-fill the extra bytes
	pagesize := uint32(os.Getpagesize() - 1)
	rounded := (r.size + pagesize) &^ pagesize

	r.data, err = unix.Mmap(int(f.Fd()), 0, int(rounded), unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	return r, err
}
