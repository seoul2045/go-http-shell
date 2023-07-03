package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type CmdRequest struct {
	Request string `json:"request"`
	// Response string `json:"response"`
	// Response []byte `json:"response"`
	// Error    error  `json:"error"`
}

func (c *CmdRequest) exCmd() ([]byte, error) {
	fmt.Println("exCmd c.Request:", c.Request)

	sCmd := strings.Fields(c.Request)
	cmd := exec.Command(sCmd[0], sCmd[1:]...)
	fmt.Printf("sCmd: %v\n", cmd)

	// cmd.Dir = "/Users/greg/Workspace/enableit-api"
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output err:", err.Error())
		return nil, err
	}
	fmt.Println("out:", string(output))
	fmt.Printf("\ncmd output: %v\n", output)
	return output, nil
}

// Handler for POST requests with a shell command
func handleCmdPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c := new(CmdRequest)
	body := make([]byte, r.ContentLength)
	_, err := io.ReadFull(r.Body, body)
	if err != nil {
		fmt.Println(err) // TODO: handler err properly
		return
	}
	err = json.Unmarshal(body, c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Printf("cmd decoded: %+v\n", c)

	// Call function to execute the shell command and return err
	output, err := c.exCmd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	bs := make([]byte, r.ContentLength)
	r.Body.Read(bs)

	_, err = fmt.Fprintf(w, "\nRequest Body: %v\nHTTP Method: %s\nHeader: %v\n", string(bs), r.Method, r.Header)
	if err != nil {
		fmt.Println("handleCmdPost error:", err)
		return
	}

	// Write response reply
	w.Header()
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
