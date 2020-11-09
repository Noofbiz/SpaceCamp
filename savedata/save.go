package savedata

type SaveData struct{}

var CurrentSave = &SaveData{}

func (d *SaveData) LoadSaveData() {}
