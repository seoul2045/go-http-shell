package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type CmdRequest struct {
	Request string `json:"request"`
}

// Method to execute a shell command
func (c *CmdRequest) exCmd() ([]byte, error) {
	sReq := strings.Fields(c.Request)
	cmd := new(exec.Cmd)
	switch {
	case len(sReq) == 0:
		return nil, fmt.Errorf("missing command %v", sReq)
	case len(sReq) == 1:
		cmd = exec.Command(sReq[0])
	default:
		cmd = exec.Command(sReq[0], sReq[1:]...)
	}
	output, err := cmd.Output()
	fmt.Printf("25 cmd.ProcessState: %v\n", cmd.ProcessState)
	fmt.Printf("26 output: %v\n", output)
	if err != nil {
		return nil, fmt.Errorf("bad command %v: %w", sReq, err)
	}
	switch {
	case len(output) == 0:
		output = fmt.Appendf(output, "successful command %v: %v\n", sReq, cmd.ProcessState)
		return output, nil
	default:
		return output, nil
	}
}

// Handler for POST requests with a shell command
func handleCmdPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cr := new(CmdRequest)
	dc := json.NewDecoder(r.Body)
	dc.DisallowUnknownFields()
	err := dc.Decode(cr)
	if err != nil {
		http.Error(w, err.Error()+"; expected \"request\"", http.StatusNotFound)
		return
	}

	// Call method to execute the shell command
	output, err := cr.exCmd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Write response reply
	w.Header().Set("Content-type", "text/plain")
	_, err = w.Write(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/api/cmd", handleCmdPost)
	fmt.Println("starting the server...")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
