package main

import (
	"os"
)

func main() {
	for _, f := range os.Args[1:] {
		rd, err := os.Open(f)
		if err != nil {
			panic(err)
		}

		p := NewParser(NewTokeniser(rd))

		if err := p.parse(); err != nil {
			panic(err)
		}

		if _, err := NewFormatter(os.Stdout).Format(p); err != nil {
			panic(err)
		}
	}
}
