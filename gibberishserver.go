package main

import (
	"flag"
	"fmt"
	"github.com/gilesp/markov"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type ChainHandler struct {
	chain markov.Chain
}

func (ch *ChainHandler) handle(w http.ResponseWriter, r *http.Request) {
	var maxLength int = 500
	err := r.ParseForm()
	if err == nil {
		v, ok := r.Form["max"]
		if ok {
			maxLength, err = strconv.Atoi(v[0])
			if err != nil {
				maxLength = 500
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, Response{"gibberish": ch.chain.Generate(maxLength)})
}

func main() {
	var filename = flag.String("file", "corpus.txt", "Filename containing source text to imitate.")
	var order = flag.Int("order", 4, "Length of prefix to use.")
	var port = flag.String("port", "7070", "Port number to listen on.")
	var root = flag.String("root", "public_html", "Folder to serve")
	var splitterFlag = flag.String("splitter", "naive", "Sentence splitter to use. Currently accepts \"naive\" or \"illinois\"")
	flag.Parse()

	fmt.Printf("Generating chain of order %d using file(s) \"%s\"...\n", *order, *filename)
	var splitter markov.Splitter
	switch string(*splitterFlag) {
	case "illinois":
		fmt.Println("Using illinois sentence splitter")
		splitter = markov.NewIllinoisSplitter()
	case "naive":
		fallthrough
	default:
		fmt.Println("Using naive sentence splitter")
		splitter = markov.NewNaiveSplitter()
	}
	chain := markov.NewChainWithSplitter(*order, splitter)
	populateChain(*chain, strings.Split(*filename, ","))
	chainHandler := &ChainHandler{*chain}

	http.HandleFunc("/gibberish", chainHandler.handle)
	http.Handle("/", http.FileServer(http.Dir(*root)))
	fmt.Printf("Listening on port " + *port + "...\n")
	http.ListenAndServe(":"+*port, nil)
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

func populateChain(chain markov.Chain, filenames []string) {
	for _, filename := range filenames {
		fmt.Printf("Parsing %s...\n", filename)
		text := loadTextFromFile(filename)
		if text != "" {
			chain.Populate(text)
		}
	}
}
