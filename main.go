package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"onchange/pkg/watcher"
)

// Execute a cancellable command
func Exec(ctx context.Context, cmd string, args []string) error {
	c := exec.CommandContext(ctx, cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		log.Println("failed to execute command:", err)
		return err
	}

	return c.Wait()
}

// Clear terminal (linux only)
func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

// Function to be called when -h flag is used
func Usage() {
	out := flag.CommandLine.Output()

	fmt.Fprint(out, "Run a command when watched files change.\n\n")
	fmt.Fprintf(out, "Usage:\n%s [options] command [args]\n\n", os.Args[0])
	fmt.Fprint(out, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	var (
		pat, dir, cmd                  string
		rec, clr, immediate, kill, dry bool
		rundelay, pollInt              int64
		cmdargs                        []string
	)

	flag.StringVar(&dir, "w", ".",
		"directory to search from")
	flag.StringVar(&pat, "f", ".+",
		"file pattern regular expression")
	flag.BoolVar(&rec, "r", false,
		"watch directory recursively")
	flag.BoolVar(&clr, "c", false,
		"clear screen before each execution")
	flag.BoolVar(&immediate, "i", false,
		"run command immediately (don't wait for first event)")
	flag.BoolVar(&kill, "k", false,
		"kill long-running command on new events")
	flag.BoolVar(&dry, "dry", false,
		"show a list of watched paths and exit")
	flag.Int64Var(&pollInt, "p", 100,
		"file poll interval (ms)")
	flag.Int64Var(&rundelay, "d", 10,
		"ignore too fast reruns within time limit (ms)")

	flag.Usage = Usage
	flag.Parse()

	// Create a file watcher
	w := watcher.New()
	w.SetMaxEvents(1) // We only care that "something" has happened
	w.FilterOps(watcher.Create, watcher.Write, watcher.Remove, watcher.Rename,
		watcher.Move)

	// Use a file pattern if it exists
	if len(pat) > 0 {
		r, err := regexp.Compile(pat)
		if err != nil {
			log.Fatalln("Failed to compile match expression", err)
		}
		w.AddFilterHook(watcher.RegexFilterHook(r, false))
	}

	// Watch folder for changes, recursively if specified
	if rec {
		if err := w.AddRecursive(dir); err != nil {
			log.Fatalln(err)
		}
	} else {
		if err := w.Add(dir); err != nil {
			log.Fatalln(err)
		}
	}

	// Print a list of all paths watched
	if dry {
		paths := w.WatchedFiles()
		fmt.Printf("Watching %d paths:\n", len(paths))
		for path, _ := range paths {
			fmt.Println(path)
		}
		os.Exit(0)
	}

	// Set up variables containing the command to be executed
	switch flag.NArg() {
	case 0:
		fmt.Fprintln(os.Stdout, "error: missing command\n")
		Usage()
		os.Exit(1)
	case 1:
		cmd = flag.Arg(0)
	default:
		cmd = flag.Arg(0)
		cmdargs = flag.Args()[1:]
	}

	// Use a context to allow a running command to be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Run command immediately if specified
	if immediate {
		if clr {
			ClearScreen()
		}
		go Exec(ctx, cmd, cmdargs)
	}

	go func() {
		lastrun := time.Now()
		for {
			select {
			case <-w.Event:
				// Ignore events that are happening too fast
				if time.Since(lastrun) < time.Duration(100+rundelay)*time.Millisecond {
					continue
				}

				// Cancel last run, in case it's still running
				if kill {
					cancel()
				}
				ctx, cancel = context.WithCancel(context.Background())

				// Clean output from previous run
				if clr {
					ClearScreen()
				}
				go Exec(ctx, cmd, cmdargs)
				lastrun = time.Now()
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				cancel()
				return
			}
		}
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * time.Duration(pollInt)); err != nil {
		log.Fatalln(err)
	}
}
