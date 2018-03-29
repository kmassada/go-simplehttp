// simplehttp is just a simple http server
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func handle(w http.ResponseWriter, r *http.Request) {
	sleepLength := 5
	var l int // Length

	b := make([]byte, 2048*1024)
	for {
		n, err := r.Body.Read(b)
		l += n // update length
		fmt.Fprintf(w, "Read %d more bytes, total is now %d bytes. Sleeping %d ms\n", n, l, sleepLength)
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error caught: %v", err)
		}
		if err == io.EOF {
			break
		}
		time.Sleep(time.Duration(sleepLength) * time.Millisecond) // Simulate a slow upload
	}
	r.Body.Close()
	fmt.Fprintf(w, "\n\nRead %d total mb.\n", l/1024/1024)
}

func headers(w http.ResponseWriter, r *http.Request) {
	for hkey, hval := range r.Header {
		fmt.Fprintf(w, "%s => %q\n", hkey, hval)
	}
	fmt.Fprintf(w, "\n-------------\n\n")
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		headers(w, r)

		if r.Method != http.MethodPost {
			fmt.Fprintln(os.Stdout, "Send a POST request, please.")
			return
		}
		handle(w, r)
	})

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    time.Minute,
		WriteTimeout:   time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
