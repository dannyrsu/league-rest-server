package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	leagueapi "github.com/dannyrsu/league-api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

type server struct {
	router *chi.Mux
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the League of Draaaaven"))
}

func (*server) getSummonerStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryValues := r.URL.Query()

	summonerProfile := leagueapi.GetSummonerProfile(chi.URLParam(r, "summonername"), queryValues.Get("region"))

	results := map[string]interface{}{
		"summonerProfile": summonerProfile,
	}

	json.NewEncoder(w).Encode(results)
}

func (*server) getMatchDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryValues := r.URL.Query()
	matchID, err := strconv.ParseInt(chi.URLParam(r, "matchid"), 10, 64)
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

func (*server) getChampionByKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	champion := leagueapi.GetChampionByKey(chi.URLParam(r, "championkey"))

	json.NewEncoder(w).Encode(champion)
}

func (s *server) middleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Use(middleware.Timeout(60 * time.Second))
}

func (s *server) routes() {
	s.router.Get("/", defaultHandler)
	s.router.Get("/v1/summoner/{summonername}/stats", s.getSummonerStatsHandler)
	s.router.Get("/v1/match/{matchid}", s.getMatchDetailHandler)
	s.router.Get("/v1/champion/{championkey}", s.getChampionByKeyHandler)
}

func main() {
	server := &server{
		router: chi.NewRouter(),
	}

	server.middleware()
	server.routes()

	handler := cors.Default().Handler(server.router)
	log.Println("Starting League Server...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
