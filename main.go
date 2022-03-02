package main

import (
	_ "embed"
	"os"

	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/billglover/now-playing/lastfm"
)

//go:embed templates/albums.html
var indexTemplate string

type app struct {
	lfm  *lastfm.LFM
	page *template.Template
}

func main() {

	apiKey := os.Getenv("LASTFM_API_KEY")
	if apiKey == "" {
		fmt.Println("Must specify LASTFM_API_KEY")
		os.Exit(1)
	}

	lfmUser := os.Getenv("LASTFM_USER")
	if lfmUser == "" {
		fmt.Println("Must specify LASTFM_USER")
		os.Exit(1)
	}

	addr := os.Getenv("LISTEN_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	lfm := lastfm.LFM{
		Client:  http.DefaultClient,
		BaseURL: "https://ws.audioscrobbler.com/2.0/",
		User:    lfmUser,
		APIKey:  apiKey,
	}

	tmpl := template.Must(template.New("index").Parse(indexTemplate))

	app := app{
		lfm:  &lfm,
		page: tmpl,
	}

	http.HandleFunc("/", app.handleBase)
	http.HandleFunc("/api", app.handleAPI)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (app *app) handleBase(w http.ResponseWriter, r *http.Request) {

	// TODO:
	// * remove template parsing on each request

	data := struct{}{}
	//tmpl := template.Must(template.ParseFiles("templates/albums.html"))

	err := app.page.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func (app *app) handleAPI(w http.ResponseWriter, r *http.Request) {

	// TODO:
	// * prevent requests being made too frequently
	// * if we are now playing, get more data about the track

	data := struct {
		NowPlaying  bool           `json:"now_playing"`
		RecentTrack *lastfm.Track  `json:"recent_track"`
		Albums      []lastfm.Album `json:"albums,omitempty"`
		User        *lastfm.User   `json:"user,omitempty"`
	}{}

	// Get recent tracks
	tracks, err := app.lfm.RecentTracks()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(tracks) > 0 {
		data.RecentTrack = &tracks[0]
	}

	if tracks[0].NowPlaying {
		data.NowPlaying = true

		respondJSON(w, r, data)
		return
	}

	albums, err := app.lfm.TopAlbums()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := app.lfm.UserInfo()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data.NowPlaying = false
	data.Albums = albums
	data.User = user
	respondJSON(w, r, data)
}

func respondJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
