package main

type ChampionCard struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	Lore string `json:"lore"`
}

type Champion struct {	
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	Lore string `json:"blurb"`
}

type ChampionResponse struct {
	Data map[string]Champion `json:"data"`
}


type GuessRequest struct {
	Champion string `json:"champion"`
}

type GuessResponse struct  {
	Same bool `json:"correct"`
}

// type Spell struct {
// 	ID string `json:"id"`
// 	Name string `json:"name"`

// }

// type SpellResponse struct{
// 	Spells map[string]Spell `json:"spells"`
// }

// type ChampionDetailResponse struct {
// 	Data map[string]SpellResponse `json:"data"`
// }