package main

import (
	"encoding/json"
	"fmt"
	"github.com/gilesp/markov"
	"net/http"
	"strconv"
)

type Response map[string]interface{}

func (r Response) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

type ChainHandler struct {
	chain markov.Chain
}

func (ch *ChainHandler) gibber(maxLength int) string {
	return ch.chain.Generate(maxLength)
}

func (ch *ChainHandler) respond(message string, maxLength int) string {
	response, err := ch.chain.GenerateResponse(message, maxLength)
	if err != "" {
		response = err
	}
	return response
}

func (ch *ChainHandler) handle(w http.ResponseWriter, r *http.Request) {
	var maxLength int = 500

	var message string

	err := r.ParseForm()
	if err == nil {
		v, ok := r.Form["max"]
		if ok {
			maxLength, err = strconv.Atoi(v[0])
			if err != nil {
				maxLength = 500
			}
		}
		v, ok = r.Form["msg"]
		if ok {
			message = v[0]
		}
	}
	w.Header().Set("Content-Type", "application/json")
	var response string
	if message != "" {
		response = ch.respond(message, maxLength)
	} else {
		response = ch.gibber(maxLength)
	}
	fmt.Fprint(w, Response{"gibberish": response})
}
