package constants

import (
	"github.com/dropbox/godropbox/container/set"
)

const (
	Solid  = "solid"
	Slow   = "slow"
	Medium = "medium"
	Fast   = "fast"
)

var Modes = set.NewSet(Solid, Slow, Medium, Fast)
