package generator

type WeightedOption struct {
	Value  string
	Weight int
}

var (
	RacesRarity = map[string][]string{
		"Common":    {"Humans", "Gnome", "Orc", "Giant", "Kobold", "Goblin", "Skeleton"},
		"Uncommon":  {"Drawf", "Elf", "Centaur", "Minotaur", "Cyclops", "Mushfolk", "Beastfolk", "Lamia", "Undead", "Harpy"},
		"Rare":      {"Dullahan", "Merfolk", "Fairy", "Druid", "Vampire", "Werewolf", "Ghost"},
		"Epic":      {"Demon", "Angel", "Djinn", "Wizard/Witch"},
		"Legandary": {"God", "Dragonborn"},
	}

	VitalityRarity = map[string][]string{
		"Common":    {"Average"},
		"Uncommon":  {"Weak", "Heathly"},
		"Rare":      {"Frail", "Robust"},
		"Epic":      {"Helpless", "Vigorous"},
		"Legandary": {"Helpless-", "Vigorous+"},
	}

	SpeedRarity = map[string][]string{
		"Common":    {"Average"},
		"Uncommon":  {"Slow", "Fast"},
		"Rare":      {"Sluggish", "Accelerated"},
		"Epic":      {"Crippled", "Supersonic"},
		"Legandary": {"Torid", "Hypersonic"},
	}

	StrengthRarity = map[string][]string{
		"Common":    {"Average"},
		"Uncommon":  {"Weak", "Strong"},
		"Rare":      {"Scrwny", "Formidable"},
		"Epic":      {"Forceless", "Overpowering"},
		"Legandary": {"Forceless-", "Overpowering+"},
	}

	DurabilityRarity = map[string][]string{
		"Common":    {"Average"},
		"Uncommon":  {"Vincible", "Reinforced"},
		"Rare":      {"Vulnerable", "Armored"},
		"Epic":      {"Defenseless", "Fortified"},
		"Legandary": {"Defenseless-", "Fortified+"},
	}

	IntelligenceRarity = map[string][]string{
		"Common":    {"Average"},
		"Uncommon":  {"Dumb", "Smart"},
		"Rare":      {"Lobotomized", "Genius"},
		"Epic":      {"Mindless", "Prodigious"},
		"Legandary": {"Mindless-", "Prodigious+"},
	}

	ManaFlowRarity = map[string][]string{
		"Common":    {"Average"},
		"Uncommon":  {"Hexed", "Enchanted"},
		"Rare":      {"Lowly", "Conjuring"},
		"Epic":      {"Mana-Less", "Overflowing"},
		"Legandary": {"No Mana", "Overflowing+"},
	}

	SkillLevelRarity = map[string][]string{
		"Common":    {"Average"},
		"Uncommon":  {"Amateur", "Skilled"},
		"Rare":      {"Novice", "Expert"},
		"Epic":      {"Skill-less", "Mastered"},
		"Legandary": {"Skill-less-", "Mastered+"},
	}

	ExtraTraitRarity = map[string][]string{
		"Common":    {"None", "Swift", "Quick Thinker", "Rough Skin", "Castle Training"},
		"Uncommon":  {"Fast Learner", "Abounding Flow", "Big Boned"},
		"Rare":      {"Druid's Blessing", "Naturally Skilled"},
		"Epic":      {"Call of Hercules", "Speed Force"},
		"Legandary": {"Blessed", "Isekai Protag"},
	}

	WeaknessOptions = []WeightedOption{
		{Value: "None", Weight: 40},
		{Value: "Fragile Bone", Weight: 10},
		{Value: "STD", Weight: 10},
		{Value: "Cancer", Weight: 10},
		{Value: "Delayed Reaction", Weight: 10},
		{Value: "Testicuilar Torsion", Weight: 10},
		{Value: "Amputee", Weight: 10},
		{Value: "Blindness", Weight: 10},
		{Value: "Too Young", Weight: 10},
		{Value: "Too Old", Weight: 10},
		{Value: "One Eye", Weight: 10},
		{Value: "Lobotomized", Weight: 10},
		{Value: "Auto Immune Disease", Weight: 10},
		{Value: "Claustrophobia", Weight: 10},
		{Value: "Paranoid", Weight: 10},
		{Value: "Schizophrenia", Weight: 10},
		{Value: "Cursed", Weight: 10},
	}

	WeaponOptions = []WeightedOption{
		{Value: "None", Weight: 20},
		{Value: "Basic Sword", Weight: 10},
		{Value: "Excalibur", Weight: 10},
		{Value: "Bow", Weight: 10},
		{Value: "Crossbow", Weight: 10},
		{Value: "Flintlock", Weight: 10},
		{Value: "Bardiche", Weight: 10},
		{Value: "Spear", Weight: 10},
		{Value: "Rapier", Weight: 10},
		{Value: "Shield", Weight: 10},
		{Value: "Sword and Shield", Weight: 10},
		{Value: "Whip", Weight: 10},
		{Value: "Anchor", Weight: 10},
		{Value: "Mace", Weight: 10},
		{Value: "Dagger", Weight: 10},
		{Value: "War Hammer", Weight: 10},
		{Value: "Battle Axe", Weight: 10},
		{Value: "Glaive", Weight: 10},
		{Value: "Scythe", Weight: 10},
		{Value: "Twinblade", Weight: 10},
		{Value: "Cutlass", Weight: 10},
		{Value: "Club", Weight: 10},
		{Value: "Whole Ass Tree Log", Weight: 10},
		{Value: "Katana", Weight: 10},
		{Value: "Big Ass Rock", Weight: 10},
		{Value: "Halberd", Weight: 10},
		{Value: "Sickle", Weight: 10},
		{Value: "SlingShot", Weight: 10},
		{Value: "Stake", Weight: 10},
		{Value: "MorningStar", Weight: 10},
		{Value: "Quarter Staff", Weight: 10},
		{Value: "Spiked Club", Weight: 10},
		{Value: "Lance", Weight: 10},
		{Value: "Bec De Corbin", Weight: 10},
		{Value: "Short Sword", Weight: 10},
		{Value: "Flail", Weight: 10},
		{Value: "Caestus", Weight: 10},
		{Value: "Magic Wand", Weight: 10},
		{Value: "Magic Staff", Weight: 10},
		{Value: "Magic Grimoire", Weight: 10},
	}

	XFactorOptions = []WeightedOption{
		{Value: "None", Weight: 20},
		{Value: "Tarnished", Weight: 10},
		{Value: "Elemental", Weight: 10},
		{Value: "Companionship", Weight: 10},
		{Value: "Partners in Crime", Weight: 10},
		{Value: "Avatar Of Elements", Weight: 10},
		{Value: "Weapon Smith", Weight: 10},
		{Value: "Training of Ten Ten", Weight: 10},
		{Value: "Halfling", Weight: 10},
		{Value: "Aizen's Plan", Weight: 10},
		{Value: "Growth Spurt", Weight: 10},
		{Value: "Naturally Buffed", Weight: 10},
	}

	ElementOptions = []WeightedOption{
		{Value: "None", Weight: 25},
		{Value: "Fire", Weight: 10},
		{Value: "Water", Weight: 10},
		{Value: "Earth", Weight: 10},
		{Value: "Wind", Weight: 10},
		{Value: "Nature", Weight: 8},
		{Value: "Toxic", Weight: 8},
		{Value: "Lighting", Weight: 5},
		{Value: "Sound", Weight: 8},
		{Value: "Dark", Weight: 5},
		{Value: "Light", Weight: 8},
		{Value: "Frost", Weight: 8},
		{Value: "Gravity", Weight: 5},
		{Value: "Crystal", Weight: 5},
		{Value: "Arcane", Weight: 5},
		{Value: "Time", Weight: 5},
	}

	AlignmentOptions = []WeightedOption{
		{Value: "Civilian", Weight: 5},
		{Value: "Anti-Hero", Weight: 1},
		{Value: "Villain", Weight: 2},
		{Value: "Hero", Weight: 2},
	}

	CompanionOptions = []string{
		"Small Dragon",
		"Dragon",
		"Wyvern",
		"Unicorn",
		"Horse",
		"Phoenix",
		"Giant Bat",
		"Hexed Cat",
		"Dog",
		"Slave",
		"Falcon",
		"Griffin",
		"Rat",
		"Fairy Spirit",
		"Snake",
		"Jackalope",
		"Bard",
		"Knight",
		"Wizard/Witch",
		"Imp",
		"Chocobo",
		"Owl",
		"Giant Raptor",
		"Big Cat",
		"Talking Object",
		"Finrir",
		"Fox",
		"Leprechaun",
	}

	HeightOptions = []string{
		"4'8",
		"4'9",
		"4'10",
		"4'11",
		"5'0",
		"5'1",
		"5'2",
		"5'3",
		"5'4",
		"5'5",
		"5'6",
		"5'7",
		"5'8",
		"5'9",
		"5'10",
		"5'11",
		"6'0",
		"6'1",
		"6'2",
		"6'3",
		"6'4",
		"6'5",
		"6'6",
		"6'7",
		"6'8",
		"6'9",
		"6'10",
		"6'11",
		"7'0",
		"7'1",
		"7'2",
		"7'3",
		"7'4",
		"7'5",
		"7'6",
		"7'7",
		"7'8",
		"7'9",
		"7'10",
		"7'11",
		"8'0",
	}
)
