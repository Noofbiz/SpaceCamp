package main

import (
	"github.com/Noofbiz/SpaceCamp/scenes"

	"github.com/EngoEngine/engo"
)

func main() {
	title := scenes.TitleScene{}
	engo.RegisterScene(&title)
	engo.RegisterScene(&scenes.NewGameScene{})
	engo.Run(engo.RunOptions{
		Title:                      "To INFINITY!!!",
		Width:                      640, //512, //16
		Height:                     360, //288, //9
		ScaleOnResize:              true,
		FPSLimit:                   60,
		ApplicationMajorVersion:    0,
		ApplicationMinorVersion:    1,
		ApplicationRevisionVersion: 0,
	}, &title)
}
