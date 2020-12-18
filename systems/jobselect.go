package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type JobSelectPauseMessage struct{}

func (JobSelectPauseMessage) Type() string {
	return "Job Select Pause Message"
}

type JobSelectUnpauseMessage struct {
	IDs []*ecs.BasicEntity
}

func (JobSelectUnpauseMessage) Type() string {
	return "Job Select Unpause Message"
}

type JobSelectComponent struct {
	Atk, Def, Spd int
	Job, Name     string
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
	JobSelectFace
}

type jobSelectEntity struct {
	*ecs.BasicEntity
	*CursorComponent
	*JobSelectComponent
}

type JobSelectSystem struct {
	Fnt    *common.Font
	LogSnd *common.Player

	entities []jobSelectEntity

	background                             sprite
	atkText, defText, spdText, spcText     sprite
	atkIcon0, defIcon0, spdIcon0, spcIcon0 sprite
	atkIcon1, defIcon1, spdIcon1, spcIcon1 sprite
	atkIcon2, defIcon2, spdIcon2, spcIcon2 sprite
	atkIcon3, defIcon3, spdIcon3, spcIcon3 sprite
	titles                                 []*selection
	positions                              []engo.Point

	paused, skipNextFrame bool

	curSys    *CursorSystem
	acceptSys *AcceptSystem

	w *ecs.World
}

func (s *JobSelectSystem) pause() {
	s.paused = true
	s.hideAll()
	for _, title := range s.titles {
		s.curSys.Remove(title.BasicEntity)
	}
}

func (s *JobSelectSystem) unpause(ids []*ecs.BasicEntity) {
	s.paused = false
	s.skipNextFrame = true
	s.showAll()
	for _, title := range s.titles {
		s.curSys.Add(&title.BasicEntity, &title.SpaceComponent, &title.RenderComponent, &title.CursorComponent)
	}
	for _, id := range ids {
		s.curSys.Remove(*id)
	}
	engo.Mailbox.Dispatch(CursorSetMessage{ID: s.titles[0]})
}

func (s *JobSelectSystem) New(w *ecs.World) {
	s.w = w

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CursorSystem:
			s.curSys = sys
		case *AcceptSystem:
			s.acceptSys = sys
		}
	}

	engo.Mailbox.Listen("Job Select Pause Message", func(m engo.Message) {
		_, ok := m.(JobSelectPauseMessage)
		if !ok {
			return
		}
		s.pause()
	})

	engo.Mailbox.Listen("Job Select Unpause Message", func(m engo.Message) {
		msg, ok := m.(JobSelectUnpauseMessage)
		if !ok {
			return
		}
		s.unpause(msg.IDs)
	})

	s.positions = []engo.Point{
		engo.Point{X: 130, Y: 110},
		engo.Point{X: 130, Y: 150},
		engo.Point{X: 130, Y: 190},
		engo.Point{X: 130, Y: 235},
		engo.Point{X: 130, Y: 275},
		engo.Point{X: 130, Y: 315},
	}

	s.background = sprite{BasicEntity: ecs.NewBasic()}
	s.background.Drawable, _ = common.LoadedSprite("start/jobSelect.png")
	s.background.SetZIndex(1)
	s.background.Position = engo.Point{X: 71, Y: 90}
	w.AddEntity(&s.background)

	s.atkText = sprite{BasicEntity: ecs.NewBasic()}
	s.atkText.Drawable = common.Text{
		Text: "ATK:",
		Font: s.Fnt,
	}
	s.atkText.SetZIndex(2)
	s.atkText.SetCenter(engo.Point{X: 240, Y: 100})
	s.atkText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&s.atkText)

	s.defText = sprite{BasicEntity: ecs.NewBasic()}
	s.defText.Drawable = common.Text{
		Text: "DEF:",
		Font: s.Fnt,
	}
	s.defText.SetZIndex(2)
	s.defText.SetCenter(engo.Point{X: 240, Y: 160})
	s.defText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&s.defText)

	s.spdText = sprite{BasicEntity: ecs.NewBasic()}
	s.spdText.Drawable = common.Text{
		Text: "SPD:",
		Font: s.Fnt,
	}
	s.spdText.SetZIndex(2)
	s.spdText.SetCenter(engo.Point{X: 240, Y: 220})
	s.spdText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&s.spdText)

	s.spcText = sprite{BasicEntity: ecs.NewBasic()}
	s.spcText.Drawable = common.Text{
		Text: "SPC:",
		Font: s.Fnt,
	}
	s.spcText.SetZIndex(2)
	s.spcText.SetCenter(engo.Point{X: 240, Y: 285})
	s.spcText.Scale = engo.Point{X: 0.5, Y: 0.5}
	w.AddEntity(&s.spcText)

	s.atkIcon0 = sprite{BasicEntity: ecs.NewBasic()}
	s.atkIcon0.Drawable, _ = common.LoadedSprite("start/atk.png")
	s.atkIcon0.SetZIndex(2)
	s.atkIcon0.SetCenter(engo.Point{X: 360, Y: 100})
	w.AddEntity(&s.atkIcon0)

	s.atkIcon1 = sprite{BasicEntity: ecs.NewBasic()}
	s.atkIcon1.Drawable, _ = common.LoadedSprite("start/atk.png")
	s.atkIcon1.SetZIndex(2)
	s.atkIcon1.SetCenter(engo.Point{X: 390, Y: 100})
	w.AddEntity(&s.atkIcon1)

	s.atkIcon2 = sprite{BasicEntity: ecs.NewBasic()}
	s.atkIcon2.Drawable, _ = common.LoadedSprite("start/atk.png")
	s.atkIcon2.SetZIndex(2)
	s.atkIcon2.SetCenter(engo.Point{X: 420, Y: 100})
	w.AddEntity(&s.atkIcon2)

	s.atkIcon3 = sprite{BasicEntity: ecs.NewBasic()}
	s.atkIcon3.Drawable, _ = common.LoadedSprite("start/atk.png")
	s.atkIcon3.SetZIndex(2)
	s.atkIcon3.SetCenter(engo.Point{X: 450, Y: 100})
	w.AddEntity(&s.atkIcon3)

	s.defIcon0 = sprite{BasicEntity: ecs.NewBasic()}
	s.defIcon0.Drawable, _ = common.LoadedSprite("start/def.png")
	s.defIcon0.SetZIndex(2)
	s.defIcon0.SetCenter(engo.Point{X: 360, Y: 160})
	w.AddEntity(&s.defIcon0)

	s.defIcon1 = sprite{BasicEntity: ecs.NewBasic()}
	s.defIcon1.Drawable, _ = common.LoadedSprite("start/def.png")
	s.defIcon1.SetZIndex(2)
	s.defIcon1.SetCenter(engo.Point{X: 390, Y: 160})
	w.AddEntity(&s.defIcon1)

	//def2
	s.defIcon2 = sprite{BasicEntity: ecs.NewBasic()}
	s.defIcon2.Drawable, _ = common.LoadedSprite("start/def.png")
	s.defIcon2.SetZIndex(2)
	s.defIcon2.SetCenter(engo.Point{X: 420, Y: 160})
	w.AddEntity(&s.defIcon2)

	//def3
	s.defIcon3 = sprite{BasicEntity: ecs.NewBasic()}
	s.defIcon3.Drawable, _ = common.LoadedSprite("start/def.png")
	s.defIcon3.SetZIndex(2)
	s.defIcon3.SetCenter(engo.Point{X: 450, Y: 160})
	w.AddEntity(&s.defIcon3)

	s.spdIcon0 = sprite{BasicEntity: ecs.NewBasic()}
	s.spdIcon0.Drawable, _ = common.LoadedSprite("start/spd.png")
	s.spdIcon0.SetZIndex(2)
	s.spdIcon0.SetCenter(engo.Point{X: 360, Y: 220})
	w.AddEntity(&s.spdIcon0)

	s.spdIcon1 = sprite{BasicEntity: ecs.NewBasic()}
	s.spdIcon1.Drawable, _ = common.LoadedSprite("start/spd.png")
	s.spdIcon1.SetZIndex(2)
	s.spdIcon1.SetCenter(engo.Point{X: 390, Y: 220})
	w.AddEntity(&s.spdIcon1)

	s.spdIcon2 = sprite{BasicEntity: ecs.NewBasic()}
	s.spdIcon2.Drawable, _ = common.LoadedSprite("start/spd.png")
	s.spdIcon2.SetZIndex(2)
	s.spdIcon2.SetCenter(engo.Point{X: 420, Y: 220})
	w.AddEntity(&s.spdIcon2)

	s.spdIcon3 = sprite{BasicEntity: ecs.NewBasic()}
	s.spdIcon3.Drawable, _ = common.LoadedSprite("start/spd.png")
	s.spdIcon3.SetZIndex(2)
	s.spdIcon3.SetCenter(engo.Point{X: 450, Y: 220})
	w.AddEntity(&s.spdIcon3)

	s.spcIcon0 = sprite{BasicEntity: ecs.NewBasic()}
	s.spcIcon0.Drawable, _ = common.LoadedSprite("start/spc.png")
	s.spcIcon0.SetZIndex(2)
	s.spcIcon0.SetCenter(engo.Point{X: 360, Y: 285})
	w.AddEntity(&s.spcIcon0)

	s.spcIcon1 = sprite{BasicEntity: ecs.NewBasic()}
	s.spcIcon1.Drawable, _ = common.LoadedSprite("start/spc.png")
	s.spcIcon1.SetZIndex(2)
	s.spcIcon1.SetCenter(engo.Point{X: 390, Y: 285})
	w.AddEntity(&s.spcIcon1)

	s.spcIcon2 = sprite{BasicEntity: ecs.NewBasic()}
	s.spcIcon2.Drawable, _ = common.LoadedSprite("start/spc.png")
	s.spcIcon2.SetZIndex(2)
	s.spcIcon2.SetCenter(engo.Point{X: 420, Y: 285})
	w.AddEntity(&s.spcIcon2)

	s.spcIcon3 = sprite{BasicEntity: ecs.NewBasic()}
	s.spcIcon3.Drawable, _ = common.LoadedSprite("start/spc.png")
	s.spcIcon3.SetZIndex(2)
	s.spcIcon3.SetCenter(engo.Point{X: 450, Y: 285})
	w.AddEntity(&s.spcIcon3)
}

func (s *JobSelectSystem) Add(basic *ecs.BasicEntity, sel *JobSelectComponent) {
	e := selection{BasicEntity: ecs.NewBasic()}
	e.Drawable = common.Text{
		Text: sel.Job,
		Font: s.Fnt,
	}
	e.SetZIndex(2)
	e.SetCenter(s.positions[len(s.entities)])
	e.Scale = engo.Point{X: 0.25, Y: 0.25}
	if len(s.titles) == 0 {
		e.Selected = true
	}
	s.w.AddEntity(&e)
	s.titles = append(s.titles, &e)
	s.entities = append(s.entities, jobSelectEntity{basic, &e.CursorComponent, sel})
}

func (s *JobSelectSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(JobSelectAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetJobSelectComponent())
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
	if s.paused {
		return
	}
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	for _, e := range s.entities {
		if e.Selected {
			if !e.shown {
				s.hideAllButBG()
				s.show(e)
				e.shown = true
			}
			if engo.Input.Button("A").JustPressed() {
				if !s.logDone() {
					continue
				}
				s.pause()
				msgs := []string{
					"So a " + e.Job,
					"Is that correct?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  s.Fnt,
						Clip: s.LogSnd,
					})
				}
				s.acceptSys.Add(e.BasicEntity, e.Job)
			}
		} else {
			if e.shown {
				e.shown = false
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
	for _, title := range s.titles {
		title.Hidden = true
	}
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

func (s *JobSelectSystem) showAll() {
	s.background.Hidden = false
	s.atkText.Hidden = false
	s.defText.Hidden = false
	s.spdText.Hidden = false
	s.spcText.Hidden = false
	for _, title := range s.titles {
		title.Hidden = false
	}
}

func (s *JobSelectSystem) show(e jobSelectEntity) {
	spcTex, _ := common.LoadedSprite(e.SpecialURL)
	s.spcIcon0.Drawable = spcTex
	s.spcIcon0.Hidden = false
	s.spcIcon1.Drawable = spcTex
	s.spcIcon1.Hidden = false
	s.spcIcon2.Drawable = spcTex
	s.spcIcon2.Hidden = false
	s.spcIcon3.Drawable = spcTex
	s.spcIcon3.Hidden = false
	s.spcText.Drawable = common.Text{
		Text: e.SpecialName,
		Font: s.Fnt,
	}
	s.spcText.Hidden = false
	switch e.Atk {
	case 1:
		s.atkIcon0.Hidden = false
	case 2:
		s.atkIcon0.Hidden = false
		s.atkIcon1.Hidden = false
	case 3:
		s.atkIcon0.Hidden = false
		s.atkIcon1.Hidden = false
		s.atkIcon2.Hidden = false
	case 4:
		s.atkIcon0.Hidden = false
		s.atkIcon1.Hidden = false
		s.atkIcon2.Hidden = false
		s.atkIcon3.Hidden = false
	}
	switch e.Def {
	case 1:
		s.defIcon0.Hidden = false
	case 2:
		s.defIcon0.Hidden = false
		s.defIcon1.Hidden = false
	case 3:
		s.defIcon0.Hidden = false
		s.defIcon1.Hidden = false
		s.defIcon2.Hidden = false
	case 4:
		s.defIcon0.Hidden = false
		s.defIcon1.Hidden = false
		s.defIcon2.Hidden = false
		s.defIcon3.Hidden = false
	}
	switch e.Spd {
	case 1:
		s.spdIcon0.Hidden = false
	case 2:
		s.spdIcon0.Hidden = false
		s.spdIcon1.Hidden = false
	case 3:
		s.spdIcon0.Hidden = false
		s.spdIcon1.Hidden = false
		s.spdIcon2.Hidden = false
	case 4:
		s.spdIcon0.Hidden = false
		s.spdIcon1.Hidden = false
		s.spdIcon2.Hidden = false
		s.spdIcon3.Hidden = false
	}
}

func (s *JobSelectSystem) logDone() bool {
	msg := &CombatLogDoneMessage{}
	engo.Mailbox.Dispatch(msg)
	return msg.Done
}
