package savedata

type Character struct {
	Job, Name string
}

type Ship struct {
	Money     int
	Inventory map[string]*Item
}

type Item struct {
	Name     string
	Quantity int
}
