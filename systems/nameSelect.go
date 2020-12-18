package systems

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"math/rand"

	"github.com/Noofbiz/SpaceCamp/savedata"
)

type NameSelectPauseMessage struct{}

func (NameSelectPauseMessage) Type() string {
	return "Name Select Pause Message"
}

type NameSelectUnpauseMessage struct{}

func (NameSelectUnpauseMessage) Type() string {
	return "Name Select Unpause Message"
}

type NameSelectSystem struct {
	Fnt    *common.Font
	LogSnd *common.Player
	MaxLen int

	paused, skipNextFrame bool

	box       sprite
	nameplate sprite
	nameText  sprite
	name      string
	alphabet  map[string]*selection

	curSys    *CursorSystem
	acceptSys *AcceptSystem
}

func (s *NameSelectSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CursorSystem:
			s.curSys = sys
		case *AcceptSystem:
			s.acceptSys = sys
		}
	}

	s.box = sprite{BasicEntity: ecs.NewBasic()}
	s.box.Drawable, _ = common.LoadedSprite("start/logs.png")
	s.box.SetZIndex(1)
	s.box.SetCenter(engo.Point{X: 175, Y: 100})
	s.box.Scale = engo.Point{X: 0.6, Y: 2.8}
	s.box.Hidden = true
	w.AddEntity(&s.box)

	s.nameplate = sprite{BasicEntity: ecs.NewBasic()}
	s.nameplate.Drawable, _ = common.LoadedSprite("start/logs.png")
	s.nameplate.SetZIndex(2)
	s.nameplate.SetCenter(engo.Point{X: 200, Y: 75})
	s.nameplate.Scale = engo.Point{X: 0.5, Y: 0.5}
	s.nameplate.Hidden = true
	w.AddEntity(&s.nameplate)

	s.nameText = sprite{BasicEntity: ecs.NewBasic()}
	s.nameText.Drawable = common.Text{
		Font: s.Fnt,
		Text: "This is the name text!!!",
	}
	s.nameText.SetZIndex(3)
	s.nameText.SetCenter(engo.Point{X: 208, Y: 85})
	s.nameText.Scale = engo.Point{X: 0.25, Y: 0.25}
	s.nameText.Hidden = true
	w.AddEntity(&s.nameText)

	s.alphabet = make(map[string]*selection)
	alphabet := []string{
		"A", "B", "C", "D", "E", "F",
		"G", "H", "I", "J", "K", "L",
		"M", "N", "O", "P", "Q", "R",
		"S", "T", "U", "V", "W", "X",
		"Y", "Z", "sp", "dl", "af", "co",
	}
	for i, char := range alphabet {
		row := i / 6
		col := i % 6
		ltr := selection{BasicEntity: ecs.NewBasic()}
		if char == "sp" || char == "dl" || char == "af" || char == "co" {
			ltr.Drawable, _ = common.LoadedSprite("start/" + char + ".png")
		} else {
			ltr.Drawable = common.Text{
				Text: char,
				Font: s.Fnt,
			}
			ltr.Scale = engo.Point{X: 0.25, Y: 0.25}
		}
		ltr.SetCenter(engo.Point{X: 220 + float32(col)*40, Y: 125 + float32(row)*40})
		ltr.SetZIndex(2)
		ltr.Hidden = true
		w.AddEntity(&ltr)
		s.alphabet[char] = &ltr
		s.curSys.Remove(ltr.BasicEntity)
	}

	engo.Mailbox.Listen("Name Select Pause Message", func(m engo.Message) {
		_, ok := m.(NameSelectPauseMessage)
		if !ok {
			return
		}
		s.pause()
	})

	engo.Mailbox.Listen("Name Select Unpause Message", func(m engo.Message) {
		_, ok := m.(NameSelectUnpauseMessage)
		if !ok {
			return
		}
		s.unpause()
	})
}

func (s *NameSelectSystem) Remove(basic ecs.BasicEntity) {}

func (s *NameSelectSystem) Update(dt float32) {
	if s.paused {
		return
	}
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	for char, ent := range s.alphabet {
		if ent.Selected {
			if engo.Input.Button("A").JustPressed() {
				if !s.logDone() {
					continue
				}
				s.addLetter(char)
			}
			if engo.Input.Button("B").JustPressed() {
				if !s.logDone() {
					continue
				}
				s.deleteLetter()
			}
			if engo.Input.Button("X").JustPressed() {
				s.toConfirm()
			}
		}
	}
}

func (s *NameSelectSystem) addLetter(in string) {
	if in == "sp" {
		if len(s.name) >= s.MaxLen {
			return
		}
		s.name = s.name + " "
		s.nameText.Drawable = common.Text{
			Font: s.Fnt,
			Text: s.name,
		}
	} else if in == "dl" {
		s.deleteLetter()
	} else if in == "af" {
		s.setRandomName()
	} else if in == "co" {
		s.confirm()
	} else {
		if len(s.name) >= s.MaxLen {
			return
		}
		s.name = s.name + in
		s.nameText.Drawable = common.Text{
			Font: s.Fnt,
			Text: s.name,
		}
	}
}

func (s *NameSelectSystem) confirm() {
	s.pause()
	msgs := []string{
		"The name's",
		s.name,
		"Eh?",
	}
	for _, msg := range msgs {
		engo.Mailbox.Dispatch(CombatLogMessage{
			Msg:  msg,
			Fnt:  s.Fnt,
			Clip: s.LogSnd,
		})
	}
	s.acceptSys.Add(s.box.GetBasicEntity(), s.name)
}

func (s *NameSelectSystem) setRandomName() {
	rndNames := make([]string, 0)
	rand.Seed(time.Now().UnixNano())
	switch savedata.CurrentSave.Chara1.Job {
	case "Chef":
		rndNames = append(rndNames, "David")
		rndNames = append(rndNames, "Sanji")
		rndNames = append(rndNames, "Soma")
		rndNames = append(rndNames, "Gordon")
	case "Mecha":
		rndNames = append(rndNames, "Jack")
		rndNames = append(rndNames, "Coop")
		rndNames = append(rndNames, "Frye")
		rndNames = append(rndNames, "Jerry")
	case "Medic":
		rndNames = append(rndNames, "Mary")
		rndNames = append(rndNames, "Cuddy")
		rndNames = append(rndNames, "Strax")
		rndNames = append(rndNames, "Krissy")
	case "Pilot":
		rndNames = append(rndNames, "Mark")
		rndNames = append(rndNames, "Launchpad")
		rndNames = append(rndNames, "Mal")
		rndNames = append(rndNames, "Charles")
	case "Defen":
		rndNames = append(rndNames, "Mel")
		rndNames = append(rndNames, "Gun")
		rndNames = append(rndNames, "Sam")
		rndNames = append(rndNames, "Will")
	case "ExoBio":
		rndNames = append(rndNames, "Rich")
		rndNames = append(rndNames, "Albert")
		rndNames = append(rndNames, "Ali")
		rndNames = append(rndNames, "Neil")
	}
	if len(rndNames) <= 0 {
		println("hey now no names!")
		return
	}
	s.name = rndNames[rand.Intn(len(rndNames))]
	s.nameText.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.name,
	}
}

func (s *NameSelectSystem) deleteLetter() {
	if len(s.name) == 0 {
		return
	}
	s.name = s.name[:len(s.name)-1]
	s.nameText.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.name,
	}
}

func (s *NameSelectSystem) toConfirm() {
	engo.Mailbox.Dispatch(CursorSetMessage{ID: s.alphabet["co"]})
}

func (s *NameSelectSystem) pause() {
	s.paused = true
	s.box.Hidden = true
	s.nameplate.Hidden = true
	s.nameText.Hidden = true
	for _, sel := range s.alphabet {
		sel.Hidden = true
		s.curSys.Remove(sel.BasicEntity)
	}
	engo.Mailbox.Dispatch(CursorJumpSetMessage{Jump: 1})
}

func (s *NameSelectSystem) unpause() {
	s.paused = false
	s.box.Hidden = false
	s.nameplate.Hidden = false
	s.nameText.Hidden = false
	s.name = ""
	s.nameText.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.name,
	}
	s.skipNextFrame = true
	alphabet := []string{
		"A", "B", "C", "D", "E", "F",
		"G", "H", "I", "J", "K", "L",
		"M", "N", "O", "P", "Q", "R",
		"S", "T", "U", "V", "W", "X",
		"Y", "Z", "sp", "dl", "af", "co",
	}
	for _, char := range alphabet {
		sel := s.alphabet[char]
		sel.Hidden = false
		s.curSys.Add(&sel.BasicEntity, &sel.SpaceComponent, &sel.RenderComponent, &sel.CursorComponent)
	}
	engo.Mailbox.Dispatch(CursorSetMessage{ID: s.alphabet["A"]})
	engo.Mailbox.Dispatch(CursorJumpSetMessage{Jump: 6})
}

func (s *NameSelectSystem) logDone() bool {
	msg := &CombatLogDoneMessage{}
	engo.Mailbox.Dispatch(msg)
	return msg.Done
}
