package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//BySummoner search by summoner name return
type Summoner struct {
	ID        string `json:"id"`
	AccountID string `json:"accountId"`
	Puuid     string `json:"puuid"`
	Name      string `json:"name"`
}

func summonerByName(name, payload string) (summoner Summoner) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + name + "?" + payload)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &summoner)
	if err != nil {
		log.Fatalln(err)
	}
	return
}
