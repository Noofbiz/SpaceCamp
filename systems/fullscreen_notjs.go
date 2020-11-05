//+build !js

package systems

import (
	"github.com/EngoEngine/engo"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func setFullScreenImpl() {
	monitor := glfw.GetPrimaryMonitor()
	var mode *glfw.VidMode
	if monitor != nil {
		mode = monitor.GetVideoMode()
	} else {
		// Initialize default values if no monitor is found
		mode = &glfw.VidMode{
			Width:       1,
			Height:      1,
			RedBits:     8,
			GreenBits:   8,
			BlueBits:    8,
			RefreshRate: 60,
		}
	}
	engo.Window.SetAttrib(glfw.Decorated, 0)
	engo.Window.SetSize(mode.Width, mode.Height)
	engo.Window.SetPos(0, 0)
}
