package main

import "log"

func main() {
	log.SetFlags(0)

	customizedGOROOT, err := CreateCustomizedGOROOT()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("customizedGOROOT=%#+v", customizedGOROOT)
}
