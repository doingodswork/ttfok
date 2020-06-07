package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
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

	// Check if app is already running
	http.DefaultClient.Timeout = time.Second
	_, err := http.Get(appURL)
	if err == nil {
		log.Fatal("The app seems to be running already")
	}
	urlErr, ok := err.(*url.Error)
	if !ok {
		log.Fatalf("Expected `*url.Error` during check if the app is already running, but error type %T; error: %v", err, err)
	}
	if !strings.Contains(urlErr.Error(), "connection refused") {
		log.Fatalf("Expected \"connection refused\" error, but was: %v", urlErr.Error())
	}

	// Run app
	cmd := exec.Command(app)
	start := time.Now()
	if err := cmd.Start(); err != nil {
		log.Fatalf("Couldn't start app: %v", err)
	}
	// Kill process at the end
	waitExceeded := false
	defer func() {
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Couldn't kill process (PID: %v): %v", cmd.Process.Pid, err)
		}
		if waitExceeded {
			os.Exit(1)
		}
	}()

	// Make requests
	http.DefaultClient.Timeout = time.Millisecond
	for {
		res, err := http.Get(appURL)
		if err == nil && res.StatusCode == http.StatusOK {
			break
		}
		if time.Since(start) >= time.Second {
			log.Println("App didn't start within 1s. Exiting...")
			waitExceeded = true
			break
		}
	}

	if !waitExceeded {
		log.Printf("Time to first OK: %v", time.Since(start))
	}
}
