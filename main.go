package main

import (
	"flag"
	"log"
	"net/http"
	"webshell/lib/httpHandlers"
)

func main() {
	httpHandlers.RequiredPassword = flag.String("pass", "mypass", "passwd")
	addr := flag.String("addr", ":9090", "bind addr and port")

	flag.Parse()

	http.HandleFunc("/", httpHandlers.Handler)

	log.Printf("start http server\n")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// curl -X POST -d '{"pass": "mypass", "cmd": "bash", "opt": "-c", "args": "ls -l ~; echo hello"}' http://localhost:9090
