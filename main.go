package main

import (
	"flag"
	"github.com/zachhuff386/hue-alert/cmd"
)

func main() {
	flag.Parse()

	var err error

	switch flag.Arg(0) {
	case "google":
		err = cmd.Google()
	case "accounts":
		err = cmd.Accounts()
	case "remove":
		err = cmd.Remove(flag.Arg(1), flag.Arg(2))
	}

	if err != nil {
		panic(err)
	}
}
