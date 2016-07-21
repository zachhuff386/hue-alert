package main

import (
	"github.com/zachhuff386/hue-alert/cmd"
)

func main() {
	err := cmd.Alert()
	if err != nil {
		panic(err)
	}
}
