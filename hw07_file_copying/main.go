package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if from == "" || to == "" || from == to {
		flag.Usage()
		log.Fatal("Invalid arguments: -from and -to must be specified and not be the same.")
	}

	fmt.Printf("Copying from %s to %s with offset %d and limit %d\n", from, to, offset, limit)

	err := Copy(from, to, offset, limit)
	if err != nil {
		log.Fatalf("copy files: %v", err)
	}

	fmt.Printf("File copied to %s successfully\n", to)
}
