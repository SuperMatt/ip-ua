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

	fset.Parse(os.Args[1:])

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		host := r.Host

		fullURL := path.Join(host, uri)

		var o string

		if fullURL == *ipURL {
			a := r.RemoteAddr
			i, _, err := net.SplitHostPort(a)
			if err != nil {
				o = "error"
				w.WriteHeader(http.StatusInternalServerError)
			}
			o = i
		} else if fullURL == *uaURL {
			o = r.UserAgent()
		} else {
			o = "error"
			w.WriteHeader(http.StatusBadRequest)
		}

		fmt.Fprintf(w, o)

	})

	p := fmt.Sprintf(":%d", *port)

	log.Fatal(http.ListenAndServe(p, nil))
}
