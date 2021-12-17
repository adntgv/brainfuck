package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	filename := "hello.bf"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	processor := NewProcessor()

	bz, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	for _, b := range bz {
		err := processor.process(b, true)
		if err != nil {
			log.Fatal(err)
		}
	}
}
