//+build js

package systems

import (
	"syscall/js"

	"github.com/EngoEngine/engo"
)

func setFullScreenImpl() {
	doc := js.Global().Get("document")
	window := js.Global().Get("window")
	body := doc.Get("body")
	canvas := body.Call("getElementsByTagName", "canvas").Index(0)
	if fs := canvas.Get("webkitRequestFullScreen"); fs.Truthy() {
		canvas.Call("webkitRequestFullScreen")
	} else {
		canvas.Call("mozRequestFullScreen")
	}
	newW := window.Get("innerWidth").Float()
	newH := window.Get("innerHeight").Float()
	canvas.Set("width", newW)
	canvas.Set("height", newH)
	engo.Gl.Viewport(0, 0, int(newW), int(newH))
}
