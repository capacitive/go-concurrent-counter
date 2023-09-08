package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type FirstSetup struct {
	mu            sync.Mutex
	firstCallFlag bool
}

var (
	counter    int        = 0
	firstSetup FirstSetup = FirstSetup{firstCallFlag: true}
)

// NOTE: Channels should be made and cleaned up in the controlling function, in this case main(). For the purposes of the exercise they're being declared and initialized here so they are global. They're also not being cleaned up as that should really happen in the main() function.
var counterChannel chan int = make(chan int)

// setup sets up a listener thread to make changes to the counter variable as the values are sent to the channel.
func setup() {
	log.Println("IN the setup call")
	go func() {
		for counterDelta := range counterChannel {
			counter = counter + counterDelta
			log.Printf("counter changed to: %v", counter)
		}
	}()
}

func get(writer http.ResponseWriter, req *http.Request) {
	log.Printf("GET counter request: %v", counter)
	fmt.Fprintf(writer, "Counter is at: %d\n", counter)
}

func set(writer http.ResponseWriter, req *http.Request) {
	log.Printf("SET counter request: %v", req.RequestURI)
	value := req.URL.Query().Get("value")
	intval, err := strconv.Atoi(value)

	if err != nil {
		log.Println("SET handler: non-integer parameter value.")
	}

	counter = intval
	log.Printf("counter set to: %v", counter)
	fmt.Fprintf(writer, "Counter set to: %d\n", counter)
}

func inc(_ http.ResponseWriter, _ *http.Request) {
	// normally setting up the channel listeners should happen in the controlling thread but I'm doing it here for the exercise.

	// for almost all the inc() calls this will be false. checking a bool is way faster than trying to grab a mutex
	if firstSetup.firstCallFlag {
		// for the threads that make it here, try to grab the mutex so only 1 thread does the setup
		firstSetup.mu.Lock()
		defer firstSetup.mu.Unlock()
		// one more check of the flag, if it's already set to false then return so setup() is not called multiple times
		if firstSetup.firstCallFlag {
			firstSetup.firstCallFlag = false
			setup()
		}
	}

	// spin up a small go func to add to the channel so this thread can return quickly
	go func() {
		counterChannel <- 1
	}()

	// moved this print statement to the listener thread on the counter channel in the setup() function
}

// decrements the counter by 1 by sending -1 to the channel. channels make this sort of concurrency super easy
func dec(_ http.ResponseWriter, _ *http.Request) {
	go func() {
		counterChannel <- -1
	}()
}

func main() {
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
