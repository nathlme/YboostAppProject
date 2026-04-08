package main


type Champion struct {	
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	Lore string   `json:"blurb"`
}


type ChampionCard struct {
	ID 		 string `json:"id"`
	Name 	 string `json:"name"`
	ImageUrl string `json:"imageUrl"`
	Lore	 string `json:"lore"`
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
