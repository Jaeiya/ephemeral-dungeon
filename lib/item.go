package lib

type ItemFlags int

const (
	Equippable ItemFlags = 1 << iota
	Breakable
	Sellable
)

type Item struct {
	Name  string
	Desc  string
	Flags ItemFlags
}

func NewItem(name, desc string, flags ItemFlags) Item {
	return Item{name, desc, flags}
}

func (i Item) Is(f ItemFlags) bool {
	return i.Flags&f == f
}

func (i Item) IsAny(f ItemFlags) bool {
	return f&i.Flags != 0
}
