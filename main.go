package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ipInfo struct {
	IP      string
	Updated time.Time
}

var ips = make(map[string]ipInfo)

func main() {
	addr := "localhost:42514"
	log.Printf("Using %s", addr)

	http.HandleFunc("/", handle)
	error := http.ListenAndServe(addr, nil)

	if error != nil {
		log.Fatal(error)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]

	switch r.Method {
	case "GET":
		get(name, w)
	case "PUT":
		put(name, w, r)
	default:
		http.Error(w, "Unsupported method", 405)
	}
}

func get(name string, w http.ResponseWriter) {
	if info, ok := ips[name]; ok {
		w.Header().Set("X-Updated", info.Updated.String())
		fmt.Fprintf(w, info.IP)

		log.Printf("Returned %s = %s", name, info.IP)
	} else {
		http.Error(w, fmt.Sprintf("%s not found", name), 404)

		log.Printf("Returned %s = not found", name)
	}
}

func put(name string, w http.ResponseWriter, r *http.Request) {
	if len(name) == 0 {
		http.Error(w, "Empty names not allowed", 400)
		return
	}

	info := ipInfo{
		IP:      requestIP(r),
		Updated: time.Now().UTC(),
	}
	ips[name] = info

	w.WriteHeader(http.StatusCreated)

	log.Printf("Set %s = %s", name, info.IP)
}

func requestIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); len(xff) > 0 {
		return xff
	}

	return strings.Split(r.RemoteAddr, ":")[0]
}
