package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type ShopKeeperSystem struct {
	keeper   animation
	keeperSS *common.Spritesheet

	KeeperURL string

	talking bool
}

func (s *ShopKeeperSystem) New(w *ecs.World) {
	s.keeperSS = common.NewSpritesheetWithBorderFromFile(s.KeeperURL, 225, 150, 1, 1)
	s.keeper = animation{BasicEntity: ecs.NewBasic()}
	s.keeper.Drawable = s.keeperSS.Drawable(0)
	s.keeper.SetZIndex(1)
	s.keeper.Position = engo.Point{X: 375, Y: 100}
	s.keeper.AnimationComponent = common.NewAnimationComponent(s.keeperSS.Drawables(), 0.1)
	s.keeper.AddAnimation(&common.Animation{Name: "talk", Frames: []int{7, 8, 7}})
	s.keeper.AddAnimation(&common.Animation{Name: "blink", Frames: []int{0, 1, 2, 3, 4, 5, 6}})
	w.AddEntity(&s.keeper)
}

func (s *ShopKeeperSystem) Remove(basic ecs.BasicEntity) {}

func (s *ShopKeeperSystem) Update(dt float32) {
	if !s.logDone() {
		s.talking = true
		if s.keeper.CurrentAnimation == nil {
			s.keeper.SelectAnimationByName("talk")
		}
	} else {
		if s.talking {
			s.keeper.SelectAnimationByName("blink")
			s.talking = false
		}
	}
}

func (s *ShopKeeperSystem) logDone() bool {
	msg := &CombatLogDoneMessage{}
	engo.Mailbox.Dispatch(msg)
	return msg.Done
}
