package systems

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// ExitSystem exits the game when you press esc for 3 seconds.
type ExitSystem struct {
	f      *common.Font
	entity struct {
		*ecs.BasicEntity
		*common.RenderComponent
		*common.SpaceComponent
	}
	time float32
}

type warn struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (e *ExitSystem) New(w *ecs.World) {
	e.f = &common.Font{
		Size: 12,
		FG:   color.Black,
		URL:  "title/PressStart.ttf",
	}
	e.f.CreatePreloaded()

	warning := warn{BasicEntity: ecs.NewBasic()}
	warning.SpaceComponent = common.SpaceComponent{
		Width:  15,
		Height: 15,
	}
	warning.RenderComponent = common.RenderComponent{
		Drawable: common.Text{
			Font: e.f,
			Text: "exiting",
		},
		Hidden:      true,
		StartZIndex: 5,
	}

	w.AddEntity(&warning)
	e.entity.BasicEntity = &warning.BasicEntity
	e.entity.SpaceComponent = &warning.SpaceComponent
	e.entity.RenderComponent = &warning.RenderComponent
}

func (e *ExitSystem) Remove(basic ecs.BasicEntity) {}

func (e *ExitSystem) Update(dt float32) {
	if engo.Input.Button("Exit").Down() {
		e.entity.Hidden = false
		e.time += dt
		if e.time < 0.3 {
			e.entity.Drawable = common.Text{
				Text: "exiting .",
				Font: e.f,
			}
		} else if e.time < 0.7 {
			e.entity.Drawable = common.Text{
				Text: "exiting . .",
				Font: e.f,
			}
		} else if e.time < 1 {
			e.entity.Drawable = common.Text{
				Text: "exiting . . .",
				Font: e.f,
			}
		} else {
			engo.Exit()
		}
	} else {
		e.entity.Drawable = common.Text{
			Text: "exiting",
			Font: e.f,
		}
		e.entity.Hidden = true
		e.time = 0
	}
}
