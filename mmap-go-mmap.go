package main

import (
	"fmt"
	"log"
	"os"

	mmap "github.com/edsrzf/mmap-go"
)

type mmappedFileFromMmapGo struct {
	name string
	size uint32
	data mmap.MMap
}

func (f *mmappedFileFromMmapGo) Read(off, sz uint32) ([]byte, error) {
	if off > off+sz || off+sz > uint32(len(f.data)) {
		return nil, fmt.Errorf("out of bounds: %d, len %d, name %s", off+sz, len(f.data), f.name)
	}
	return f.data[off : off+sz], nil
}

func (f *mmappedFileFromMmapGo) Name() string {
	return f.name
}

func (f *mmappedFileFromMmapGo) Size() (uint32, error) {
	return f.size, nil
}

func (f *mmappedFileFromMmapGo) Close() {
	if err := f.data.Unmap(); err != nil {
		log.Printf("WARN failed to memory unmap %s: %v", f.name, err)
	}
}

func MapFileUsingMmapGo(f *os.File) (MmappedFile, error) {
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	sz := fi.Size()
	if sz >= maxUInt32 {
		return nil, fmt.Errorf("file %s too large: %d", f.Name(), sz)
	}
	r := &mmappedFileFromMmapGo{
		name: f.Name(),
		size: uint32(sz),
	}

	// round up to the OS page size because mmap likes to align on pages
	// mmap will zero-fill the extra bytes
	pagesize := uint32(os.Getpagesize() - 1)
	rounded := (r.size + pagesize) &^ pagesize

	r.data, err = mmap.MapRegion(f, int(rounded), mmap.RDONLY, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("MapFileUsingMmapGo: unable to memory map %s: %w", f.Name(), err)
	}

	return r, err
}
