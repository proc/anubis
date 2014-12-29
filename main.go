package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
)

func HomePage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Anubis")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomePage)

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
