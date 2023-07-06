package main

import (
	"encoding/json"
	"errors"
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
	cmd := new(exec.Cmd)
	sReq := strings.Fields(c.Request)
	switch {
	case len(sReq) == 1:
		cmd = exec.Command(sReq[0])
	default:
		cmd = exec.Command(sReq[0], sReq[1:]...)
	}
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("bad command %v: %w", sReq, err)
	}
	switch {
	case len(output) == 0: // no output (e.g. $touch), return success message
		output = fmt.Appendf(output, "successful command %v: %v\n", sReq, cmd.ProcessState)
		return output, nil
	default:
		return output, nil
	}
}

func validate(r *http.Request) (*CmdRequest, error) {
	cr := new(CmdRequest)
	dc := json.NewDecoder(r.Body)
	dc.DisallowUnknownFields() // request only or return err
	err := dc.Decode(cr)
	if err != nil {
		err = errors.New(err.Error() + "; expected \"request\"")
		return nil, err
	}
	sReq := strings.Fields(strings.ToLower(cr.Request))
	switch {
	case len(sReq) == 0:
		err = fmt.Errorf("missing command %v", sReq)
		return nil, err
	case strings.Contains(sReq[0], "sudo"):
		err = errors.New("SUDO command not allowed")
		return nil, err
	}
	return cr, nil
}

// Handler for POST requests with a shell command
func handleCmdPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate JSON field and invalidate a SUDO command
	cr, err := validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
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
