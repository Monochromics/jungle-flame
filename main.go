package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

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

	kill, death, assist, champN, gamet := matchFeedCheck(ps.ByName("name"), payloadE, list)

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

	fmt.Println(list)

	tEvent, tX, tY := killAssistLocale(ps.ByName("name"), payloadE, list)
	avgX := tX / tEvent
	avgY := tY / tEvent
	fmt.Fprintf(w, "Average X:  "+fmt.Sprint(avgX)+"\n")
	fmt.Fprintf(w, "Average Y:  "+fmt.Sprint(avgY)+"\n")

	fmt.Println("Endpoint Hit: Avg Kill ")

}

func testJS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl, err := template.ParseFiles("./static/killmap.gohtml")
	if err != nil {
		panic(err)
	}

	type Coordinates struct {
		Coord string
	}

	entry := Coordinates{"[4123,1235]"}
	err = tmpl.Execute(w, entry)
	if err != nil {
		panic(err)
	}

	// http.ServeFile(w, r, "./static/killmap.gohtml")

}

func jungleLiveKL(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	payload := url.Values{}
	payload.Add("api_key", os.Getenv("APIKEY"))
	payloadE := payload.Encode()

	_, id := summonerByName(ps.ByName("name"), payloadE)
	junglerNames := liveGameJunglers(id, payloadE)

	for _, a := range junglerNames {
		println(a)
		jSumm, _ := summonerByName(a, payloadE)
		if jSumm == "" {
			println("FAILED TO GET SUMM ID FOR : " + a)
			break
		}

		time.Sleep(2 * time.Second)
		println(jSumm)
		matches := matchesByRole(jSumm, payloadE, "JUNGLE")
		tEvent, tX, tY := killAssistLocale(a, payloadE, matches)
		avgX := tX / tEvent
		avgY := tY / tEvent
		fmt.Fprintf(w, a+"\n")
		fmt.Fprintf(w, "Average X:  "+fmt.Sprint(avgX)+"\n")
		fmt.Fprintf(w, "Average Y:  "+fmt.Sprint(avgY)+"\n")
	}

}

func handleRequests() {

	router := httprouter.New()
	router.GET("/", homePage)
	router.GET("/kills/:name", avgKillCoord)
	router.GET("/summoner/:name", fuckingFeeder)
	router.GET("/jkl/:name", jungleLiveKL)
	router.GET("/test", testJS)

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	log.Fatal(http.ListenAndServe(":10000", router))
}

func main() {
	LoadConfiguration("api_conf.json")
	handleRequests()
}
