package main

import (
	"github.com/Liam-E2/mainevent/eventsource"
	"os"
	"log"
	"fmt"
)

func main(){
	// Run Router
	host, specified := os.LookupEnv("EVENTSOURCEHOST")
	if !specified {
		host = "localhost"
	}
	port, specified := os.LookupEnv("EVENTSOURCEPORT")
	if !specified {
		port = "9019"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Fatal(eventsource.RunServer(addr))
}
