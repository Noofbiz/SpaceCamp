package systems

import (
	"image/color"
	"strconv"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Noofbiz/SpaceCamp/savedata"
)

type ShopPurchasePauseMessage struct{}

func (ShopPurchasePauseMessage) Type() string { return "Shop Purchase Pause Message" }

type ShopPurchaseUnpauseMessage struct {
	ID ecs.Identifier
}

func (ShopPurchaseUnpauseMessage) Type() string { return "Shop Purchase Unpause Message" }

type ShopPurchaseComponent struct {
	MaxQuantity int
	Price       int
}

func (c *ShopPurchaseComponent) GetShopPurchaseComponent() *ShopPurchaseComponent {
	return c
}

type ShopPurchaseFace interface {
	GetShopPurchaseComponent() *ShopPurchaseComponent
}

type ShopPurchaseAble interface {
	common.BasicFace
	ShopItemSelectFace
	ShopPurchaseFace
}

type NotShopPurchaseComponent struct{}

func (c *NotShopPurchaseComponent) GetNotShopPurchaseComponent() *NotShopPurchaseComponent {
	return c
}

type NotShopPurchaseAble interface {
	GetNotShopPurchaseComponent() *NotShopPurchaseComponent
}

type shopPurchaseEntity struct {
	*ecs.BasicEntity
	*ShopItemSelectComponent
	*ShopPurchaseComponent
}

type ShopPurchaseSystem struct {
	entities []shopPurchaseEntity
	entidx   int

	Fnt *common.Font
	Snd *common.Player

	bg    sprite
	BGURL string

	moneyBG    sprite
	MoneyBGURL string
	moneyAmt   sprite

	number          sprite
	currentNumber   int
	displayedNumber int

	curSys *CursorSystem

	paused, skipNextFrame bool
	switchWhenLogDone     bool
}

func (s *ShopPurchaseSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CursorSystem:
			s.curSys = sys
		}
	}

	s.currentNumber = 1
	s.paused = true

	s.bg = sprite{BasicEntity: ecs.NewBasic()}
	s.bg.Drawable = common.Rectangle{}
	s.bg.Color = color.White
	s.bg.Height = 125
	s.bg.Width = 200
	s.bg.SetZIndex(1)
	s.bg.Position = engo.Point{X: 120, Y: 100}
	s.bg.Hidden = true
	w.AddEntity(&s.bg)

	s.moneyBG = sprite{BasicEntity: ecs.NewBasic()}
	s.moneyBG.Drawable = common.Rectangle{}
	s.moneyBG.Color = color.RGBA{0x00, 0xff, 0x00, 0xff}
	s.moneyBG.Height = 50
	s.moneyBG.Width = 125
	s.moneyBG.SetZIndex(3)
	s.moneyBG.Position = engo.Point{X: 340, Y: 75}
	w.AddEntity(&s.moneyBG)

	s.moneyAmt = sprite{BasicEntity: ecs.NewBasic()}
	s.moneyAmt.Drawable = common.Text{
		Font: s.Fnt,
		Text: "$" + strconv.Itoa(savedata.CurrentSave.S.Money),
	}
	s.moneyAmt.SetZIndex(4)
	s.moneyAmt.Scale = engo.Point{X: 0.35, Y: 0.35}
	s.moneyAmt.Position = engo.Point{X: 345, Y: 95}
	w.AddEntity(&s.moneyAmt)

	s.number = sprite{BasicEntity: ecs.NewBasic()}
	s.number.Drawable = common.Text{
		Font: s.Fnt,
		Text: strconv.Itoa(s.currentNumber),
	}
	s.number.SetZIndex(2)
	s.number.Scale = engo.Point{X: 0.75, Y: 0.75}
	s.number.Position = engo.Point{X: 170, Y: 138}
	s.number.Hidden = true
	w.AddEntity(&s.number)

	engo.Mailbox.Listen("Shop Purchase Pause Message", func(m engo.Message) {
		_, ok := m.(ShopItemSelectPauseMessage)
		if !ok {
			return
		}
		s.pause()
	})
	engo.Mailbox.Listen("Shop Purchase Unpause Message", func(m engo.Message) {
		msg, ok := m.(ShopPurchaseUnpauseMessage)
		if !ok {
			return
		}
		idx := -1
		for i, ent := range s.entities {
			if msg.ID.ID() == ent.ID() {
				idx = i
			}
		}
		if idx < 0 {
			engo.Mailbox.Dispatch(ShopItemSelectUnpauseMessage{})
			return
		}
		if s.entities[idx].Name == "Leave" {
			//set new scene!
			msgs := []string{
				"All done?",
				"Good luck on your journey!",
			}

			for _, msg := range msgs {
				engo.Mailbox.Dispatch(CombatLogMessage{
					Msg:  msg,
					Fnt:  s.Fnt,
					Clip: s.Snd,
				})
			}
			engo.Mailbox.Dispatch(LogFinishedSetSceneMessage{
				To:       "Take Off Scene",
				NewWorld: true,
			})
			return
		}

		msgs := []string{
			"Ah, interested in",
			s.entities[idx].FullName + "?",
		}

		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  s.Fnt,
				Clip: s.Snd,
			})
		}

		if savedata.CurrentSave.GetItemQty(s.entities[idx].Name) >= s.entities[idx].MaxQuantity {
			msgs := []string{
				"It appears you already have",
				"more than enough",
				s.entities[idx].FullName,
			}

			for _, msg := range msgs {
				engo.Mailbox.Dispatch(CombatLogMessage{
					Msg:  msg,
					Fnt:  s.Fnt,
					Clip: s.Snd,
				})
			}
			engo.Mailbox.Dispatch(ShopItemSelectUnpauseMessage{})
			return
		}

		if savedata.CurrentSave.S.Money < s.entities[idx].Price {
			msgs := []string{
				"You don't have enough for",
				s.entities[idx].FullName,
				"Sorry about that!",
			}

			for _, msg := range msgs {
				engo.Mailbox.Dispatch(CombatLogMessage{
					Msg:  msg,
					Fnt:  s.Fnt,
					Clip: s.Snd,
				})
			}
			engo.Mailbox.Dispatch(ShopItemSelectUnpauseMessage{})
			return
		}
		engo.Mailbox.Dispatch(CombatLogMessage{
			Msg:  "How many would you like?",
			Fnt:  s.Fnt,
			Clip: s.Snd,
		})
		s.unpause(msg.ID.ID())
	})
}

func (s *ShopPurchaseSystem) Add(basic *ecs.BasicEntity, sel *ShopItemSelectComponent, pur *ShopPurchaseComponent) {
	s.entities = append(s.entities, shopPurchaseEntity{basic, sel, pur})
}

func (s *ShopPurchaseSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(ShopPurchaseAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetShopItemSelectComponent(), o.GetShopPurchaseComponent())
}

func (s *ShopPurchaseSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *ShopPurchaseSystem) Update(dt float32) {
	if s.paused {
		return
	}
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if engo.Input.Button("up").JustPressed() || engo.Input.Button("right").JustPressed() {
		s.currentNumber++
	}
	if engo.Input.Button("down").JustPressed() || engo.Input.Button("left").JustPressed() {
		s.currentNumber--
	}
	if s.currentNumber <= 0 {
		s.currentNumber = 1
	}
	if max := s.entities[s.entidx].MaxQuantity - savedata.CurrentSave.GetItemQty(s.entities[s.entidx].Name); s.currentNumber > max {
		s.currentNumber = max
	}
	if max := savedata.CurrentSave.S.Money / s.entities[s.entidx].Price; s.currentNumber > max {
		s.currentNumber = max
	}
	if s.currentNumber != s.displayedNumber {
		s.number.Drawable = common.Text{
			Font: s.Fnt,
			Text: strconv.Itoa(s.currentNumber),
		}
		amount := savedata.CurrentSave.S.Money
		amount -= s.entities[s.entidx].Price * s.currentNumber
		s.moneyAmt.Drawable = common.Text{
			Font: s.Fnt,
			Text: "$" + strconv.Itoa(amount),
		}
		s.displayedNumber = s.currentNumber
	}
	if engo.Input.Button("A").JustPressed() && s.logDone() {
		savedata.CurrentSave.SetItemQty(s.entities[s.entidx].Name, savedata.CurrentSave.GetItemQty(s.entities[s.entidx].Name)+s.currentNumber)
		savedata.CurrentSave.S.Money -= s.entities[s.entidx].Price * s.currentNumber

		msgs := []string{
			"Let met get that for you!",
			strconv.Itoa(s.currentNumber) + " " + s.entities[s.entidx].FullName + "s",
			"Have been loaded on your ship",
			"Can I get you anything else?",
		}

		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  s.Fnt,
				Clip: s.Snd,
			})
		}
		s.pause()
		engo.Mailbox.Dispatch(ShopItemSelectUnpauseMessage{})
	}
	if engo.Input.Button("B").JustPressed() && s.logDone() {
		msgs := []string{
			"Not interested after all?",
			"That's cool!",
			"Would you like anything else?",
		}

		for _, msg := range msgs {
			engo.Mailbox.Dispatch(CombatLogMessage{
				Msg:  msg,
				Fnt:  s.Fnt,
				Clip: s.Snd,
			})
		}
		s.pause()
		engo.Mailbox.Dispatch(ShopItemSelectUnpauseMessage{})
	}
}

func (s *ShopPurchaseSystem) pause() {
	s.paused = true
	s.bg.Hidden = true
	s.number.Hidden = true
	s.moneyAmt.Drawable = common.Text{
		Font: s.Fnt,
		Text: "$" + strconv.Itoa(savedata.CurrentSave.S.Money),
	}
}

func (s *ShopPurchaseSystem) unpause(id uint64) {
	s.paused = false
	s.skipNextFrame = true
	s.bg.Hidden = false
	s.number.Hidden = false
	s.currentNumber = 1
	s.displayedNumber = -1
	s.entidx = 0
	for i, ent := range s.entities {
		if ent.ID() == id {
			s.entidx = i
			break
		}
	}
}

func (s *ShopPurchaseSystem) logDone() bool {
	msg := &CombatLogDoneMessage{}
	engo.Mailbox.Dispatch(msg)
	return msg.Done
}
