package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
)

func main() {
	fset := flag.NewFlagSet("main", flag.ExitOnError)
	ipURL := fset.String("ip", "localhost:8088/ip", "URL for IP reporter")
	uaURL := fset.String("ua", "localhost:8088/ua", "URL for User Agent String reporter")
	port := fset.Int("port", 8088, "listen port")
	debug := fset.Bool("debug", false, "Debugging messages to stdout")

	fset.Parse(os.Args[1:])

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		host := r.Host

		fullURL := path.Join(host, uri)

		var o string

		if fullURL == *ipURL {
			if *debug {
				fmt.Println("Matched address:", *ipURL)
			}
			a := r.RemoteAddr
			i, _, err := net.SplitHostPort(a)
			if err != nil {
				o = "error"
				w.WriteHeader(http.StatusInternalServerError)
			}
			x, ok := r.Header["X-Forwarded-For"]
			if ok {
				i = x[0]
			}
			o = i
		} else if fullURL == *uaURL {
			if *debug {
				fmt.Println("Matched address:", *uaURL)
			}
			o = r.UserAgent()
		} else {
			if *debug {
				fmt.Println("Could not match", fullURL, "to either", *ipURL, "or", *uaURL)
			}
			o = "error"
			w.WriteHeader(http.StatusBadRequest)
		}

		o = fmt.Sprintf("%s\n", o)
		fmt.Fprintf(w, o)

	})

	p := fmt.Sprintf(":%d", *port)

	log.Fatal(http.ListenAndServe(p, nil))
}
