package domain

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// NewID returns a new identifier.
func NewID() string { return <-ids }

// A channel which yields ids.
var ids = make(chan string, 1)

func init() {
	rnd := ulid.Monotonic(
		rand.New(
			rand.NewSource(time.Now().UnixNano()),
		),
		0,
	)

	go func() {
		for {
			ids <- ulid.MustNew(ulid.Now(), rnd).String()
		}
	}()
}

func IsID(id string) bool {
	_, err := ulid.Parse(id)
	return err == nil
}
