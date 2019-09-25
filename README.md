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

```
$ go install github.com/stigok/go-onchange
```

Make sure you have `$GOPATH/bin` added to your `$PATH`.

## Building

```
$ go build
```

## License

The underlying library used to watch files is written by Benjamin Radovsky (BSD 3-Clause "New" or "Revised" License).
Read the full license text in [pkg/github.com/radovskyb/watcher/LICENSE]().
