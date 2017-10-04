package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/google/gopacket/dumpcommand"
	"github.com/google/gopacket/pcapgo"
)

func main() {
	flag.Parse()
	err := dump(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
}

func dump(r io.Reader) error {
	pr, err := pcapgo.NewReader(r)
	if err != nil {
		return err
	}
	dumpcommand.Run(pr)
	return nil
}
