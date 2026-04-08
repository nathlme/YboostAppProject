package main 


type Spell struct {
	ID				string	`json:"id"`
	Name			string	`json:"name"`
	Description 	string 	`json:"description"`
	ImageUrl		Image	`json:"image"`
}

type SpellCard struct {
	ID				string	`json:"id"`
	Name			string	`json:"name"`
	Description 	string 	`json:"description"`
	ImageUrl		Image	`json:"image"`
}

type SpellResponse struct {
	Data map[string]Spell `json:"data"`
}

type Image struct {
	Full string `json:"full"`
}