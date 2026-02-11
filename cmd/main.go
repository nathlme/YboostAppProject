package main

import  (
	"fmt"
)

func main () {
	champsL:= GetChampion()

	champsSlice := SelecChamp(champsL)

    randomChamp := GetRandomChampion(champsSlice)

    fmt.Println(randomChamp.Name)
}




