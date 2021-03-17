package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Noofbiz/SpaceCamp/savedata"
)

type SelectPhase uint

const (
	JobSelectPhaseOne = iota
	NameSelectPhaseOne
	JobSelectPhaseTwo
	NameSelectPhaseTwo
	JobSelectPhaseThree
	NameSelectPhaseThree
)

var CurrentPhase SelectPhase = JobSelectPhaseOne

type AcceptSystem struct {
	Fnt    *common.Font
	LogSnd *common.Player

	paused, skipNextFrame bool

	box     sprite
	yes, no selection
	job     string

	curSys *CursorSystem

	latestSelection *ecs.BasicEntity
	selections      []*ecs.BasicEntity
}

func (s *AcceptSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CursorSystem:
			s.curSys = sys
		}
	}

	s.box = sprite{BasicEntity: ecs.NewBasic()}
	s.box.Drawable, _ = common.LoadedSprite("start/logs.png")
	s.box.SetZIndex(1)
	s.box.SetCenter(engo.Point{X: 73, Y: 100})
	s.box.Scale = engo.Point{X: 1, Y: 2}
	s.box.Hidden = true
	w.AddEntity(&s.box)

	s.yes = selection{BasicEntity: ecs.NewBasic()}
	s.yes.Drawable = common.Text{
		Text: "yep",
		Font: s.Fnt,
	}
	s.yes.SetZIndex(2)
	s.yes.SetCenter(engo.Point{X: 105, Y: 160})
	s.yes.Scale = engo.Point{X: 0.7, Y: 0.7}
	s.yes.Hidden = true
	w.AddEntity(&s.yes)
	s.curSys.Remove(s.yes.BasicEntity)

	s.no = selection{BasicEntity: ecs.NewBasic()}
	s.no.Drawable = common.Text{
		Text: "nope",
		Font: s.Fnt,
	}
	s.no.SetZIndex(2)
	s.no.SetCenter(engo.Point{X: 400, Y: 160})
	s.no.Scale = engo.Point{X: 0.7, Y: 0.7}
	s.no.Hidden = true
	w.AddEntity(&s.no)
	s.curSys.Remove(s.no.BasicEntity)
}

func (s *AcceptSystem) Add(i *ecs.BasicEntity, job string) {
	s.unpause()
	s.job = job
	s.latestSelection = i
	s.selections = append(s.selections, i)
}

func (s *AcceptSystem) Remove(basic ecs.BasicEntity) {
	s.pause()
	s.job = ""
	delete := -1
	for index, entity := range s.selections {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.selections = append(s.selections[:delete], s.selections[delete+1:]...)
	}
}

func (s *AcceptSystem) Update(dt float32) {
	if s.paused {
		return
	}
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if s.yes.Selected {
		if !s.logDone() {
			return
		}
		if engo.Input.Button("A").JustPressed() {
			s.pause()
			switch CurrentPhase {
			case JobSelectPhaseOne:
				msgs := []string{
					"One of the best!",
					s.job,
					"What was your name, again?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  s.Fnt,
						Clip: s.LogSnd,
					})
				}
				savedata.CurrentSave.Chara1.Job = s.job
				CurrentPhase = NameSelectPhaseOne
				engo.Mailbox.Dispatch(NameSelectUnpauseMessage{})
			case NameSelectPhaseOne:
				msgs := []string{
					"Ah! Now I remember!",
					s.job,
					"Perfect name for a leader!",
					"Now, let's move on.",
					"Your second in command...",
					"What's their job around here?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  s.Fnt,
						Clip: s.LogSnd,
					})
				}
				savedata.CurrentSave.Chara1.Name = s.job
				CurrentPhase = JobSelectPhaseTwo
				engo.Mailbox.Dispatch(JobSelectUnpauseMessage{IDs: s.selections})
			case JobSelectPhaseTwo:
				msgs := []string{
					"Second to none at",
					s.job,
					"What a great ally!",
					"What's their name?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  s.Fnt,
						Clip: s.LogSnd,
					})
				}
				savedata.CurrentSave.Chara2.Job = s.job
				CurrentPhase = NameSelectPhaseTwo
				engo.Mailbox.Dispatch(NameSelectUnpauseMessage{})
			case NameSelectPhaseTwo:
				msgs := []string{
					"So that's their name!",
					s.job,
					"Your second in command is ready!",
					"What about your left hand?",
					"What do they do?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  s.Fnt,
						Clip: s.LogSnd,
					})
				}
				savedata.CurrentSave.Chara2.Name = s.job
				CurrentPhase = JobSelectPhaseThree
				engo.Mailbox.Dispatch(JobSelectUnpauseMessage{IDs: s.selections})
			case JobSelectPhaseThree:
				msgs := []string{
					"Can't go wrong with",
					s.job,
					"Right?",
					"And their name is?",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  s.Fnt,
						Clip: s.LogSnd,
					})
				}
				savedata.CurrentSave.Chara3.Job = s.job
				CurrentPhase = NameSelectPhaseThree
				engo.Mailbox.Dispatch(NameSelectUnpauseMessage{})
			case NameSelectPhaseThree:
				msgs := []string{
					s.job,
					"Is really the perfect name!",
					"I kinda want it!",
					"Now to pack for your journey!",
					"Good luck!",
				}
				for _, msg := range msgs {
					engo.Mailbox.Dispatch(CombatLogMessage{
						Msg:  msg,
						Fnt:  s.Fnt,
						Clip: s.LogSnd,
					})
				}
				savedata.CurrentSave.Chara3.Name = s.job
				engo.Mailbox.Dispatch(LogFinishedSetSceneMessage{
					To:       "Shop Scene",
					NewWorld: true,
				})
			}
		}
		if engo.Input.Button("B").JustPressed() {
			s.Remove(*s.latestSelection)
			s.back()
		}
	}
	if s.no.Selected {
		if !s.logDone() {
			return
		}
		if engo.Input.Button("A").JustPressed() {
			s.Remove(*s.latestSelection)
			s.back()
		}
		if engo.Input.Button("B").JustPressed() {
			s.Remove(*s.latestSelection)
			s.back()
		}
	}
}

func (s *AcceptSystem) pause() {
	s.paused = true
	s.box.Hidden = true
	s.yes.Hidden = true
	s.no.Hidden = true
	s.curSys.Remove(s.yes.BasicEntity)
	s.curSys.Remove(s.no.BasicEntity)
}

func (s *AcceptSystem) unpause() {
	s.paused = false
	s.skipNextFrame = true
	s.box.Hidden = false
	s.yes.Hidden = false
	s.no.Hidden = false
	s.curSys.Add(&s.yes.BasicEntity, &s.yes.SpaceComponent, &s.yes.RenderComponent, &s.yes.CursorComponent)
	s.curSys.Add(&s.no.BasicEntity, &s.no.SpaceComponent, &s.no.RenderComponent, &s.no.CursorComponent)
	engo.Mailbox.Dispatch(CursorSetMessage{ID: s.yes})
}

func (s *AcceptSystem) back() {
	s.pause()
	switch CurrentPhase {
	case JobSelectPhaseOne, JobSelectPhaseTwo, JobSelectPhaseThree:
		s.backToJobSelect()
	case NameSelectPhaseOne, NameSelectPhaseTwo, NameSelectPhaseThree:
		s.backToNameSelect()
	}
}

func (s *AcceptSystem) backToJobSelect() {
	msgs := []string{
		"Well then...",
		"What was that?",
	}
	for _, msg := range msgs {
		engo.Mailbox.Dispatch(CombatLogMessage{
			Msg:  msg,
			Fnt:  s.Fnt,
			Clip: s.LogSnd,
		})
	}
	engo.Mailbox.Dispatch(JobSelectUnpauseMessage{IDs: s.selections})
}

func (s *AcceptSystem) backToNameSelect() {
	msgs := []string{
		"Hmm...",
		"What was that name, again?",
		"I couldn't hear you",
	}
	for _, msg := range msgs {
		engo.Mailbox.Dispatch(CombatLogMessage{
			Msg:  msg,
			Fnt:  s.Fnt,
			Clip: s.LogSnd,
		})
	}
	engo.Mailbox.Dispatch(NameSelectUnpauseMessage{})
}

func (s *AcceptSystem) logDone() bool {
	msg := &CombatLogDoneMessage{}
	engo.Mailbox.Dispatch(msg)
	return msg.Done
}
