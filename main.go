package main

import (
	"mymodule/mypackage"

	"github.com/faiface/pixel/pixelgl"
)

func run() {
	mypackage.CreateGame()
}

func main() {
	pixelgl.Run(run)
}
