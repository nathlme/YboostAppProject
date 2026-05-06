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

var cachedVersion string



// Print the latest version of the League of legends API 
func GetVersion() string{
	if cachedVersion != "" {
        return cachedVersion
    }

	url := "https://ddragon.leagueoflegends.com/api/versions.json"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors du GET :", err)
		return ""
	}

	defer resp.Body.Close()


	body,err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du Body:", err)
		return ""
	}

	var versions []string
	
	err = json.Unmarshal(body,&versions)
	if err != nil {
		fmt.Println("Erreur JSON :", err)
		return ""
	}
	

	cachedVersion = versions[0]
    return cachedVersion
}


//	Return a map of the ChampionResponse struct fill with every champion in the Lol API response 
func GetChampion() (map[string]Champion, error) {
	version := GetVersion()
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/fr_FR/champion.json", version)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du GET: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code inattendu: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du body: %w", err)
	}

	var result ChampionResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("erreur JSON: %w", err)
	}

	return result.Data, nil
}


// Fill an object of the struct ChampionCard and add it to the championCards list
func LoadChampionCards() error{
	
	championCards = []ChampionCard{}

	Dico, err := GetChampion()
	if err != nil {
		return err
	}

	version := GetVersion()
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/img/champion/",version)

	for _, champ := range Dico {
		var card ChampionCard
		card.ID = champ.ID
		card.Name = champ.Name
		card.Lore = champ.Lore
		card.ImageUrl = BuildUrl(url, card.ID)
		championCards = append(championCards, card)
	}
	
	return nil 
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

	dayliChamp, err := GetOrCreateDailyChampion("daily_champion")
	if err != nil {
		return "", "", "", err
	}

	lang := "fr_FR"
	version := GetVersion()
	championUrl := fmt.Sprintf(
		"https://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion/%s.json",
		version,
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
	version := GetVersion()
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/fr_FR/item.json", version)

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

	version := GetVersion()
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/img/item/",version)

	for _, item := range dico {
		if !IsValidChallengeItem(item) {
			continue
		}

		var card ItemCard

		card.Name 			= item.Name
		card.ImageId 		= item.Image.Full
		card.ImageUrl 		= BuildUrl(url, card.ImageId)
		card.Price		    = item.Gold.Total
		card.Stats 			= item.Stats
		card.SpecialRecipe 	= item.SpecialRecipe
		card.From 			= item.From 
		
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

	version := GetVersion()
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/img/item/",version)

	for _, id := range fromIDs {
		if comp, ok := dico[id]; ok {
			components = append(components, ItemComponent{
				Name:     comp.Name,
				ImageUrl: BuildUrl(url, comp.Image.Full),
			})
		}
	}

	return components
}