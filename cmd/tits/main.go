package main

import (
	"github.com/boobsrate/core/internal/entrypoint/tits"
)

func main() {
	err := tits.Run()
	if err != nil {
		panic(err)
	}
}
