package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//MatchHistory of target summoner
type MatchHistory struct {
	Matches    []Matches `json:"matches"`
	StartIndex int       `json:"startIndex"`
	EndIndex   int       `json:"endIndex"`
	TotalGames int       `json:"totalGames"`
}

//Matches detail of given summoner
type Matches struct {
	PlatformID string `json:"platformId"`
	GameID     int64  `json:"gameId"`
	Champion   int    `json:"champion"`
	Queue      int    `json:"queue"`
	Season     int    `json:"season"`
	Timestamp  int64  `json:"timestamp"`
	Role       string `json:"role"`
	Lane       string `json:"lane"`
}

func matchesByAcc(summoner, payload string) (list []int64) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/match/v4/matchlists/by-account/" + summoner + "?" + payload + "&endIndex=15")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var matchbody MatchHistory
	err = json.Unmarshal(body, &matchbody)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range matchbody.Matches {
		list = append(list, v.GameID)
	}
	return
}

func matchesByRole(accID, payload, lane string) (listByRole []int64) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/match/v4/matchlists/by-account/" + accID + "?" + payload + "&endIndex=20")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	println("https://na1.api.riotgames.com/lol/match/v4/matchlists/by-account/" + accID + "?" + payload + "&endIndex=20")
	var matchbody MatchHistory
	err = json.Unmarshal(body, &matchbody)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range matchbody.Matches {
		if v.Lane == strings.ToUpper(lane) {
			listByRole = append(listByRole, v.GameID)
		}
	}
	return
}
