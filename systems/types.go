package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type sprite struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type selection struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
	CursorComponent
}
