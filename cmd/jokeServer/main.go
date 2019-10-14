package main

import (
	j "appleTakeHome/pkg/jokester"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var jokester j.Jokester
	portPtr := flag.Int("port", 5000, "port to listen for requests on")
	flag.Parse()
	sigCh := make(chan os.Signal)
	err := jokester.Init()
	if err != nil {
		log.Printf("Error initializing jokester: %v", err)
	}
	// Notify on the following os.Signals for graceful shutdown
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", jokester.HandleNameJoke)
	go http.ListenAndServe(fmt.Sprintf(":%v", *portPtr), nil)
	log.Printf("Listening on port: %v\n", *portPtr)

	//block for signal & deinit upon receipt
	s := <-sigCh
	jokester.Deinit()
	log.Printf("Received: %v\nExiting", s.String())
}
