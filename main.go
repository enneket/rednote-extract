package main

import (
	"flag"
)

var (
	keyword = flag.String("k", "", "keyword to extract")
)

func main() {
	flag.Parse()
	if *keyword == "" {
		flag.Usage()
		return
	}
}
