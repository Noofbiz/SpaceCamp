package systems

import (
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type cursorChangeMessage struct{}

func (cursorChangeMessage) Type() string {
	return "Cursor Change Message"
}

type CursorSetMessage struct {
	Index int
}

func (CursorSetMessage) Type() string {
	return "Cursor Set Message"
}

type CursorEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.RenderComponent
	*CursorComponent
}

type pointer struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type CursorComponent struct {
	Selected bool
}

func (c *CursorComponent) GetCursorComponent() *CursorComponent {
	return c
}

type NotCursorComponent struct{}

func (n *NotCursorComponent) GetNotCursorComponent() *NotCursorComponent {
	return n
}

type CursorSystem struct {
	entities []CursorEntity
	ptr      pointer
}

func (s *CursorSystem) New(w *ecs.World) {
	s.ptr = pointer{BasicEntity: ecs.NewBasic()}
	pointerTex, err := common.LoadedSprite("title/cursor.png")
	if err != nil {
		log.Fatalf("Unable to load pointer.png Error was: %v", err)
	}
	s.ptr.Drawable = pointerTex
	s.ptr.Hidden = true
	s.ptr.Width = s.ptr.Drawable.Width()
	s.ptr.Height = s.ptr.Drawable.Height()
	s.ptr.SetZIndex(100)
	w.AddEntity(&s.ptr)
	engo.Mailbox.Listen("Cursor Set Message", func(msg engo.Message) {
		m, ok := msg.(CursorSetMessage)
		if !ok {
			return
		}
		s.setPointer(m.Index)
	})
}

func (s *CursorSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent, render *common.RenderComponent, selection *CursorComponent) {
	s.entities = append(s.entities, CursorEntity{basic, space, render, selection})
	if selection.Selected {
		s.setPointer(len(s.entities) - 1)
	}
}

func (s *CursorSystem) AddByInterface(id ecs.Identifier) {
	o, ok := id.(CursorAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetSpaceComponent(), o.GetRenderComponent(), o.GetCursorComponent())
}

func (s *CursorSystem) Remove(basic ecs.BasicEntity) {
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

func (s *CursorSystem) Update(dt float32) {
	for i := 0; i < len(s.entities); i++ {
		if s.entities[i].Selected == true {
			if engo.Input.Button("down").JustPressed() || engo.Input.Button("right").JustPressed() {
				s.entities[i].Selected = false
				if i+1 >= len(s.entities) && len(s.entities) > 0 {
					s.entities[i].Selected = true
					s.setPointer(i)
					return
				} else {
					s.entities[i+1].Selected = true
					s.setPointer(i + 1)
					return
				}
			} else if engo.Input.Button("up").JustPressed() || engo.Input.Button("left").JustPressed() {
				s.entities[i].Selected = false
				if i-1 < 0 && len(s.entities) > 0 {
					s.entities[i].Selected = true
					s.setPointer(i)
					return
				} else {
					s.entities[i-1].Selected = true
					s.setPointer(i - 1)
					return
				}
			}
		}
	}
}

func (s *CursorSystem) setPointer(i int) {
	if len(s.entities) == 0 || i > len(s.entities)-1 || i < 0 {
		s.ptr.Hidden = true
		return
	}
	ent := s.entities[i]
	s.ptr.Hidden = false
	s.ptr.Position.X = ent.Position.X - s.ptr.Width - 2
	s.ptr.Position.Y = ent.Position.Y + (ent.Height / 2) + 6
}

type CursorFace interface {
	GetCursorComponent() *CursorComponent
}

type CursorAble interface {
	common.BasicFace
	common.SpaceFace
	common.RenderFace
	CursorFace
}

type NotCursorAble interface {
	GetNotCursorComponent() *NotCursorComponent
}
