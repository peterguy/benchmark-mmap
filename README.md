# benchmark-mmap

Benchmark Go `mmap` implementations

## The `golang.org/x/sys/unix` package

Platform-specific: works on Linux and Unix variants, but not Windows.

## The `github.com/edsrzf/mmap-go` package

One of several cross-platform `mmap` wrapper packages that I found.

It has the cleanest API of the ones I found.

Under the hood it uses `mmap` on *Nix systems and on Windows the `CreateFileMapping` and `MapViewOfFile` Windows API methods.

It exposes a `[]byte`, which matches the API of `unix.Mmap`, and decorates that with methods to manipulate it, like `Read()` and `Close()`.

## `os.File.ReadAt()`

which isn't _really_ memory mapping, but it is randomly reading parts of a file.

It was/is included as a non-*nix "memory map" option in the [zoekt project](https://github.com/sourcegraph/zoekt), which is why it's included here.

It's so much slower than memory mapping, it doesn't really deserve consideration.

# Notes

The benchmark setup involves creating a 1GB-sized file in the current working directory filled with random bytes.

Each benchmark consists of opening/mapping that file, reading 20,000,000 50k chunks at random locations, and then unmapping/closing the file.

The cleanup at the end of the benchmarks deletes the 1GB file.

To run the benchmarks
```bash
go test -bench .
```

You can filter the benchmarks by passing regular expressions to `-bench` instead of the dot (`.`). To avoid the _slooooow_ file-based "mapping" benchmark, use the regular expression `'.*(Unix|Mmap)'`

Example of using Docker for benchmarking on other platforms
```bash
docker run --rm --volume $(pwd):/src --volume ${GOPATH:-${HOME}/go}:/go --workdir /src --interactive golang <<<"go test -bench ."
```