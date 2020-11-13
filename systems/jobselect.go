package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type JobSelectComponent struct {
	Atk, Def, Spd int
	SpecialName   string
	SpecialURL    string
	shown         bool
}

func (c *JobSelectComponent) GetJobSelectComponent() *JobSelectComponent {
	return c
}

type NotJobSelectComponent struct{}

func (c *NotJobSelectComponent) GetNotJobSelectComponent() *NotJobSelectComponent {
	return c
}

type NotJobSelectAble interface {
	GetNotJobSelectComponent() *NotJobSelectComponent
}

type JobSelectFace interface {
	GetJobSelectComponent() *JobSelectComponent
}

type JobSelectAble interface {
	common.BasicFace
	CursorFace
	JobSelectFace
}

type jobSelectEntity struct {
	*ecs.BasicEntity
	*CursorComponent
	*JobSelectComponent
}

type JobSelectSystem struct {
	Fnt *common.Font

	entities []jobSelectEntity

	background                             sprite
	atkText, defText, spdText, spcText     sprite
	atkIcon0, defIcon0, spdIcon0, spcIcon0 sprite
	atkIcon1, defIcon1, spdIcon1, spcIcon1 sprite
	atkIcon2, defIcon2, spdIcon2, spcIcon2 sprite
	atkIcon3, defIcon3, spdIcon3, spcIcon3 sprite
}

func (s *JobSelectSystem) New(w *ecs.World) {
	s.background = sprite{BasicEntity: ecs.NewBasic()}
	s.background.Drawable, _ = common.LoadedSprite("start/jobSelect.png")
	s.background.SetZIndex(1)
	s.background.Position = engo.Point{X: 72, Y: 100}
	w.AddEntity(&s.background)

	s.atkText = sprite{BasicEntity: ecs.NewBasic()}
	s.atkText.Drawable = common.Text{
		Text: "ATK:",
		Font: s.Fnt,
	}
	s.atkText.SetZIndex(2)
	s.atkText.SetCenter(engo.Point{X: 240, Y: 125})
	s.atkText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&s.atkText)

	s.defText = sprite{BasicEntity: ecs.NewBasic()}
	s.defText.Drawable = common.Text{
		Text: "DEF:",
		Font: s.Fnt,
	}
	s.defText.SetZIndex(2)
	s.defText.SetCenter(engo.Point{X: 240, Y: 185})
	s.defText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&s.defText)

	s.spdText = sprite{BasicEntity: ecs.NewBasic()}
	s.spdText.Drawable = common.Text{
		Text: "SPD:",
		Font: s.Fnt,
	}
	s.spdText.SetZIndex(2)
	s.spdText.SetCenter(engo.Point{X: 240, Y: 245})
	s.spdText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&s.spdText)

	s.spcText = sprite{BasicEntity: ecs.NewBasic()}
	s.spcText.Drawable = common.Text{
		Text: "SPC:",
		Font: s.Fnt,
	}
	s.spcText.SetZIndex(2)
	s.spcText.SetCenter(engo.Point{X: 240, Y: 310})
	s.spcText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&s.spcText)

	s.atkIcon0 = sprite{BasicEntity: ecs.NewBasic()}
	s.atkIcon0.Drawable, _ = common.LoadedSprite("start/atk.png")
	s.atkIcon0.SetZIndex(2)
	s.atkIcon0.SetCenter(engo.Point{X: 360, Y: 125})
	w.AddEntity(&s.atkIcon0)

	s.atkIcon1 = sprite{BasicEntity: ecs.NewBasic()}
	s.atkIcon1.Drawable, _ = common.LoadedSprite("start/atk.png")
	s.atkIcon1.SetZIndex(2)
	s.atkIcon1.SetCenter(engo.Point{X: 390, Y: 125})
	w.AddEntity(&s.atkIcon1)

	s.atkIcon2 = sprite{BasicEntity: ecs.NewBasic()}
	s.atkIcon2.Drawable, _ = common.LoadedSprite("start/atk.png")
	s.atkIcon2.SetZIndex(2)
	s.atkIcon2.SetCenter(engo.Point{X: 420, Y: 125})
	w.AddEntity(&s.atkIcon2)

	s.atkIcon3 = sprite{BasicEntity: ecs.NewBasic()}
	s.atkIcon3.Drawable, _ = common.LoadedSprite("start/atk.png")
	s.atkIcon3.SetZIndex(2)
	s.atkIcon3.SetCenter(engo.Point{X: 450, Y: 125})
	w.AddEntity(&s.atkIcon3)

	s.defIcon0 = sprite{BasicEntity: ecs.NewBasic()}
	s.defIcon0.Drawable, _ = common.LoadedSprite("start/def.png")
	s.defIcon0.SetZIndex(2)
	s.defIcon0.SetCenter(engo.Point{X: 360, Y: 185})
	w.AddEntity(&s.defIcon0)

	s.defIcon1 = sprite{BasicEntity: ecs.NewBasic()}
	s.defIcon1.Drawable, _ = common.LoadedSprite("start/def.png")
	s.defIcon1.SetZIndex(2)
	s.defIcon1.SetCenter(engo.Point{X: 390, Y: 185})
	w.AddEntity(&s.defIcon1)

	//def2
	s.defIcon2 = sprite{BasicEntity: ecs.NewBasic()}
	s.defIcon2.Drawable, _ = common.LoadedSprite("start/def.png")
	s.defIcon2.SetZIndex(2)
	s.defIcon2.SetCenter(engo.Point{X: 420, Y: 185})
	w.AddEntity(&s.defIcon2)

	//def3
	s.defIcon3 = sprite{BasicEntity: ecs.NewBasic()}
	s.defIcon3.Drawable, _ = common.LoadedSprite("start/def.png")
	s.defIcon3.SetZIndex(2)
	s.defIcon3.SetCenter(engo.Point{X: 450, Y: 185})
	w.AddEntity(&s.defIcon3)

	s.spdIcon0 = sprite{BasicEntity: ecs.NewBasic()}
	s.spdIcon0.Drawable, _ = common.LoadedSprite("start/spd.png")
	s.spdIcon0.SetZIndex(2)
	s.spdIcon0.SetCenter(engo.Point{X: 360, Y: 245})
	w.AddEntity(&s.spdIcon0)

	s.spdIcon1 = sprite{BasicEntity: ecs.NewBasic()}
	s.spdIcon1.Drawable, _ = common.LoadedSprite("start/spd.png")
	s.spdIcon1.SetZIndex(2)
	s.spdIcon1.SetCenter(engo.Point{X: 390, Y: 245})
	w.AddEntity(&s.spdIcon1)

	s.spdIcon2 = sprite{BasicEntity: ecs.NewBasic()}
	s.spdIcon2.Drawable, _ = common.LoadedSprite("start/spd.png")
	s.spdIcon2.SetZIndex(2)
	s.spdIcon2.SetCenter(engo.Point{X: 420, Y: 245})
	w.AddEntity(&s.spdIcon2)

	s.spdIcon3 = sprite{BasicEntity: ecs.NewBasic()}
	s.spdIcon3.Drawable, _ = common.LoadedSprite("start/spd.png")
	s.spdIcon3.SetZIndex(2)
	s.spdIcon3.SetCenter(engo.Point{X: 450, Y: 245})
	w.AddEntity(&s.spdIcon3)

	s.spcIcon0 = sprite{BasicEntity: ecs.NewBasic()}
	s.spcIcon0.Drawable, _ = common.LoadedSprite("start/spc.png")
	s.spcIcon0.SetZIndex(2)
	s.spcIcon0.SetCenter(engo.Point{X: 360, Y: 310})
	w.AddEntity(&s.spcIcon0)

	s.spcIcon1 = sprite{BasicEntity: ecs.NewBasic()}
	s.spcIcon1.Drawable, _ = common.LoadedSprite("start/spc.png")
	s.spcIcon1.SetZIndex(2)
	s.spcIcon1.SetCenter(engo.Point{X: 390, Y: 310})
	w.AddEntity(&s.spcIcon1)

	s.spcIcon2 = sprite{BasicEntity: ecs.NewBasic()}
	s.spcIcon2.Drawable, _ = common.LoadedSprite("start/spc.png")
	s.spcIcon2.SetZIndex(2)
	s.spcIcon2.SetCenter(engo.Point{X: 420, Y: 310})
	w.AddEntity(&s.spcIcon2)

	s.spcIcon3 = sprite{BasicEntity: ecs.NewBasic()}
	s.spcIcon3.Drawable, _ = common.LoadedSprite("start/spc.png")
	s.spcIcon3.SetZIndex(2)
	s.spcIcon3.SetCenter(engo.Point{X: 450, Y: 310})
	w.AddEntity(&s.spcIcon3)
}

func (s *JobSelectSystem) Add(basic *ecs.BasicEntity, cursor *CursorComponent, sel *JobSelectComponent) {
	s.entities = append(s.entities, jobSelectEntity{basic, cursor, sel})
}

func (s *JobSelectSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(JobSelectAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetCursorComponent(), o.GetJobSelectComponent())
}

func (s *JobSelectSystem) Remove(basic ecs.BasicEntity) {
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

func (s *JobSelectSystem) Update(dt float32) {
	for _, e := range s.entities {
		if e.Selected {
			if !e.shown {
				s.hideAllButBG()
				s.show(e)
				e.shown = true
			}
		}
	}
}

func (s *JobSelectSystem) hideAll() {
	s.hideAllButBG()
	s.background.Hidden = true
	s.atkText.Hidden = true
	s.defText.Hidden = true
	s.spdText.Hidden = true
	s.spcText.Hidden = true
}

func (s *JobSelectSystem) hideAllButBG() {
	s.atkIcon0.Hidden = true
	s.atkIcon1.Hidden = true
	s.atkIcon2.Hidden = true
	s.atkIcon3.Hidden = true
	s.defIcon0.Hidden = true
	s.defIcon1.Hidden = true
	s.defIcon2.Hidden = true
	s.defIcon3.Hidden = true
	s.spdIcon0.Hidden = true
	s.spdIcon1.Hidden = true
	s.spdIcon2.Hidden = true
	s.spdIcon3.Hidden = true
	s.spcIcon0.Hidden = true
	s.spcIcon1.Hidden = true
	s.spcIcon2.Hidden = true
	s.spcIcon3.Hidden = true
}

func (s *JobSelectSystem) show(e jobSelectEntity) {}
