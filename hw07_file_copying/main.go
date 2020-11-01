package main

import (
	"flag"
	"fmt"
	"os"
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

const helpText = `
Usage: go-cp [OPTION]... -from SOURCE -to DEST
Copy SOURCE file to DEST file.

Options:
- to 		(mandatory)	path of toination file
- from 		(mandatory)	path of source file to copy
- limit 	(optional)	maximum bytes to copy
- offset 	(optional)	offset in source file

Examples:

With minimum options
	cp -from /tmp/from.txt -to /tmp/to.txt

With maximum options
	cp -from /tmp/from.txt -to /tmp/to.txt -offset 10 -limit 5`

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println(helpText)
		os.Exit(0)
	}

	if from == "" {
		fmt.Println("Source file is missing. Please specify it with -from option.")
		os.Exit(1)
	}

	if to == "" {
		fmt.Println("Destination file is missing. Please specify it with -to option.")
		os.Exit(1)
	}

	if err := Copy(from, to, int64(offset), int64(limit)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
