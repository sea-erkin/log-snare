package main

import (
	"flag"
	"log"
	"log-snare/web/server"
)

// cli
var (
	configFileFlag = flag.String("c", "logs-snare-web.yml", "path to config yaml file")
	debugFileFlag  = flag.Bool("d", false, "debug")
	resetDbFlag    = flag.Bool("r", false, "reset db")
	listenHost     = flag.String("p", "0.0.0.0:8080", "listen host")
)

func main() {
	flag.Parse()
	log.Fatal(server.Run(*configFileFlag, *debugFileFlag, *resetDbFlag, *listenHost))
}
