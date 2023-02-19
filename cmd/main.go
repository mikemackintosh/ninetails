package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"ninetails/config"
	"ninetails/version"
	"os"
	"sync"
)

var (
	/* */
	flagWithFilename bool
	flagWithLinenum  bool
	flagVersion      bool

	/* */
	flagConfig string
)

func init() {
	flag.BoolVar(&flagWithFilename, "H", false, "Display filename")
	flag.BoolVar(&flagWithLinenum, "n", false, "Display linenum")
	flag.BoolVar(&flagVersion, "v", false, "Display version")
	flag.StringVar(&flagConfig, "c", ".ninetail.yml", "Configuration file")
}

func main() {
	// Parse the config
	flag.Parse()

	if flagVersion {
		fmt.Printf("%s - %s\n", version.Version, version.CommitHash)
		os.Exit(0)
	}

	if err := config.Parse(flagConfig); err != nil {
		log.Fatal(err)
	}

	// make a channel
	messages := make(chan string)

	var wg = &sync.WaitGroup{}

	go printChannelData(messages)

	// Check if the config was provided from stdin. To do so, we need
	// to see if stdin was provided via a pipe, before we start blocking
	// on the scanner read loop.
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if len(scanner.Text()) > 0 {
				messages <- scanner.Text()
			}
		}
	}

	var args = flag.Args()
	if len(args) > 0 {
		for _, f := range args {
			wg.Add(1)
			go readFile(wg, f, messages)
		}
	}

	wg.Wait()
}

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

	scanner := bufio.NewScanner(file)
	var i = 0
	for scanner.Scan() {
		line := scanner.Text()

		var formattedLine = ""

		if flagWithFilename {
			if flagWithLinenum {
				i = i + 1
				formattedLine = fmt.Sprintf("%s:%d | "+formattedLine, f, i)
			} else {
				formattedLine = fmt.Sprintf("%s | "+formattedLine, f)
			}
		}

		c <- fmt.Sprintf(formattedLine+"%s", line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func printChannelData(c chan string) {
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

func formatLines(line string) {
	if v, ok := config.Replace(line); ok {
		fmt.Println(v)
		return
	}

	fmt.Println(line)
}
