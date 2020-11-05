package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
)

type FullScreenSystem struct{}

func (*FullScreenSystem) Remove(basic ecs.BasicEntity) {}

func (f *FullScreenSystem) Update(float32) {
	if engo.Input.Button("FullScreen").JustPressed() {
		setFullScreenImpl()
	}
}
