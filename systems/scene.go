package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type SceneSwitchComponent struct {
	//scene to switch to
	To       string
	NewWorld bool
}

func (c *SceneSwitchComponent) GetSceneSwitchComponent() *SceneSwitchComponent {
	return c
}

type NotSceneSwitchComponent struct{}

func (c *NotSceneSwitchComponent) GetNotSceneSwitchComponent() *NotSceneSwitchComponent {
	return c
}

type NotSceneSwitchAble interface {
	GetNotSceneSwitchComponent() *NotSceneSwitchComponent
}

type SceneSwitchFace interface {
	GetSceneSwitchComponent() *SceneSwitchComponent
}

type SceneSwitchAble interface {
	common.BasicFace
	CursorFace
	SceneSwitchFace
}

type sceneSwitchEntity struct {
	*ecs.BasicEntity
	*CursorComponent
	*SceneSwitchComponent
}

type SceneSwitchSystem struct {
	entities []sceneSwitchEntity
}

func (s *SceneSwitchSystem) Add(basic *ecs.BasicEntity, cursor *CursorComponent, scene *SceneSwitchComponent) {
	s.entities = append(s.entities, sceneSwitchEntity{basic, cursor, scene})
}

func (s *SceneSwitchSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(SceneSwitchAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetCursorComponent(), o.GetSceneSwitchComponent())
}

func (s *SceneSwitchSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, entity := range s.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *SceneSwitchSystem) Update(dt float32) {
	for _, e := range s.entities {
		if engo.Input.Button("A").JustPressed() {
			if e.Selected {
				engo.SetSceneByName(e.To, e.NewWorld)
			}
		}
	}
}
