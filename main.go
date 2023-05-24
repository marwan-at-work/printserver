package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	status = flag.Int("status", 200, "response status to send")
	body   = flag.String("body", "ok", "response body to send")
	addr   = flag.String("addr", "127.0.0.1:4545", "server address to listen on")
	wait   = flag.Duration("wait", 0, "simulate a longer response by sleeping for a duration")
)

func main() {
	flag.Parse()
	fmt.Printf("listening on http://%s\n", *addr)
	http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("%s: %s\n", r.Method, r.URL.String())
		fmt.Println("Headers:")
		for k := range r.Header {
			fmt.Printf("%s: %s\n", k, r.Header.Get(k))
		}
		fmt.Println("Body:")
		if r.Header.Get("content-type") == "application/json" {
			var buf bytes.Buffer
			src, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Println("ERROR READING BODY:", err)
				return
			}
			err = json.Indent(&buf, src, "", "\t")
			if err != nil {
				fmt.Println("ERROR INDENTING JSON:", err)
			}
			fmt.Println(buf.String())
		} else {
			io.Copy(os.Stdout, r.Body)
			fmt.Println()
		}
		if *wait > 0 {
			time.Sleep(*wait)
		}
		w.WriteHeader(*status)
		fmt.Fprintf(w, "%s\n", *body)
		fmt.Printf("response status: %d duration: %s\n", *status, time.Since(start))
	}))
}
