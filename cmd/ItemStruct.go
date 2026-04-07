package main 


type Item struct {
	Name  string            `json:"name"`
	Maps  map[string]bool   `json:"maps"`
	Gold  GoldData          `json:"gold"`
	Image ImageData         `json:"image"`
	Stats map[string]float64 `json:"stats"`
	SpecialRecipe	int		`json:"specialRecipe"`
	From []string			`json:"from"`
}


type ItemCard struct {
	Name 	 	string  			`json:"name"`
	ImageId	 	string  			`json:"full"`
	ImageUrl 	string  			`json:"imageUrl"`
	Price		int					`json:"price"`
	Stats 		map[string]float64	`json:"stats"`
	SpecialRecipe	int		 		`json:"specialRecipe"`
	From          []string           `json:"from"`
}


type GoldData struct {
	Base        int  `json:"base"`
	Total       int  `json:"total"`
	Sell        int  `json:"sell"`
	Purchasable bool `json:"purchasable"`
}


type ItemComponent struct {
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}


type ImageData struct {
	Full string `json:"full"`
}


type ItemResponse struct {
	Data	map[string]Item	  `json:"data"`
}


type ItemGuessRequest struct {
	Item string `json:"item"`
}