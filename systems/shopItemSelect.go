package systems

import (
	"strconv"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

	"github.com/Noofbiz/SpaceCamp/savedata"
)

type ShopItemSelectComponent struct {
	Name, FullName, Desc, Cost string
}

func (c *ShopItemSelectComponent) GetShopItemSelectComponent() *ShopItemSelectComponent {
	return c
}

type ShopItemSelectFace interface {
	GetShopItemSelectComponent() *ShopItemSelectComponent
}

type ShopItemSelectAble interface {
	common.BasicFace
	ShopItemSelectFace
}

type NotShopItemSelectComponent struct{}

func (n *NotShopItemSelectComponent) GetNotShopItemSelectComponent() *NotShopItemSelectComponent {
	return n
}

type NotShopItemSelectAble interface {
	GetNotShopItemSelectComponent() *NotShopItemSelectComponent
}

type shopItemSelectEntity struct {
	*ecs.BasicEntity
	*ShopItemSelectComponent
}

type ShopItemSelectPauseMessage struct{}

func (ShopItemSelectPauseMessage) Type() string { return "Shop Item Select Pause Message" }

type ShopItemSelectUnpauseMessage struct{}

func (ShopItemSelectUnpauseMessage) Type() string { return "Shop Item Select Unpause Message" }

type ShopItemSelectSystem struct {
	entities []shopItemSelectEntity

	bg                         sprite
	URL                        string
	Fnt                        *common.Font
	item1, item2, item3, item4 selection
	name, description          sprite
	price, qty                 sprite

	idx int

	curSys *CursorSystem

	paused, skipNextFrame bool
}

func (s *ShopItemSelectSystem) Priority() int {
	return 1
}

func (s *ShopItemSelectSystem) New(w *ecs.World) {
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *CursorSystem:
			s.curSys = sys
		}
	}

	s.bg = sprite{BasicEntity: ecs.NewBasic()}
	s.bg.Drawable, _ = common.LoadedSprite(s.URL)
	s.bg.SetZIndex(1)
	s.bg.Position = engo.Point{X: 70, Y: 100}
	w.AddEntity(&s.bg)

	s.item1 = selection{BasicEntity: ecs.NewBasic()}
	s.item1.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Item 1",
	}
	s.item1.SetZIndex(2)
	s.item1.Scale = engo.Point{X: 0.3, Y: 0.3}
	s.item1.Position = engo.Point{X: 100, Y: 115}
	s.item1.Selected = true
	w.AddEntity(&s.item1)

	s.item2 = selection{BasicEntity: ecs.NewBasic()}
	s.item2.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Item 2",
	}
	s.item2.SetZIndex(2)
	s.item2.Scale = engo.Point{X: 0.3, Y: 0.3}
	s.item2.Position = engo.Point{X: 100, Y: 157.5}
	w.AddEntity(&s.item2)

	s.item3 = selection{BasicEntity: ecs.NewBasic()}
	s.item3.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Item 3",
	}
	s.item3.SetZIndex(2)
	s.item3.Scale = engo.Point{X: 0.3, Y: 0.3}
	s.item3.Position = engo.Point{X: 100, Y: 200}
	w.AddEntity(&s.item3)

	s.item4 = selection{BasicEntity: ecs.NewBasic()}
	s.item4.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Item 4",
	}
	s.item4.SetZIndex(2)
	s.item4.Scale = engo.Point{X: 0.3, Y: 0.3}
	s.item4.Position = engo.Point{X: 100, Y: 237.5}
	w.AddEntity(&s.item4)

	s.name = sprite{BasicEntity: ecs.NewBasic()}
	s.name.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Name",
	}
	s.name.SetZIndex(2)
	s.name.Scale = engo.Point{X: 0.2, Y: 0.2}
	s.name.Position = engo.Point{X: 255, Y: 125}
	w.AddEntity(&s.name)

	s.description = sprite{BasicEntity: ecs.NewBasic()}
	s.description.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Desc Line 1 \nDesc Line 2 \nDesc Line 3",
	}
	s.description.SetZIndex(2)
	s.description.Scale = engo.Point{X: 0.2, Y: 0.2}
	s.description.Position = engo.Point{X: 255, Y: 150}
	w.AddEntity(&s.description)

	s.price = sprite{BasicEntity: ecs.NewBasic()}
	s.price.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Cost: XXX",
	}
	s.price.SetZIndex(2)
	s.price.Scale = engo.Point{X: 0.2, Y: 0.2}
	s.price.Position = engo.Point{X: 255, Y: 215}
	w.AddEntity(&s.price)

	s.qty = sprite{BasicEntity: ecs.NewBasic()}
	s.qty.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Owned: XX",
	}
	s.qty.SetZIndex(2)
	s.qty.Scale = engo.Point{X: 0.2, Y: 0.2}
	s.qty.Position = engo.Point{X: 255, Y: 240}
	w.AddEntity(&s.qty)

	engo.Mailbox.Listen("Shop Item Select Pause Message", func(m engo.Message) {
		_, ok := m.(ShopItemSelectPauseMessage)
		if !ok {
			return
		}
		s.pause()
	})
	engo.Mailbox.Listen("Shop Item Select Unpause Message", func(m engo.Message) {
		_, ok := m.(ShopItemSelectUnpauseMessage)
		if !ok {
			return
		}
		s.unpause()
	})
}

func (s *ShopItemSelectSystem) Add(basic *ecs.BasicEntity, item *ShopItemSelectComponent) {
	if len(s.entities) == 0 {
		s.item1.Drawable = common.Text{
			Font: s.Fnt,
			Text: item.Name,
		}
		s.name.Drawable = common.Text{
			Font: s.Fnt,
			Text: item.FullName,
		}
		s.description.Drawable = common.Text{
			Font: s.Fnt,
			Text: item.Desc,
		}
		s.price.Drawable = common.Text{
			Font: s.Fnt,
			Text: "Cost: " + item.Cost,
		}
		s.qty.Drawable = common.Text{
			Font: s.Fnt,
			Text: "Owned: " + strconv.Itoa(savedata.CurrentSave.GetItemQty(item.Name)),
		}
	} else if len(s.entities) == 1 {
		s.item2.Drawable = common.Text{
			Font: s.Fnt,
			Text: item.Name,
		}
	} else if len(s.entities) == 2 {
		s.item3.Drawable = common.Text{
			Font: s.Fnt,
			Text: item.Name,
		}
	} else if len(s.entities) == 3 {
		s.item4.Drawable = common.Text{
			Font: s.Fnt,
			Text: item.Name,
		}
	}
	s.entities = append(s.entities, shopItemSelectEntity{basic, item})
}

func (s *ShopItemSelectSystem) AddByInterface(i ecs.Identifier) {
	o, ok := i.(ShopItemSelectAble)
	if !ok {
		return
	}
	s.Add(o.GetBasicEntity(), o.GetShopItemSelectComponent())
}

func (s *ShopItemSelectSystem) Remove(basic ecs.BasicEntity) {
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

func (s *ShopItemSelectSystem) Update(dt float32) {
	if s.paused {
		return
	}
	if s.skipNextFrame {
		s.skipNextFrame = false
		return
	}
	if engo.Input.Button("X").JustPressed() {
		engo.Mailbox.Dispatch(CursorSetMessage{ID: s.item4})
		s.idx = len(s.entities) - 4
		s.item1.Drawable = common.Text{
			Font: s.Fnt,
			Text: s.entities[s.idx].Name,
		}
		s.item2.Drawable = common.Text{
			Font: s.Fnt,
			Text: s.entities[s.idx+1].Name,
		}
		s.item3.Drawable = common.Text{
			Font: s.Fnt,
			Text: s.entities[s.idx+2].Name,
		}
		s.item4.Drawable = common.Text{
			Font: s.Fnt,
			Text: s.entities[s.idx+3].Name,
		}
	}
	index := -1
	if s.item1.Selected {
		if engo.Input.Button("up").JustPressed() {
			s.idx--
			if s.idx < 0 {
				s.idx = 0
			}
			s.item1.Drawable = common.Text{
				Font: s.Fnt,
				Text: s.entities[s.idx].Name,
			}
			s.item2.Drawable = common.Text{
				Font: s.Fnt,
				Text: s.entities[s.idx+1].Name,
			}
			s.item3.Drawable = common.Text{
				Font: s.Fnt,
				Text: s.entities[s.idx+2].Name,
			}
			s.item4.Drawable = common.Text{
				Font: s.Fnt,
				Text: s.entities[s.idx+3].Name,
			}
		}
		index = s.idx
	} else if s.item2.Selected {
		index = s.idx + 1
	} else if s.item3.Selected {
		index = s.idx + 2
	} else if s.item4.Selected {
		if engo.Input.Button("down").JustPressed() {
			s.idx++
			if s.idx+3 >= len(s.entities) {
				s.idx = len(s.entities) - 4
			}
			s.item1.Drawable = common.Text{
				Font: s.Fnt,
				Text: s.entities[s.idx].Name,
			}
			s.item2.Drawable = common.Text{
				Font: s.Fnt,
				Text: s.entities[s.idx+1].Name,
			}
			s.item3.Drawable = common.Text{
				Font: s.Fnt,
				Text: s.entities[s.idx+2].Name,
			}
			s.item4.Drawable = common.Text{
				Font: s.Fnt,
				Text: s.entities[s.idx+3].Name,
			}
		}
		index = s.idx + 3
	}
	if index < 0 || index >= len(s.entities) {
		return
	}
	s.name.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.entities[index].FullName,
	}
	s.description.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.entities[index].Desc,
	}
	s.price.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Cost: " + s.entities[index].Cost,
	}
	s.qty.Drawable = common.Text{
		Font: s.Fnt,
		Text: "Owned: " + strconv.Itoa(savedata.CurrentSave.GetItemQty(s.entities[index].Name)),
	}
	if engo.Input.Button("A").JustPressed() && s.logDone() {
		s.pause()
		engo.Mailbox.Dispatch(ShopPurchaseUnpauseMessage{ID: s.entities[index]})
	}
}

func (s *ShopItemSelectSystem) unpause() {
	s.paused = false
	s.skipNextFrame = true
	s.bg.Hidden = false
	s.description.Hidden = false
	s.price.Hidden = false
	s.qty.Hidden = false
	s.name.Hidden = false
	s.item1.Hidden = false
	s.curSys.AddByInterface(&s.item1)
	s.item1.Selected = true
	s.item2.Hidden = false
	s.curSys.AddByInterface(&s.item2)
	s.item3.Hidden = false
	s.curSys.AddByInterface(&s.item3)
	s.item4.Hidden = false
	s.curSys.AddByInterface(&s.item4)
	engo.Mailbox.Dispatch(CursorSetMessage{s.item1})
	s.idx = 0
	s.item1.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.entities[s.idx].Name,
	}
	s.item2.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.entities[s.idx+1].Name,
	}
	s.item3.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.entities[s.idx+2].Name,
	}
	s.item4.Drawable = common.Text{
		Font: s.Fnt,
		Text: s.entities[s.idx+3].Name,
	}
}

func (s *ShopItemSelectSystem) pause() {
	s.paused = true
	s.bg.Hidden = true
	s.description.Hidden = true
	s.price.Hidden = true
	s.qty.Hidden = true
	s.name.Hidden = true
	s.item1.Hidden = true
	s.curSys.Remove(s.item1.BasicEntity)
	engo.Mailbox.Dispatch(CursorSetMessage{ID: s.item1})
	s.item2.Hidden = true
	s.curSys.Remove(s.item2.BasicEntity)
	s.item3.Hidden = true
	s.curSys.Remove(s.item3.BasicEntity)
	s.item4.Hidden = true
	s.curSys.Remove(s.item4.BasicEntity)
}

func (s *ShopItemSelectSystem) logDone() bool {
	msg := &CombatLogDoneMessage{}
	engo.Mailbox.Dispatch(msg)
	return msg.Done
}
