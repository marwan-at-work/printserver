package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	status = flag.Int("status", 200, "response status to send")
	body   = flag.String("body", "ok", "response body to send")
	port   = flag.String("port", "4545", "port number the server runs on")
)

func main() {
	flag.Parse()
	http.ListenAndServe(":"+*port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(*status)
		fmt.Fprintf(w, "%s\n", *body)
	}))
}
