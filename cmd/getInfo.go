package main

import(
	"encoding/json"
	"fmt"
	"io"
	"net/http"

)


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


func GetChampion() {
	url := "https://ddragon.leagueoflegends.com/cdn/16.2.1/data/en_US/champion.json"

	resp,err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors du GET :", err)
		return 
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
		if err != nil {
		fmt.Println("Erreur lors de la lecture du Body:", err)
		return
	}

	var result ChampionResponse
	
	err = json.Unmarshal(body,&result)
	if err != nil {
		fmt.Println("Error JSON :", err)
		return
	}

	count := 0 
	for key, champ := range result.Data{
		fmt.Printf("%s -> %s (%v)\n",key,champ.Name,champ.Tags)
		count++
		if count == 5{
			break
		}
	}
	
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