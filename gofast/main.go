// Copyright 2013 Davide Muzzarelli <d.muzzarelli@dav-muz.net>. All rights reserved.
// Use of this source code is governed by a BSD-style

/*
Example of a Go web site with FastCGI
Read the full article on http://www.dav-muz.net/blog/2013/09/how-to-use-go-and-fastcgi/
*/

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"runtime"
)

var (
	local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")
	tcp   = flag.String("tcp", "", "serve as FCGI via TCP, example: 0.0.0.0:8000")
	unix  = flag.String("unix", "", "serve as FCGI via UNIX socket, example: /tmp/myprogram.sock")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func homeView(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")
	io.WriteString(w, "<html><head></head><body><p>It works!</p></body></html>")
}

func info(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")
	io.WriteString(w, "<html><head></head><body><p>Request Info</p><table>")
	io.WriteString(w, fmt.Sprintf("<tr><td>Method</td><td>%s</td></tr>", r.Method))
	io.WriteString(w, fmt.Sprintf("<tr><td>URL</td><td>%s</td></tr>", r.URL))
	io.WriteString(w, fmt.Sprintf("<tr><td>URL.Path</td><td>%s</td></tr>", r.URL.Path))
	io.WriteString(w, fmt.Sprintf("<tr><td>Proto</td><td>%s</td></tr>", r.Proto))
	io.WriteString(w, fmt.Sprintf("<tr><td>Host</td><td>%s</td></tr>", r.Host))
	io.WriteString(w, fmt.Sprintf("<tr><td>RemoteAddr</td><td>%s</td></tr>", r.RemoteAddr))
	io.WriteString(w, fmt.Sprintf("<tr><td>RequestURI</td><td>%s</td></tr>", r.RequestURI))
	io.WriteString(w, fmt.Sprintf("<tr><td>Header</td><td>%s</td></tr>", r.Header))
	io.WriteString(w, fmt.Sprintf("<tr><td>Body</td><td>%s</td></tr>", r.Body))
	io.WriteString(w, "</table></body></html>")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeView)
	r.HandleFunc("/info", info)

	flag.Parse()
	var err error

	if *local != "" { // Run as a local web server
		err = http.ListenAndServe(*local, r)
	} else if *tcp != "" { // Run as FCGI via TCP
		listener, err := net.Listen("tcp", *tcp)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, r)
	} else if *unix != "" { // Run as FCGI via UNIX socket
		listener, err := net.Listen("unix", *unix)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, r)
	} else { // Run as FCGI via standard I/O
		err = fcgi.Serve(nil, r)
	}
	if err != nil {
		log.Fatal(err)
	}
}
