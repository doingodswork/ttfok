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

var (
	timeout = flag.Duration("t", time.Millisecond, "Timeout for the request")
	wait    = flag.Duration("w", time.Second, "Duration to wait for app start")
)

func main() {
	flag.Parse()

	// Precondition check
	argCount := len(flag.Args())
	if argCount < 2 {
		log.Fatalf("Bad number of arguments. Requires 2, but was %v", argCount)
	}

	// Read arguments
	app := flag.Arg(0)
	appURL := flag.Arg(argCount - 1)
	var args []string
	if argCount > 2 {
		args = append(args, flag.Args()[1:argCount-1]...)
	}

	// Check if app is already running
	http.DefaultClient.Timeout = *timeout + time.Second
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
	cmd := exec.Command(app, args...)
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
	http.DefaultClient.Timeout = *timeout
	for {
		res, err := http.Get(appURL)
		if err == nil && res.StatusCode == http.StatusOK {
			break
		}
		if time.Since(start) >= *wait {
			log.Printf("App didn't start within %v. Exiting...", *wait)
			waitExceeded = true
			break
		}
	}

	if !waitExceeded {
		log.Printf("Time to first OK: %v", time.Since(start))
	}
}
