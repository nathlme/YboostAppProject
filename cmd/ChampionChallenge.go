package main

import (
	"math/rand"
	"time"
)

func SelecChamp(champs map[string]Champion) []Champion{
	var list []Champion 
	for _,c := range champs{
		list = append(list,c)
		print(c.ID)
	}
	return list 
}

func GetRandomChampion(list []Champion) Champion {
    rand.Seed(time.Now().UnixNano())
    return list[rand.Intn(len(list))]
}