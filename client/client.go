package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Request body is a JSON string
	fmt.Println("Received request with JSON body: ", string(body))

	// Respond with a JSON object
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"type": "FRAME","message": "from clientServer"}`))
}
func main() {
	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":6969", nil)
}
