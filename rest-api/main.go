package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const config = "/etc/sample-api/info.cfg"

func main() {
	http.HandleFunc("/", Handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(config)
	if err != nil {
		log.Fatalf("Error: Unable to open %s: %s", config, err)
		io.Copy(w, strings.NewReader("Unable to open config file"))
		return
	}
	io.Copy(w, f)
}
