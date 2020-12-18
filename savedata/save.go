package savedata

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

var location = filepath.Join("outtahere", "spacecamp.sav")

func (data *SaveData) Save() error {
	os.Remove(location)
	f, err := os.Create(location)
	if err != nil {
		return err
	}
	defer f.Close()
	var d []byte
	d, err = json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err = f.Write(d); err != nil {
		return err
	}
	return nil
}

func (dat *SaveData) Load() error {
	// for creating executables
	// exep, err := os.Executable()
	// if err != nil {
	//	return err
	//}
	//location = filepath.Join(filepath.Dir(exep), "data", "opts.json")
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	location = filepath.Join(pwd, "outtahere", "spacecamp.sav")
	f, err := os.Open(location)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, CurrentSave)
	if err != nil {
		CurrentSave = &SaveData{
			Chara1: &Character{},
			Chara2: &Character{},
			Chara3: &Character{},
			S: &Ship{
				Inventory: make(map[string]*Item),
			},
		}
		return err
	}
	return nil
}

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
