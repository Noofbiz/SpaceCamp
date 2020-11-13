package systems

import (
	"image/color"
	"sync"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type CombatLogMessage struct {
	Msg  string
	Fnt  *common.Font
	Clip *common.Player
}

func (m CombatLogMessage) Type() string {
	return "CombatLogMessage"
}

type CombatLogSystem struct {
	lock                                      sync.RWMutex
	log                                       []CombatLogMessage
	idx, charAt                               int
	done, moved                               bool
	elapsed                                   float32
	BackgroundURL, FontURL, DotURL            string
	LineDelay, LetterDelay                    float32
	bg, dot1, dot2, dot3, line1, line2, line3 sprite
}

func (s *CombatLogSystem) New(w *ecs.World) {
	//bg
	s.bg = sprite{BasicEntity: ecs.NewBasic()}
	s.bg.Drawable, _ = common.LoadedSprite(s.BackgroundURL)
	s.bg.SetZIndex(1)
	s.bg.Width = s.bg.Drawable.Width()
	s.bg.Height = s.bg.Drawable.Height()
	s.bg.SetCenter(engo.Point{X: 320, Y: s.bg.Height / 2})
	w.AddEntity(&s.bg)

	dotTex, _ := common.LoadedSprite(s.DotURL)
	//dot1
	s.dot1 = sprite{BasicEntity: ecs.NewBasic()}
	s.dot1.Drawable = dotTex
	s.dot1.SetZIndex(2)
	s.dot1.SetCenter(engo.Point{X: 84, Y: 55})
	s.dot1.Hidden = true
	w.AddEntity(&s.dot1)
	//dot2
	s.dot2 = sprite{BasicEntity: ecs.NewBasic()}
	s.dot2.Drawable = dotTex
	s.dot2.SetZIndex(2)
	s.dot2.SetCenter(engo.Point{X: 84, Y: 35})
	s.dot2.Hidden = true
	w.AddEntity(&s.dot2)
	//dot3
	s.dot3 = sprite{BasicEntity: ecs.NewBasic()}
	s.dot3.Drawable = dotTex
	s.dot3.SetZIndex(2)
	s.dot3.SetCenter(engo.Point{X: 84, Y: 15})
	s.dot3.Hidden = true
	w.AddEntity(&s.dot3)

	logFont := &common.Font{
		Size: 64,
		FG:   color.Black,
		URL:  s.FontURL,
	}
	logFont.CreatePreloaded()
	//line1
	s.line1 = sprite{BasicEntity: ecs.NewBasic()}
	s.line1.Drawable = common.Text{
		Font: logFont,
		Text: "",
	}
	s.line1.Scale = engo.Point{X: 0.2, Y: 0.2}
	s.line1.SetZIndex(2)
	s.line1.Position = engo.Point{X: 99, Y: 54}
	w.AddEntity(&s.line1)
	//line2
	s.line2 = sprite{BasicEntity: ecs.NewBasic()}
	s.line2.Drawable = common.Text{
		Font: logFont,
		Text: "",
	}
	s.line2.Scale = engo.Point{X: 0.2, Y: 0.2}
	s.line2.SetZIndex(2)
	s.line2.Position = engo.Point{X: 99, Y: 34}
	w.AddEntity(&s.line2)
	//line3
	s.line3 = sprite{BasicEntity: ecs.NewBasic()}
	s.line3.Drawable = common.Text{
		Font: logFont,
		Text: "",
	}
	s.line3.Scale = engo.Point{X: 0.2, Y: 0.2}
	s.line3.SetZIndex(2)
	s.line3.Position = engo.Point{X: 99, Y: 14}
	w.AddEntity(&s.line3)

	engo.Mailbox.Listen("CombatLogMessage", func(message engo.Message) {
		msg, ok := message.(CombatLogMessage)
		if !ok {
			return
		}
		s.lock.Lock()
		defer s.lock.Unlock()
		s.log = append(s.log, msg)
	})
}

func (s *CombatLogSystem) Remove(basic ecs.BasicEntity) {}

func (s *CombatLogSystem) Update(dt float32) {
	s.elapsed += dt
	if s.done {
		if s.idx < len(s.log)-1 {
			s.idx++
			s.moved = false
			s.done = false
		}
	} else {
		if !s.moved && len(s.log) > 0 {
			if s.elapsed < s.LineDelay {
				return
			}
			s.elapsed = 0
			s.dot1.Hidden = false
			s.line3.Drawable = s.line2.Drawable
			s.line2.Drawable = s.line1.Drawable
			txt := s.line1.Drawable.(common.Text)
			txt.Font = s.log[s.idx].Fnt
			txt.Text = ""
			if !s.log[s.idx].Clip.IsPlaying() {
				s.log[s.idx].Clip.Rewind()
				s.log[s.idx].Clip.Play()
			}
			s.line1.Drawable = txt
			s.moved = true
			txt2 := s.line2.Drawable.(common.Text)
			txt3 := s.line3.Drawable.(common.Text)
			if txt2.Text != "" {
				s.dot2.Hidden = false
			}
			if txt3.Text != "" {
				s.dot3.Hidden = false
			}
		}
		if len(s.log) > 0 && s.elapsed > s.LetterDelay {
			s.charAt++
			txt := s.line1.Drawable.(common.Text)
			txt.Text = s.log[s.idx].Msg[:s.charAt]
			s.line1.Drawable = txt
			s.elapsed = 0
		}
		if engo.Input.Button("A").JustPressed() {
			txt := s.line1.Drawable.(common.Text)
			txt.Text = s.log[s.idx].Msg
			s.line1.Drawable = txt
			s.charAt = 0
			s.elapsed = 0
			s.done = true
		}
		if len(s.log) > 0 && s.charAt >= len(s.log[s.idx].Msg) {
			s.charAt = 0
			s.elapsed = 0
			s.done = true
		}
	}
}
