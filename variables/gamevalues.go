package variables

const (
	Common_Chance    = 50
	Uncommon_Chance  = 25
	Rare_Chance      = 15
	Epic_Chance      = 8
	Legendary_Chance = 2
)

// Combat system constants
const (
	// Stat conversion rates (10:1 ratio as requested)
	VitalityToHPRatio         = 10 // 1 Vitality = 10 HP
	DurabilityToDefenseRatio  = 10 // 1 Durability = 10 Defense
	SpeedToInitiativeRatio    = 10 // 1 Speed = 10 Initiative
	StrengthToDamageRatio     = 10 // 1 Strength = 10 base physical damage
	IntelligenceToDamageRatio = 10 // 1 Intelligence = 10 base magic damage
	ManaToPoolRatio           = 10 // 1 Mana = 10 MP
	MasteryToAccuracyRatio    = 10 // 1 Mastery = 10 accuracy points

	// Combat limits
	MaxStatValue   = 300 // Maximum value for any base stat
	MaxHP          = MaxStatValue * VitalityToHPRatio
	MaxMP          = MaxStatValue * ManaToPoolRatio
	MaxDodgeChance = 30 // Maximum dodge chance (30%)

	// Base values
	BaseAccuracy         = 70  // Base accuracy percentage
	BaseDodgeChance      = 5   // Base dodge chance percentage
	BaseCritChance       = 5   // Base critical hit chance percentage
	CritDamageMultiplier = 1.5 // Damage multiplier on critical hit

	// Action costs
	PhysicalAttackManaCost  = 0  // Mana cost for physical attacks
	MagicAttackBaseManaCost = 15 // Base mana cost for magic attacks
)

var (
	VitalityValue = map[string]float64{
		"Helpless-": 0,
		"Helpless":  45.75,
		"Frail":     85.75,
		"Weak":      110,
		"Average":   127.5,
		"Heathly":   145.5,
		"Robust":    166,
		"Vigorous":  180,
		"Vigorous+": 200,
	}

	IntelligenceValue = map[string]float64{
		"Mindless-":   0,
		"Mindless":    45.75,
		"Lobotomized": 85.75,
		"Dumb":        110,
		"Average":     127.5,
		"Smart":       145.5,
		"Genius":      166,
		"Prodigious":  180,
		"Prodigious+": 200,
	}

	StrengthValue = map[string]float64{
		"Forceless-":    0,
		"Forceless":     45.75,
		"Scrawny":       85.75,
		"Weak":          110,
		"Average":       127.5,
		"Strong":        145.5,
		"Formidable":    166,
		"Overpowering":  180,
		"Overpowering+": 200,
	}

	SpeedValue = map[string]float64{
		"Torpid":      0,
		"Crippled":    45.75,
		"Sluggish":    85.75,
		"Slow":        110,
		"Average":     127.5,
		"Fast":        145.5,
		"Accelerated": 166,
		"Supersonic":  180,
		"Hypersonic":  200,
	}

	DurabilityValue = map[string]float64{
		"Defenseless-": 0,
		"Defenseless":  45.75,
		"Vulnerable":   85.75,
		"Vincible":     110,
		"Average":      127.5,
		"Reinforced":   145.5,
		"Armored":      166,
		"Fortified":    180,
		"Fortified+":   200,
	}

	ManaValue = map[string]float64{
		"No-Mana":      0,
		"Mana-less":    45.75,
		"Lowly":        85.75,
		"Hexed":        110,
		"Average":      127.5,
		"Enchanted":    145.5,
		"Conjuring":    166,
		"Overflowing":  180,
		"Overflowing+": 200,
	}

	MasteryValue = map[string]float64{
		"Skill-Less-": 0,
		"Skill-Less":  45.75,
		"Novice":      85.75,
		"Amateur":     110,
		"Average":     127.5,
		"Skilled":     145.5,
		"Expert":      166,
		"Mastered":    180,
		"Mastered+":   200,
	}

	InnateValues = map[string][]map[string]float64{
		"Blessed":           {{"Vitality": 50, "Strength": 50, "Intelligence": 50, "Durability": 50, "Speed": 50, "Mastery": 50, "Mana": 50}},
		"Speed Force":       {{"Speed": 75}},
		"Fast Learner":      {{"Intelligence": 75}},
		"Abounding Flow":    {{"Mana": 75}},
		"Rough Skin":        {{"Durability": 50}},
		"Castle Training":   {{"Strength": 50}},
		"Call of Hercules":  {{"Strength": 75}},
		"Isekai Protag":     {{"Vitality": 25, "Strength": 25, "Intelligence": 25, "Durability": 25, "Speed": 25, "Mastery": 25, "Mana": 25}},
		"Druid's Blessing":  {{"Vitality": 75}},
		"Naturally Skilled": {{"Mastery": 75}},
		"Swift":             {{"Speed": 25}},
		"Big Boned":         {{"Vitality": 25, "Strength": 25, "Durability": 25}},
		"Quick Thinker":     {{"Intelligence": 25}},
		"None":              {{}},
	}

	InadequacyValues = map[string][]map[string]float64{
		"Fragile Bone":        {{"Vitality": 25, "Durability": 25}},
		"STD":                 {{"Strength": 25}},
		"Cancer":              {{"Strength": 25}},
		"Delayed Reaction":    {{"Speed": 25}},
		"Testicular Torsion":  {{"Speed": 50}},
		"Amputee":             {{"Mastery": 25}},
		"Blindness":           {{"Mastery": 50, "Speed": 50}},
		"Too Young":           {{"Mastery": 25, "Strength": 25}},
		"Too Old":             {{"Speed": 25, "Strength": 25, "Durability": 25}},
		"One Eye":             {{"Speed": 25, "Mastery": 25}},
		"Lobotomized":         {{"Vitality": 50, "Strength": 50, "Intelligence": 50, "Durability": 50, "Speed": 50, "Mastery": 50, "Mana": 50}},
		"Auto Immune Disease": {{"Vitality": 25, "Strength": 25, "Intelligence": 25, "Durability": 25, "Speed": 25, "Mastery": 25, "Mana": 25}},
		"Depression":          {{"Vitality": 25, "Strength": 25, "Speed": 25}},
		"Claustrophobia":      {{"Speed": 50}},
		"Paranoid":            {{"Mastery": 25, "Vitality": 50}},
		"Cursed":              {{"Vitality": 75, "Strength": 75, "Intelligence": 75, "Durability": 75, "Speed": 75, "Mastery": 75, "Mana": 75}},
		"None":                {{}},
	}

	RaceValues = map[string][]map[string][]map[string]float64{
		"Elf":          {{"Buff": {{"Intelligence": 50, "Speed": 50}}, "Weakness": {{"Durability": 50}}}},
		"Dwarf":        {{"Buff": {{"Durability": 75, "Mastery": 75}}, "Weakness": {{"Speed": 50}}}},
		"Demon":        {{"Buff": {{"Strength": 75}}, "Weakness": {{"Durability": 50}}}},
		"Angel":        {{"Buff": {{"Vitality": 50, "Strength": 50, "Intelligence": 50, "Durability": 50, "Speed": 50, "Mastery": 50, "Mana": 50}}}},
		"Dullahan":     {{"Buff": {{"Strength": 75, "Speed": 50}}, "Weakness": {{"Durability": 50}}}},
		"Lamia":        {{"Buff": {{"Durability": 50, "Vitality": 50}}, "Weakness": {{"Speed": 25}}}},
		"Gnome":        {{"Buff": {{"Intelligence": 75}}, "Weakness": {{"Height": 75}}}},
		"Orc":          {{"Buff": {{"Strength": 75, "Vitality": 75}}, "Weakness": {{"Intelligence": 75}}}},
		"Goblin":       {{"Buff": {{"Intelligence": 50, "Strength": 25}}, "Weakness": {{"Height": 75}}}},
		"Giant":        {{"Buff": {{"Height": 75, "Strength": 50}}, "Weakness": {{"Intelligence": 50}}}},
		"Centaur":      {{"Buff": {{"Vitality": 50, "Speed": 50, "Height": 25}}}},
		"Kobold":       {{"Buff": {{"Speed": 75}}, "Weakness": {{"Height": 50, "Durability": 25}}}},
		"Beastfolk":    {{"Buff": {{"Strength": 25, "Vitality": 25, "Intelligence": 25, "Durability": 25, "Speed": 25, "Mastery": 25, "Mana": 25, "Height": 25}}}},
		"Mushfolk":     {{"Buff": {{"Mana": 75}}, "Weakness": {{"Vitality": 25, "Durability": 50, "Height": 50}}}},
		"Merfolk":      {{"Buff": {{"Mana": 50, "Strength": 25, "Durability": 25}}}},
		"Dragonborn":   {{"Buff": {{"Strength": 75, "Height": 50, "Durability": 50, "Vitality": 50}}}},
		"Fairy":        {{"Buff": {{"Speed": 75, "Mana": 75}}, "Weakness": {{"Height": 75}}}},
		"Cyclops":      {{"Buff": {{"Height": 75, "Strength": 50}}, "Weakness": {{"Intelligence": 50}}}},
		"Druid":        {{"Buff": {{"Mana": 75, "Vitality": 50}}, "Weakness": {{"Durability": 75}}}},
		"God":          {{"Buff": {{"Strength": 75, "Vitality": 75, "Intelligence": 75, "Durability": 75, "Speed": 75, "Mastery": 75, "Mana": 75, "Height": 75}}}},
		"Minotaur":     {{"Buff": {{"Vitality": 50, "Speed": 50, "Height": 25}}}},
		"Wizard/Witch": {{"Buff": {{"Mana": 75, "Mastery": 50, "Intelligence": 50}}, "Weakness": {{"Durability": 75}}}},
		"Vampire":      {{"Buff": {{"Strength": 50, "Vitality": 50, "Speed": 50, "Mana": 75}}}},
		"Werewolf":     {{"Buff": {{"Strength": 75, "Vitality": 75, "Speed": 75}}}},
		"Undead":       {{"Buff": {{"Vitality": 75, "Strength": 50}}, "Weakness": {{"Intelligence": 25}}}},
		"Ghost":        {{"Buff": {{"Vitality": 75, "Strength": 50}}}},
		"Harpy":        {{"Buff": {{"Strength": 50, "Speed": 50, "Vitality": 25}}, "Weakness": {{"Durability": 50}}}},
		"Skeleton":     {{"Weakness": {{"Strength": 25, "Vitality": 25, "Intelligence": 25, "Durability": 25, "Speed": 25, "Mastery": 25, "Mana": 25, "Height": 25}}}},
		"Djinn":        {{"Buff": {{"Strength": 50, "Vitality": 50, "Intelligence": 50, "Durability": 50, "Speed": 50, "Mastery": 50, "Mana": 50, "Height": 50}}}},
		"Humans":       {{"Buff": {{"Vitality": 25, "Intelligence": 25}}}},
		"Drawf":        {{"Buff": {{"Durability": 75, "Mastery": 75}}, "Weakness": {{"Speed": 50}}}},
	}
)
