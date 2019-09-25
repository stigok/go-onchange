# go-onchange

Run a command when watched files change.

Please report any issues you encounter.

## Usage

```
Run a command when watched files change.

Usage:
onchange [options] command [args]

Options:
  -c    clear screen before each execution
  -d int
        ignore too fast reruns within time limit (ms) (default 10)
  -dry
        show a list of watched paths and exit
  -i    run command immediately (don't wait for first event)
  -k    kill long-running command on new events
  -p string
        file pattern regular expression (default ".+")
  -r    watch directory recursively
  -w string
        directory to search from (default ".")
```

### Examples

Watch all go files recursively from current directory and run `go build` on changes:

```
onchange -c -i -p \.go$ -r go build
```

Run a built binary after it has been built, and notify when build is complete:

```
onchange -i -r bash -c "go build && notify-send go 'Build complete' && ./foobin"
```

## Install

Download source and install with `go install`

```
$ git clone https://github.com/stigok/go-onchange && cd go-onchange
$ go install
```

## Building

```
$ go get github.com/radovskyb/watcher
$ go build
```
