package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Counter struct {
	mu  sync.Mutex
	x   int
}

var locker Counter

func get(writer http.ResponseWriter, req *http.Request) {
	locker.mu.Lock()
	log.Printf("GET counter request: %v", locker.x)
	fmt.Fprintf(writer, "Counter is at: %d\n", locker.x)
	locker.mu.Unlock()
}

func set(writer http.ResponseWriter, req *http.Request) {
	log.Printf("SET counter request: %v", req.RequestURI)
	value := req.URL.Query().Get("value")
	intval, err := strconv.Atoi(value)

	if err != nil {
		log.Println("SET handler: non-integer parameter value.")
	}

	locker.mu.Lock()
	locker.x = intval
	log.Printf("counter set to: %v", locker.x)
	fmt.Fprintf(writer, "Counter set to: %d\n", locker.x)
	locker.mu.Unlock()
}

func inc(_ http.ResponseWriter, _ *http.Request) {
	// time.Sleep(1 * time.Second)
	locker.mu.Lock()
	locker.x = locker.x + 1
	log.Printf("counter incremented to: %v", locker.x)
	locker.mu.Unlock()
}

func dec(_ http.ResponseWriter, _ *http.Request) {
	locker.mu.Lock()
	locker.x = locker.x - 1
	log.Printf("counter decremented to: %v", locker.x)
	locker.mu.Unlock()
}

func main() {
	locker = Counter{
		x: 0,
		mu: sync.Mutex{},
	}
	http.HandleFunc("/counter", get)
	http.HandleFunc("/counter/set", set)
	http.HandleFunc("/increment", inc)
	http.HandleFunc("/decrement", dec)

	port := 9095
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	}
	log.Printf("Listening on port %d\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), nil))
}
