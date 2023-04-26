package main

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"os"
	"os/exec"
	"io"
)

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	fmt.Printf("Running command: %s %s\n", cmd, strings.Join(args, " "))
	return exec.Command(cmd, args...).Start()
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("Usage: ./url-launcher [<interface:port>]")
		return
	}

	addr := "0.0.0.0:8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error reading body: %v\n", err)
			http.Error(w, "Invalid body", http.StatusInternalServerError)
			return
		}

		urlString := strings.TrimSpace(string(body))
		urlString = strings.ReplaceAll(urlString, "\n", "")

		match, _ := regexp.MatchString(`\Ahttps?://`, urlString)
		if match {
			fmt.Printf("Opening URL: %s\n", urlString)
			err := open(urlString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, "OK")
		} else {
			fmt.Printf("Invalid URL: %s\n", urlString)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
		}
	})

	fmt.Printf("Server listening on %s...\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
