package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rs/cors"

	leagueapi "github.com/dannyrsu/league-api"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	router *httprouter.Router
}

func (*server) defaultHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to the League of Draaaaven")
}

func (*server) getSummonerStatsHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	queryValues := r.URL.Query()

	summonerProfile := leagueapi.GetSummonerProfile(params.ByName("summonername"), queryValues.Get("region"))

	results := map[string]interface{}{
		"summonerProfile": summonerProfile,
	}

	json.NewEncoder(w).Encode(results)
}

func (*server) getMatchDetailHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	queryValues := r.URL.Query()
	matchID, err := strconv.ParseInt(params.ByName("matchid"), 10, 64)
	if err != nil {
		log.Fatalf("Error converting match paramter: %v", err)
		matchID = 0
	}
	match := leagueapi.GetGameData(matchID, queryValues.Get("region"))

	results := map[string]interface{}{
		"match": match,
	}

	json.NewEncoder(w).Encode(results)
}

func (*server) getChampionByKeyHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	champion := leagueapi.GetChampionByKey(params.ByName("championkey"))

	json.NewEncoder(w).Encode(champion)
}

func (s *server) routes() {
	s.router.GET("/", s.defaultHandler)
	s.router.GET("/v1/summoner/:summonername/stats", s.getSummonerStatsHandler)
	s.router.GET("/v1/match/:matchid", s.getMatchDetailHandler)
	s.router.GET("/v1/champion/:championkey", s.getChampionByKeyHandler)
	s.router.ServeFiles("/static/*filepath", http.Dir("static"))
}

func main() {
	server := &server{
		router: httprouter.New(),
	}

	server.routes()

	handler := cors.Default().Handler(server.router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
