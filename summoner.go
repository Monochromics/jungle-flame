package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//BySummoner search by summoner name return
type BySummoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

func summonerByName(name, payload string) (accountID, summID string) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/summoner/v4/summoners/by-name/" + name + "?" + payload)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var sbody BySummoner

	err = json.Unmarshal(body, &sbody)
	if err != nil {
		log.Fatalln(err)
	}
	accountID = sbody.AccountID
	summID = sbody.ID
	return
}
