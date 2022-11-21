package main

import (
		"fmt"
	  "os"

		log "github.com/sirupsen/logrus"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage:\n\n")
		fmt.Printf("        %s <command>\n\n", os.Args[0])
		fmt.Printf("The commands are:\n\n")
		fmt.Printf("        geo          generate hotel geographical location data\n")
		fmt.Printf("        hotels       generate hotel profile data\n")
		fmt.Printf("        inventory    generate hotel inventory data\n")
		os.Exit(-1)
	}

	var cmd = os.Args[1]

	f := os.Stdout

	switch cmd {
	case "geo":
		GenerateGeo(f)
	case "hotels":
		GenerateHotels(f)
	case "inventory":
		GenerateInventory(f)
	default:
		log.Fatalf("unknown cmd: %s", cmd)
	}
}
