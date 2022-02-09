package main

import "github.com/boobsrate/core/internal/applications/tits"

func main() {
	err := tits.Run()
	if err != nil {
		panic(err)
	}
}
