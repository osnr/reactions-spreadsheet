// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

var addr = flag.String("addr", ":9391", "http service address")
var fs = http.FileServer(http.Dir("./static"))

func serve(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, r.URL)
	if strings.HasPrefix(r.URL.Path, "/static/") {
		r.URL.Path = strings.Replace(r.URL.Path, "/static/", "/", 1)
		fs.ServeHTTP(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// any /whatever, /blahblah, /etc URL is a spreadsheet.
	http.ServeFile(w, r, "spreadsheet.html")
}

func main() {
	flag.Parse()

	hubs := make(map[string]*Hub)

	http.HandleFunc("/", serve)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub := hubs[r.URL.RawQuery]
		if hub == nil {
			hub = newHub()
			go hub.run()
			hubs[r.URL.RawQuery] = hub
		}
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
