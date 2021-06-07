package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//KillEvent data
type KillEvent struct {
	Type      string   `json:"type"`
	Timestamp int      `json:"timestamp"`
	Pos       Position `json:"position"`
	Killer    int      `json:"killerId"`
	Assist    []int    `json:"assistingParticipantIds"`
}

//Frames of events
type Frames struct {
	Events []KillEvent `json:"events"`
}

//GameEvent data for given match
type GameEvent struct {
	Frames []Frames `json:"frames"`
}

//Position of player
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type KAEvent struct {
	EType string
	X     int
	Y     int
}

func matchEventData(payload string, gameID int64) (eventdata GameEvent) {
	resp, err := http.Get("https://na1.api.riotgames.com/lol/match/v4/timelines/by-match/" + fmt.Sprint(gameID) + "?" + payload)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, &eventdata)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func killAssistLocale(name, payload string, gamedataArray []Game, maxtime int) (kaEvents []KAEvent) {
	for _, gamedata := range gamedataArray {
		println("GameId: " + fmt.Sprint(gamedata.GameID))
		var pid int
		for i, id := range gamedata.ParticipantIdentities {
			if strings.Compare(strings.ToUpper(strings.TrimSpace(string(id.Player.SummonerName))), strings.ToUpper(name)) == 0 {
				pid = i + 1
			}
		}

		eventdata := matchEventData(payload, gamedata.GameID)

		for _, i := range eventdata.Frames {
			for _, a := range i.Events {
				champKillBool := (a.Type == "CHAMPION_KILL")
				beforeTimer := (a.Timestamp <= maxtime)
				if !beforeTimer {
					break
				}
				if champKillBool {
					if fmt.Sprint(a.Killer) == fmt.Sprint(pid) {
						kX := a.Pos.X
						kY := a.Pos.Y
						killEvent := KAEvent{"KILL", kX, kY}
						kaEvents = append(kaEvents, killEvent)
					}
					for x := range a.Assist {
						if pid == a.Assist[x] {
							kX := a.Pos.X
							kY := a.Pos.Y
							assistEvent := KAEvent{"ASSIST", kX, kY}
							kaEvents = append(kaEvents, assistEvent)
						}
					}
				}
			}
		}
	}
	return
}

func kaLocaleAverage(kaEvents []KAEvent) (avgX, avgY int) {
	totalEvent := len(kaEvents)
	totalX := 0
	totalY := 0
	for _, event := range kaEvents {
		totalX += event.X
		totalY += event.Y
	}
	avgX = totalX / totalEvent
	avgY = totalY / totalEvent

	return
}
