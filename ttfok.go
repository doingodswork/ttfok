package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func main() {
	flag.Parse()

	// Precondition check
	if len(flag.Args()) != 2 {
		log.Fatalf("Bad number of arguments. Requires 2, but was %v", len(flag.Args()))
	}

	// Read arguments
	app := flag.Arg(0)
	appURL := flag.Arg(1)

	// Run app
	cmd := exec.Command(app)
	start := time.Now()
	if err := cmd.Start(); err != nil {
		log.Fatalf("Couldn't start app: %v", err)
	}

	// Make requests
	http.DefaultClient.Timeout = time.Millisecond
	for {
		res, err := http.Get(appURL)
		if err == nil && res.StatusCode == http.StatusOK {
			break
		}
	}

	log.Printf("Time to first OK: %v", time.Since(start))
}
