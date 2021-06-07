package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/julienschmidt/httprouter"
)

// Config file containing API key
type Config struct {
	RiotKey string
}

// LoadConfiguration loads json key file
func LoadConfiguration(file string) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	os.Setenv("APIKEY", config.RiotKey)
}

func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome! Add '/summoner/$NASummonerName'  to the end of the url to get started!\n")
	fmt.Println("Endpoint Hit: HomePage")
}

func fuckingFeeder(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload := url.Values{}
	payload.Add("api_key", os.Getenv("APIKEY"))
	payloadE := payload.Encode()
	summoner, _ := summonerByName(ps.ByName("name"), payloadE)
	list := matchesByAcc(summoner, payloadE)
	matchData := matchDataGrab(payloadE, list)

	kill, death, assist, champN, gamet := matchFeedCheck(ps.ByName("name"), matchData)

	fmt.Fprintf(w, "Last time "+ps.ByName("name")+" fed:\n")
	fmt.Fprintf(w, fmt.Sprint(kill)+"/"+fmt.Sprint(death)+"/"+fmt.Sprint(assist)+"     "+champN+"\n")
	fmt.Fprintf(w, gamet+"\n")

	fmt.Println("Endpoint Hit: Summoner")
}

func avgKillCoord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload := url.Values{}
	payload.Add("api_key", os.Getenv("APIKEY"))
	payloadE := payload.Encode()
	summoner, _ := summonerByName(ps.ByName("name"), payloadE)
	list := matchesByAcc(summoner, payloadE)
	gamedataArray := matchDataGrab(payloadE, list)
	tEvent, tX, tY := killAssistLocale(ps.ByName("name"), payloadE, gamedataArray)
	avgX := tX / tEvent
	avgY := tY / tEvent

	coordArr := [2]int{avgX, avgY}
	soloKillMap(w, r, ps, coordArr[:])
}

func soloKillMap(w http.ResponseWriter, r *http.Request, ps httprouter.Params, coord []int) {
	tmpl, err := template.ParseFiles("./static/single/killmap.gohtml")
	if err != nil {
		panic(err)
	}

	type Coordinates struct {
		Coord string
	}
	s, _ := json.Marshal(coord)
	entry := Coordinates{string(s)}
	err = tmpl.Execute(w, entry)
	if err != nil {
		panic(err)
	}
}

func jungleLiveKL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload := url.Values{}
	payload.Add("api_key", os.Getenv("APIKEY"))
	payloadE := payload.Encode()

	_, id := summonerByName(ps.ByName("name"), payloadE)
	junglerNames := liveGameJunglers(id, payloadE)

	type JunglerData struct {
		JunglerName string
		Coords      []int
	}

	dataArray := [2]JunglerData{}
	for i, a := range junglerNames {
		println(a)
		jSumm, _ := summonerByName(a, payloadE)
		if jSumm == "" {
			println("FAILED TO GET SUMM ID FOR : " + a)
			break
		}
		println(jSumm)
		matches := matchesByRole(jSumm, payloadE, "JUNGLE")
		matchData := matchDataGrab(payloadE, matches)
		tEvent, tX, tY := killAssistLocale(a, payloadE, matchData)
		avgX := tX / tEvent
		avgY := tY / tEvent
		coords := []int{avgX, avgY}
		out := JunglerData{JunglerName: a, Coords: coords}
		dataArray[i] = out
	}
	jungleKillMaps(w, r, ps, dataArray[0].JunglerName, dataArray[0].Coords, dataArray[1].JunglerName, dataArray[1].Coords)

}

func jungleKillMaps(w http.ResponseWriter, r *http.Request, ps httprouter.Params, jungleA string, coordA []int, jungleB string, coordB []int) {
	tmpl, err := template.ParseFiles("./static/junglers/killmap.gohtml")
	if err != nil {
		panic(err)
	}

	type Coordinates struct {
		JungleA string
		JungleB string
		CoordA  string
		CoordB  string
	}

	cUA, _ := json.Marshal(coordA)
	cUB, _ := json.Marshal(coordB)
	entry := Coordinates{JungleA: jungleA, JungleB: jungleB, CoordA: string(cUA), CoordB: string(cUB)}
	err = tmpl.Execute(w, entry)
	if err != nil {
		panic(err)
	}
}

func handleRequests() {

	router := httprouter.New()
	router.GET("/", homePage)
	router.GET("/kills/:name", avgKillCoord)
	router.GET("/summoner/:name", fuckingFeeder)
	router.GET("/jkl/:name", jungleLiveKL)

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	LoadConfiguration("api_conf.json")
	handleRequests()
}
