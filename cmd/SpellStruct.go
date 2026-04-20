package main 
type Image struct {
	Full string `json:"full"`
}

type Spell struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       Image  `json:"image"`
}

type ChampionData struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Spells []Spell `json:"spells"`
}

type ChampionSpellResponse struct {
	Data map[string]ChampionData `json:"data"`
}

type DailySpell struct {
	ChampionID   string
	ChampionName string
	SpellSlot    string
	SpellID      string
	SpellName    string
	Description  string
	ImageURL     string
}