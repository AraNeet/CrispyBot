package variables

type StatType int
type TraitType int
type CharacteristicType int

const (
	Vitality StatType = iota
	Durability
	Strength
	Speed
	Intelligence
	Mastery
	Mana
)

const (
	Innate TraitType = iota
	Inadequacy
	X_Factor
)

const (
	Alignment CharacteristicType = iota
	Race
	Height
	Element
)
