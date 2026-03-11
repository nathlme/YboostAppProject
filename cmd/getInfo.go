package main

import(
	"encoding/json"
	"fmt"
	"io"
	"net/http"

)
var championCards []ChampionCard

func GetVersion(){
	url := "https://ddragon.leagueoflegends.com/api/versions.json"

	resp, err := http.Get(url)
	if err != nil {
			fmt.Println("Erreur lors du GET :", err)
			return 
		}

	defer resp.Body.Close()


	body,err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du Body:", err)
		return
	}

	var versions []string
	
	err = json.Unmarshal(body,&versions)
	if err != nil {
		fmt.Println("Erreur JSON :", err)
		return
	}

	fmt.Println("Version la plus récente :", versions[0])
}


func GetChampion() map[string]Champion {
	url := "https://ddragon.leagueoflegends.com/cdn/16.2.1/data/en_US/champion.json"

	resp,err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors du GET :", err) 
	}

	defer resp.Body.Close()
		
	body, err := io.ReadAll(resp.Body)
		if err != nil {
		fmt.Println("Erreur lors de la lecture du Body:", err)
		
	}

	var result ChampionResponse
	
	err = json.Unmarshal(body,&result)
	if err != nil {
		fmt.Println("Error JSON :", err)
		
	}

	return result.Data

}

func LoadChampionCards(){
	
	championCards = []ChampionCard{}

	Dico := GetChampion()

	for _, champ := range Dico {
		var card ChampionCard
		card.ID = champ.ID
		card.Name = champ.Name
		card.Lore = champ.Lore
		card.ImageUrl = BuildUrl("https://ddragon.leagueoflegends.com/cdn/16.2.1/img/champion/", card.ID)
		championCards = append(championCards,card)
	}
}

func BuildUrl(baseUrl string, champId string) string {
	newUrl := baseUrl + champId + ".png"
	return newUrl
}


















// func GetSpells(){
// 	url := "https://ddragon.leagueoflegends.com/cdn/16.2.1/data/en_US/champion/Ahri.json"

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println("Erreur lors du GET :", err)
// 		return 
// 	}

// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 		fmt.Println("Erreur lors de la lecture du Body:", err)
// 		return
// 	}

// }