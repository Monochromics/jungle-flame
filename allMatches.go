package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto"
)

//Stats of given game
type Stats struct {
	ParticipantID    int `json:"participantId"`
	Kills            int `json:"kills"`
	Deaths           int `json:"deaths"`
	Assists          int `json:"assists"`
	TotalPlayerScore int `json:"totalPlayerScore"`
	TotalScoreRank   int `json:"totalScoreRank"`
}

//Participants in given game
type Participants struct {
	ParticipantID int   `json:"participantId"`
	TeamID        int   `json:"teamId"`
	ChampionID    int   `json:"championId"`
	Spell1ID      int   `json:"spell1Id"`
	Spell2ID      int   `json:"spell2Id"`
	Stats         Stats `json:"stats"`
}

//Player in given game
type Player struct {
	PlatformID        string `json:"platformId"`
	AccountID         string `json:"accountId"`
	SummonerName      string `json:"summonerName"`
	SummonerID        string `json:"summonerId"`
	CurrentPlatformID string `json:"currentPlatformId"`
	CurrentAccountID  string `json:"currentAccountId"`
	MatchHistoryURI   string `json:"matchHistoryUri"`
	ProfileIcon       int    `json:"profileIcon"`
}

//ParticipantIdentities in given game
type ParticipantIdentities struct {
	ParticipantID int    `json:"participantId"`
	Player        Player `json:"player"`
}

//Game of league in player's history
type Game struct {
	GameID                int64                   `json:"gameId"`
	PlatformID            string                  `json:"platformId"`
	GameCreation          int64                   `json:"gameCreation"`
	GameDuration          int                     `json:"gameDuration"`
	QueueID               int                     `json:"queueId"`
	MapID                 int                     `json:"mapId"`
	SeasonID              int                     `json:"seasonId"`
	GameVersion           string                  `json:"gameVersion"`
	GameMode              string                  `json:"gameMode"`
	GameType              string                  `json:"gameType"`
	ParticipantIdentities []ParticipantIdentities `json:"participantIdentities"`
	Participants          []Participants          `json:"participants"`
}

func matchDataGrab(payload string, gameID []int64, c *ristretto.Cache) (gamedataArray []Game) {
	for _, s := range gameID {
		var gamedata Game

		value, found := c.Get(fmt.Sprint(s))

		if found {
			gamedata = value.(Game)
		} else {
			resp, err := http.Get("https://na1.api.riotgames.com/lol/match/v4/matches/" + fmt.Sprint(s) + "?" + payload)
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			err = json.Unmarshal(body, &gamedata)
			if err != nil {
				log.Fatalln(err)
			}
			c.Set(fmt.Sprint(s), gamedata, 1)
			println(fmt.Sprint(gamedata.GameVersion))
		}
		gamedataArray = append(gamedataArray, gamedata)
	}
	return
}

func matchFeedCheck(name string, gamedataArray []Game) (killz, deathz, assistz int, champName, gtime string) {
	for _, gamedata := range gamedataArray {
		var pid int
		for i := range gamedata.ParticipantIdentities {
			if strings.Compare(strings.ToUpper(strings.TrimSpace(string(gamedata.ParticipantIdentities[i].Player.SummonerName))), strings.ToUpper(name)) == 0 {
				pid = i
				break
			}
		}

		killz = gamedata.Participants[pid].Stats.Kills
		deathz = gamedata.Participants[pid].Stats.Deaths
		assistz = gamedata.Participants[pid].Stats.Assists

		if deathz >= (killz + assistz) {
			t := time.Unix(0, gamedata.GameCreation*int64(1000000))
			gtime = t.Format("01-02-2006 15:04")
			champName = CLookup(gamedata.Participants[pid].ChampionID)
			break
		}
	}
	return
}
