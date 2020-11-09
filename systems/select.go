package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type SelectComponent struct {
	ACallback, BCallback func()
}

func (c *SelectComponent) GetSelectComponent() *SelectComponent {
	return c
}

type NotSelectComponent struct{}

func (c *NotSelectComponent) GetNotSelectComponent() *NotSelectComponent {
	return c
}

type NotSelectAble interface {
	GetNotSelectComponent() *NotSelectComponent
}

type SelectFace interface {
	GetSelectComponent() *SelectComponent
}

type SelectAble interface {
	common.BasicFace
	CursorFace
	SelectFace
}

type sceneSwitchEntity struct {
	*ecs.BasicEntity
	*CursorComponent
	*SelectComponent
}

type SelectSystem struct {
	entities []sceneSwitchEntity
}

func (s *SelectSystem) Add(basic *ecs.BasicEntity, cursor *CursorComponent, scene *SelectComponent) {
	s.entities = append(s.entities, sceneSwitchEntity{basic, cursor, scene})
}

func (s *SelectSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(SelectAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetCursorComponent(), o.GetSelectComponent())
}

func (s *SelectSystem) Remove(basic ecs.BasicEntity) {
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

func (s *SelectSystem) Update(dt float32) {
	for _, e := range s.entities {
		if engo.Input.Button("A").JustPressed() {
			if e.Selected {
				e.ACallback()
			}
		}
		if engo.Input.Button("B").JustPressed() {
			if e.Selected {
				e.BCallback()
			}
		}
	}
}
