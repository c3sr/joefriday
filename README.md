joefriday
=========
[![Build Status](https://travis-ci.org/c3sr/joefriday.svg?branch=master)](https://travis-ci.org/c3sr/joefriday)
> "All we want are the facts, ma'am"  
>   - Joe Friday

JoeFriday is a group of libraries that gathers system information: cpu, disk, memory, network, system, etc. This information can be returned as Go structs, Flatbuffers serialized bytes, or JSON serialized bytes. For Flatbuffers and JSON, deserialization convenience methods are provided. When it makes sense, a Ticker based implementation is provided to enable periodic gathering of information.

JoeFriday seeks to minimize allocations and time spent gathering the information. For minimal resource usage, use the sysinfo implementations, if appropriate, as the data provided by those implementations use syscalls, which are at least an order of magnitude faster than processing the proc files.

This library only supports Linux.

See package specific READMEs for information about what those packages provide.

## Benchmarks
### Comparative Benchmarks
The `bench` package provides comparative benchmarks; showing how JoeFriday implementations compare with some other libraries that provide similar information.

The information provided by the various libraries are not the same; that should also be factored into the decision as to which library best suits your requirements.

Either a benchmark app can be compiled with `go build` or they can be run using `go test -bench=. [flags]`.

The benchmark app's output includes group, e.g. CPU, Memory, etc, package name and function called, Ops, ns/Op, B/op, and Allocs/Op. The output can be formatted as text lines (default), CSV, or Markdown Table and written to a file.

### JoeBench
JoeBench is an app that runs benchmarks of JoeFriday functionality, including serialization implementations, and generates formatted output. Output includes group, e.g. CPU, Memory, etc, package name and function called, Ops, ns/Op, B/op, and Allocs/Op. The output can be formatted as text lines (default), CSV, or Markdown Table and written to a file.

## Notes:
A big thanks to [Eric Lagergren](https://github.com/EricLagergren) for all of his help and his contributions.

## `alpha-1` branch:
The alpha-1 branch has the original JoeFriday implementation. If your code relies on the package's state prior to 7/2017, this branch should be used if you don't want to change your existing code to work with the refactored JoeFriday.

## TODO
* Provide protobuf implementations.  
* Add CPU speed info: min, max, current.  
* For utilization and usage, add output of deltas between snapshots.
* For utilization and usage revisit calculations and algorithms used, maybe add additional algorithms, where appropriate. This may be a separate library.  
