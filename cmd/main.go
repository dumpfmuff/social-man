package main

import (
	"flag"
	"fmt"
	"github.com/osimono/social-man/cmd/app/server"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var protocol = flag.String("protocol", "https", "protocol for the webserver")
	var host = flag.String("host", "localhost", "hostname for the webserver")
	var port = flag.Int("port", 443, "port for the webserver")
	flag.Parse()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)

	server.Init(*protocol, *host, *port)
}
