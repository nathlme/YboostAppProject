package main

type Champion struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ChampionResponse struct {
	Data map[string]Champion `json:"data"`
}


type Spell struct {
	ID string `json:"id"`
	Name string `json:"name"`

}

type SpellResponse struct{
	Spells map[string]Spell `json : "spells"`
}

type ChampionDetailResponse struct {
	Data map[string]SpellResponse `json:"data"`
}