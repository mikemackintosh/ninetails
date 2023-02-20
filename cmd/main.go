package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mikemackintosh/ninetails/internal/config"
	"github.com/mikemackintosh/ninetails/internal/version"
)

var (
	/* */
	flagWithFilename bool
	flagWithLinenum  bool
	flagWithFollow   bool
	flagVersion      bool

	/* */
	flagConfig string
)

func init() {
	flag.BoolVar(&flagWithFilename, "H", false, "Display filename")
	flag.BoolVar(&flagWithLinenum, "n", false, "Display linenum")
	flag.BoolVar(&flagVersion, "v", false, "Display version")
	flag.BoolVar(&flagWithFollow, "F", false, "Follow changes in the file")
	flag.StringVar(&flagConfig, "c", ".ninetails.yml", "Configuration file")
}

func main() {
	// Parse the config
	flag.Parse()

	// Show the version info only if it's requested
	if flagVersion {
		fmt.Printf("%s - %s\n", version.Version, version.CommitHash)
		os.Exit(0)
	}

	// Parse he configuration
	if err := config.Parse(flagConfig); err != nil {
		log.Fatal(err)
	}

	// make a channel
	messages := make(chan string)

	// Start the channel watcher
	go watcher(messages)

	// Check if the config was provided from stdin. To do so, we need
	// to see if stdin was provided via a pipe, before we start blocking
	// on the scanner read loop.
	var usingStdin bool
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	// If there is stdin, read it
	if fi.Mode()&os.ModeNamedPipe != 0 {
		usingStdin = true
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if len(scanner.Text()) > 0 {
				messages <- scanner.Text()
			}
		}
	}

	// Otherwise, look for files passed as arguments
	var wg = &sync.WaitGroup{}
	var args = flag.Args()
	if len(args) == 0 && !usingStdin {
		fmt.Printf("no files specified")
		os.Exit(1)
	}

	for _, f := range args {
		wg.Add(1)
		go readFile(wg, f, messages)
	}

	wg.Wait()
}

// readFile
func readFile(wg *sync.WaitGroup, f string, c chan string) {
	defer wg.Done()

	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		log.Fatal("file does not exist")
	}

	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewReader(file)
	var i = 0
	for {
		line, err := scanner.ReadString('\n')
		line = strings.Replace(line, "\n", "", -1)
		if err != nil {
			if err == io.EOF {
				if !flagWithFollow {
					break
				}
				time.Sleep(time.Millisecond * 100)
			} else {
				fmt.Println(err)
			}
		}

		var formattedLine = "%s"
		// If the user requested to show the filename, prep the string
		if flagWithFilename {
			// Only show line numbers when we are not following
			if flagWithLinenum && !flagWithFollow {
				i = i + 1
				formattedLine = fmt.Sprintf("%s:%d | ", f, i) + formattedLine
			} else {
				formattedLine = fmt.Sprintf("%s | ", f) + formattedLine
			}
		}

		// If we are following the file, strings reader will be empty
		// skip it only if we are following, otherwise we want to print

		if line == "" && flagWithFollow {
			continue
		}

		c <- fmt.Sprintf(formattedLine, line)
	}
}

// watcher will watch the channel and call the formatter when
// the data is received.
func watcher(c chan string) {
	for {
		v, ok := <-c
		if ok {
			formatLines(v)
		} else {
			fmt.Println("- End of input -")
			break
		}
	}
}

// formatLines will call the config replacer
func formatLines(line string) {
	if v, ok := config.Replace(line); ok {
		fmt.Println(v)
		return
	}

	fmt.Println(line)
}
