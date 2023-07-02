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
	Request  string `json:"request"`
	Response string `json:"response"`
	// Response []byte `json:"response"`
	// Error    error  `json:"error"`
}

func (c *CmdRequest) exCmd() error {
	sCmd := strings.Fields(c.Request)
	fmt.Printf("sCmd: %v\n", sCmd[1])
	cmd := exec.Command(sCmd[0], sCmd[1:]...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("cmd.Output err:", err.Error())
		return errors.New(err.Error())
	}
	fmt.Println("out:", string(out))
	fmt.Printf("\ncmd output: %v\n", out)
	return nil
}

// Handler for POST requests with a shell command
func handleCmdPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// new(Command) returns a pointer type which is used by Decode(*T), but maybe use var declaration instead
	c := new(CmdRequest)
	err := json.NewDecoder(r.Body).Decode(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Printf("cmd decoded: %+v\n", c)

	// Call function to execute the shell command and return err
	err = c.exCmd()
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

}

func main() {
	http.HandleFunc("/api/cmd", handleCmdPost)
	fmt.Println("starting the server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
