package savedata

type SaveData struct {
	Chara1, Chara2, Chara3 *Character
	S                      *Ship
	Scene                  string
}

var CurrentSave = &SaveData{
	Chara1: &Character{},
	Chara2: &Character{},
	Chara3: &Character{},
	S: &Ship{
		Inventory: make(map[string]*Item),
	},
}

func (d *SaveData) LoadSaveData() {}

func (d *SaveData) Save() {}

func (d *SaveData) GetItemQty(name string) int {
	val, ok := d.S.Inventory[name]
	if !ok {
		d.S.Inventory[name] = &Item{Name: name, Quantity: 0}
		return 0
	}
	return val.Quantity
}

func (d *SaveData) SetItemQty(name string, qty int) {
	_, ok := d.S.Inventory[name]
	if !ok {
		d.S.Inventory[name] = &Item{Name: name, Quantity: qty}
	} else {
		d.S.Inventory[name].Name = name
		d.S.Inventory[name].Quantity = qty
	}
}
