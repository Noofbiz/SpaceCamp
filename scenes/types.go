package scenes

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"

	"github.com/Noofbiz/SpaceCamp/systems"
)

type sprite struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type audio struct {
	ecs.BasicEntity
	common.AudioComponent
}

type playerSelectableSprite struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type selection struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	systems.CursorComponent
	systems.SelectComponent
}
