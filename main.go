package main

import (
	"github.com/Noofbiz/SpaceCamp/savedata"
	"github.com/Noofbiz/SpaceCamp/scenes"

	"github.com/EngoEngine/engo"
)

func main() {
	savedata.CurrentSave.Load()
	title := scenes.TitleScene{}
	engo.RegisterScene(&title)
	// options scene
	// credits scene
	engo.RegisterScene(&scenes.NewGameScene{})
	engo.RegisterScene(&scenes.ShopScene{})
	// pre-takeoff ship scene
	engo.RegisterScene(&scenes.TakeOffScene{})
	// loss scene
	// aftermath ship scene
	engo.Run(engo.RunOptions{
		Title:                      "To INFINITY!!!",
		Width:                      640, //512, //16
		Height:                     360, //288, //9
		ScaleOnResize:              true,
		FPSLimit:                   60,
		ApplicationMajorVersion:    0,
		ApplicationMinorVersion:    1,
		ApplicationRevisionVersion: 0,
		// }, &title)
	}, &title)
}
