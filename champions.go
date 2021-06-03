package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//Champion name
type Champion struct {
	Version string `json:"version"`
	ID      string `json:"id"`
	Key     string `json:"key"`
	Name    string `json:"name"`
}

//ChampData lookup
type ChampData struct {
	Type    string              `json:"type"`
	Format  string              `json:"format"`
	Version string              `json:"version"`
	Data    map[string]Champion `json:"data"`
}

//CLookup using to get champ info
func CLookup(champID int) (champName string) {
	//Champion grab
	resp, err := http.Get("https://ddragon.leagueoflegends.com/cdn/10.19.1/data/en_US/champion.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var cdata ChampData
	json.Unmarshal(body, &cdata)

	for i := range cdata.Data {
		if cdata.Data[i].Key == fmt.Sprint(champID) {
			champName = cdata.Data[i].Name
		}
	}
	return
}
