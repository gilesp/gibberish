package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gilesp/markov"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	filename, port, root, splitterFlag, order, lines := parseFlags()
	fmt.Printf("Generating chain of order %d using file(s) \"%s\"...\n", order, filename)

	chain := markov.NewChainWithSplitter(order, makeSplitter(splitterFlag))
	populateChain(*chain, strings.Split(filename, ","), lines)
	chainHandler := &ChainHandler{*chain}

	http.HandleFunc("/gibberish", chainHandler.handle)
	http.Handle("/", http.FileServer(http.Dir(root)))
	fmt.Printf("Listening on port " + port + "...\n")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags() (filename string, port string, root string, splitterFlag string, order int, lines bool) {
	flag.StringVar(&filename, "file", "corpus.txt", "Filename containing source text to imitate.")
	flag.StringVar(&port, "port", "7070", "Port number to listen on.")
	flag.StringVar(&root, "root", "public_html", "Folder to serve")
	flag.StringVar(&splitterFlag, "splitter", "naive", "Sentence splitter to use. Currently accepts \"naive\", \"illinois\" or \"dub\"")
	flag.IntVar(&order, "order", 4, "Length of prefix to use.")
	flag.BoolVar(&lines, "lines", false, "Read input files line by line, rather than all in one go")
	flag.Parse()
	return
}

func makeSplitter(splitterName string) markov.Splitter {
	var splitter markov.Splitter
	switch string(splitterName) {
	case "illinois":
		fmt.Println("Using illinois sentence splitter")
		splitter = markov.NewIllinoisSplitter()
	case "dub":
		fmt.Println("Swallowing the feather")
		splitter = markov.NewDubSplitter()
	case "naive":
		fallthrough
	default:
		fmt.Println("Using naive sentence splitter")
		splitter = markov.NewNaiveSplitter()
	}

	return splitter
}

func loadTextFromFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to open file ", filename)
		return ""
	} else {
		return string(content)
	}
}

func loadLinesFromFile(filename string) []string {
	lines := []string{}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: %v", err)
		return lines
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		fmt.Println("Error reading file: %v", err)
	}
	return lines
}

func populateChain(chain markov.Chain, filenames []string, lines bool) {
	for _, filename := range filenames {
		fmt.Printf("Parsing %s...\n", filename)
		if lines {
			lines := loadLinesFromFile(filename)
			if len(lines) > 0 {
				for _, line := range lines {
					chain.Populate(line)
				}
			}
		} else {
			text := loadTextFromFile(filename)
			if text != "" {
				chain.Populate(text)
			}
		}
	}
}
