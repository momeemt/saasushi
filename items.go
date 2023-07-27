package main

type Item struct {
	ImageSource string   `json:"image_source"`
	Kind        ItemKind `json:"kind"`
	Name        string   `json:"name"`
	Price       uint64   `json:"price"`
	Calorie     uint64   `json:"calorie"`
	Note        []string `json:"note"`
}

type ItemKind int

const (
	LimitedEdition ItemKind = iota
	Nigiri
	GunkanMakimono
	SideMenu
	Drink
	Dessert
)

func (ik ItemKind) MarshalJSON() ([]byte, error) {
	switch ik {
	case LimitedEdition:
		return []byte(`"期間限定"`), nil
	case Nigiri:
		return []byte(`"にぎり"`), nil
	case GunkanMakimono:
		return []byte(`"軍艦・巻物"`), nil
	case SideMenu:
		return []byte(`"サイドメニュー"`), nil
	case Drink:
		return []byte(`"ドリンク"`), nil
	case Dessert:
		return []byte(`"デザート"`), nil
	default:
		return nil, nil
	}
}
