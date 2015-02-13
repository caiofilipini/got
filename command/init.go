// Package command provides basic funcionality and implementation
// of the default commands.
package command

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
