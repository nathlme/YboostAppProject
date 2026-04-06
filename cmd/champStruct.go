package main

type ChampionCard struct {
	ID 		 string `json:"id"`
	Name 	 string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	Lore	 string `json:"lore"`
}

type Champion struct {	
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	Lore string `json:"blurb"`
}

type ChampionDetail struct {
    Data map[string]struct {
        Name  string `json:"name"`
        Title string `json:"title"`
        Lore  string `json:"lore"`
    } `json:"data"`
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


type Item struct {
	Name  string            `json:"name"`
	Maps  map[string]bool   `json:"maps"`
	Gold  GoldData          `json:"gold"`
	Image ImageData         `json:"image"`
	Stats map[string]float64 `json:"stats"`
	SpecialRecipe	int		 `json:"specialRecipe"`
}

type GoldData struct {
	Base        int  `json:"base"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
	Purchasable bool `json:"purchasable"`
}

type ImageData struct {
	Full string `json:"full"`
}

type ItemCard struct {
	Name 	 	string  			`json:"name"`
	ImageId	 	string  			`json:"full"`
	ImageUrl 	string  			`json:"imageUrl"`
	Price		int					`json:"base"`
	Stats 		map[string]float64	`json:"stats"`
	SpecialRecipe	int		 		`json:"specialRecipe"`
}

type ItemResponse struct {
	Data	map[string]Item	  `json:"data"`
}

type ItemGuessRequest struct {
	Item string `json:"item"`
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