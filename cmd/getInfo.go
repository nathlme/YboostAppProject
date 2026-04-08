package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)


// A list of Champion object 
var championCards []ChampionCard
// A list of Item object 
var itemCards []ItemCard


// Print the latest version of the League of legends API 
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


//	Return a map of the ChampionResponse struct fill with every champion in the Lol API response 
func GetChampion() map[string]Champion {
	url := "https://ddragon.leagueoflegends.com/cdn/16.2.1/data/fr_FR/champion.json"

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


// Fill an object of the struct ChampionCard and add it to the championCards list
func LoadChampionCards(){
	
	championCards = []ChampionCard{}

	Dico := GetChampion()

	for _, champ := range Dico {
		var card ChampionCard
		card.ID = champ.ID
		card.Name = champ.Name
		card.Lore = champ.Lore
		card.ImageUrl = BuildUrl("https://ddragon.leagueoflegends.com/cdn/16.2.1/img/champion/", card.ID)
		championCards = append(championCards, card)
	}
}


// Build an url based on a base and a champion id and return it
func BuildUrl(baseUrl string, champId string) string {
	var newUrl = ""

	if strings.HasSuffix(champId, ".png") {
		newUrl = baseUrl + champId
	} else {
		newUrl = baseUrl + champId + ".png"
	}
	
	return newUrl
}


// Return the name, the title and the complet lore in French of the daily champ 
func GetFullLore() (string, string, string, error) {

	dayliChamp, err := GetOrCreateDailyChampion()
	if err != nil {
		return "", "", "", err
	}

	lang := "fr_FR"

	championUrl := fmt.Sprintf(
		"https://ddragon.leagueoflegends.com/cdn/16.2.1/data/%s/champion/%s.json",
		lang,
		dayliChamp,
	)

	resp, err := http.Get(championUrl)
	if err != nil {
		return "", "", "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	var champDetails ChampionDetail

	err = json.Unmarshal(body, &champDetails)
	if err != nil {
		return "", "", "", err
	}

	champ, ok := champDetails.Data[dayliChamp]
	if !ok {
		return "", "", "", fmt.Errorf("champion introuvable dans les données JSON")
	}

	return champ.Name, champ.Title, champ.Lore, nil
}


// Return a map of the ItemResponse struct fill with every item in the Lol API response 
func GetItem() (map[string]Item, error) {
	url := "https://ddragon.leagueoflegends.com/cdn/16.2.1/data/fr_FR/item.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var itemResult ItemResponse

	err = json.Unmarshal(body, &itemResult)
	if err != nil {
		return nil, err
	}

	return itemResult.Data, nil
}


// Fill an object of the struct ItemCard and add it to the itemCards list
func LoadItemCard() error {
	itemCards = []ItemCard{}

	dico, err := GetItem()
	if err != nil {
		return err
	}

	for _, item := range dico {
		if !IsValidChallengeItem(item) {
			continue
		}

		var card ItemCard

		card.Name = item.Name
		card.ImageId = item.Image.Full
		card.ImageUrl = BuildUrl("https://ddragon.leagueoflegends.com/cdn/16.2.1/img/item/", card.ImageId)
		card.Price = item.Gold.Total
		card.Stats = item.Stats
		card.SpecialRecipe = item.SpecialRecipe
		card.From = item.From 
		
		itemCards = append(itemCards, card)
	}

	if len(itemCards) == 0 {
		return fmt.Errorf("aucun item valide chargé")
	}

	log.Printf("%d items valides chargés\n", len(itemCards))
	return nil
}


// Return a list of the object component 
func GetItemComponents(fromIDs []string, dico map[string]Item) []ItemComponent {
	components := []ItemComponent{}

	for _, id := range fromIDs {
		if comp, ok := dico[id]; ok {
			components = append(components, ItemComponent{
				Name:     comp.Name,
				ImageUrl: BuildUrl("https://ddragon.leagueoflegends.com/cdn/16.2.1/img/item/", comp.Image.Full),
			})
		}
	}

	return components
}