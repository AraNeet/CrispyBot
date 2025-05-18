package variables

type StatType int
type TraitType int

const (
	Vitality StatType = iota
	Intelligence
	Strength
	Durability
	Speed
	SkillLevel
	ManaFlow
	Height
	Alignment
)

const (
	Buff TraitType = iota
	Weakness
	X_Factor
	Race
)
