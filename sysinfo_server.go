// sysinfo_server serves a REST API of `systemd-analyze time`
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const systemdAnalyze string = "/usr/bin/systemd-analyze"
const version int = 1

type handlerMap map[string]func(http.ResponseWriter, *http.Request)

var handlers = handlerMap{
	"/":         handleRoot,
	"/version":  handleVersion,
	"/duration": handleDuration,
}

// writeError writes a 500 Internal Server Error
func writeError(w http.ResponseWriter, err error) {
	log.Println(fmt.Sprintf("Error %v", err))
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

// handleRoot handles the root request (/) by just running systemd-analyze
func handleRoot(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command(systemdAnalyze)
	cmd.Stdout = w
	if err := cmd.Run(); err != nil {
		writeError(w, err)
		return
	}
}

// handleVersion handles the version request
func handleVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", version)
}

// handleDuration handles the duration request
func handleDuration(w http.ResponseWriter, r *http.Request) {
	var out bytes.Buffer

	cmd := exec.Command(systemdAnalyze)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		writeError(w, err)
		return
	}

	var line string
	scanner := bufio.NewScanner(&out)
	for i := 0; i < 1 && scanner.Scan(); i++ {
		line = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		writeError(w, err)
		return
	}

	if !strings.Contains(line, " = ") {
		writeError(w, errors.New(fmt.Sprintf("expected '=' in string '%v'", line)))

	}

	total := strings.Split(line, " = ")[1]
	fmt.Fprintf(w, "%v", total)
}

// printReady prints that the server is ready and all available endpoints
func printReady(handlers handlerMap) {
	endpoints := make([]string, 0, len(handlers))
	for e := range handlers {
		endpoints = append(endpoints, e)
	}
	fmt.Printf("Server ready, endpoints: %v\n", strings.Join(endpoints, ", "))
}

func main() {
	for handler, handler_fn := range handlers {
		http.HandleFunc(handler, handler_fn)
	}
	printReady(handlers)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
