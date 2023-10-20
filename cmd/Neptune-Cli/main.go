package main

import (
	"github.com/m1ndo/Neptune/pkg/neptune"
)

func main() {
	App := neptune.NewApp()
	App.Sys.Init()
}
