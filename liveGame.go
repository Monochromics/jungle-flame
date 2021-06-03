package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// LiveGame of given encrypted summoner ID
type LiveGame struct {
	GameID            int64  `json:"gameId"`
	MapID             int    `json:"mapId"`
	GameMode          string `json:"gameMode"`
	GameType          string `json:"gameType"`
	GameQueueConfigID int    `json:"gameQueueConfigId"`
	Participants      []struct {
		TeamID        int    `json:"teamId"`
		Spell1ID      int    `json:"spell1Id"`
		Spell2ID      int    `json:"spell2Id"`
		ChampionID    int    `json:"championId"`
		ProfileIconID int    `json:"profileIconId"`
		SummonerName  string `json:"summonerName"`
		Bot           bool   `json:"bot"`
		SummonerID    string `json:"summonerId"`
		Perks         struct {
			PerkIds      []int `json:"perkIds"`
			PerkStyle    int   `json:"perkStyle"`
			PerkSubStyle int   `json:"perkSubStyle"`
		} `json:"perks"`
	} `json:"participants"`
	Observers struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`
	PlatformID      string `json:"platformId"`
	BannedChampions []struct {
		ChampionID int `json:"championId"`
		TeamID     int `json:"teamId"`
		PickTurn   int `json:"pickTurn"`
	} `json:"bannedChampions"`
	GameStartTime int64 `json:"gameStartTime"`
	GameLength    int   `json:"gameLength"`
}

func liveGameJunglers(summoner, payload string) (junglerRealName []string) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/" + summoner + "?" + payload)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var livegamedata LiveGame
	err = json.Unmarshal(body, &livegamedata)
	if err != nil {
		log.Fatalln(err)
	}

	println(livegamedata.GameType)

	for i := range livegamedata.Participants {
		switch {
		case livegamedata.Participants[i].Spell1ID == 11:
			junglerRealName = append(junglerRealName, livegamedata.Participants[i].SummonerName)
		case livegamedata.Participants[i].Spell2ID == 11:
			junglerRealName = append(junglerRealName, livegamedata.Participants[i].SummonerName)
		}
	}
	return
}
