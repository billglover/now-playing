package lastfm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Album struct {
	CoverImageURL string
	Artist        string
	Title         string
	PlayCount     int
}

type Track struct {
	Name       string
	Artist     string
	Album      string
	ImgURL     string
	NowPlaying bool
}

type User struct {
	PlayCount int
	Name      string
}

type LFM struct {
	Client  *http.Client
	BaseURL string
	User    string
	APIKey  string
}

func (lfm *LFM) UserInfo() (*User, error) {

	type Response struct {
		User struct {
			Country    string `json:"country"`
			Age        string `json:"age"`
			PlayCount  string `json:"playcount"`
			Subscriber string `json:"subscriber"`
			RealName   string `json:"realname"`
			Playlists  string `json:"playlists"`
			Bootstrap  string `json:"bootstrap"`
			Image      []struct {
				Size string `json:"size"`
				Text string `json:"#text"`
			} `json:"image"`
			Registered struct {
				Unixtime string `json:"unixtime"`
				Text     int    `json:"#text"`
			} `json:"registered"`
			URL    string `json:"url"`
			Gender string `json:"gender"`
			Name   string `json:"name"`
			Type   string `json:"type"`
		} `json:"user"`
	}

	req, _ := http.NewRequest(http.MethodGet, lfm.BaseURL, nil)
	q := req.URL.Query()
	q.Set("method", "user.getinfo")
	q.Set("api_key", lfm.APIKey)
	q.Set("user", lfm.User)
	q.Set("format", "json")
	req.URL.RawQuery = q.Encode()

	resp, err := lfm.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := Response{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	playCount, err := strconv.Atoi(data.User.PlayCount)
	if err != nil {
		fmt.Println(err)
		playCount = 0
	}

	user := User{
		PlayCount: playCount,
		Name:      data.User.RealName,
	}

	return &user, nil
}

func (lfm *LFM) RecentTracks() ([]Track, error) {

	type Response struct {
		Recenttracks struct {
			Track []struct {
				Artist struct {
					Mbid string `json:"mbid"`
					Text string `json:"#text"`
				} `json:"artist"`
				Streamable string `json:"streamable"`
				Image      []struct {
					Size string `json:"size"`
					Text string `json:"#text"`
				} `json:"image"`
				Mbid  string `json:"mbid"`
				Album struct {
					Mbid string `json:"mbid"`
					Text string `json:"#text"`
				} `json:"album"`
				Name string `json:"name"`
				Attr struct {
					Nowplaying string `json:"nowplaying"`
				} `json:"@attr,omitempty"`
				URL  string `json:"url"`
				Date struct {
					Uts  string `json:"uts"`
					Text string `json:"#text"`
				} `json:"date,omitempty"`
			} `json:"track"`
			Attr struct {
				User       string `json:"user"`
				TotalPages string `json:"totalPages"`
				Page       string `json:"page"`
				Total      string `json:"total"`
				PerPage    string `json:"perPage"`
			} `json:"@attr"`
		} `json:"recenttracks"`
	}

	req, _ := http.NewRequest(http.MethodGet, lfm.BaseURL, nil)
	q := req.URL.Query()
	q.Set("method", "user.getrecenttracks")
	q.Set("api_key", lfm.APIKey)
	q.Set("user", lfm.User)
	q.Set("format", "json")
	q.Set("limit", "2")
	req.URL.RawQuery = q.Encode()

	resp, err := lfm.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := Response{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	tracks := []Track{}

	for _, t := range data.Recenttracks.Track {
		curTrack := Track{
			Name:   t.Name,
			Artist: t.Artist.Text,
			Album:  t.Album.Text,
			ImgURL: t.Image[3].Text,
		}
		if t.Attr.Nowplaying == "true" {
			curTrack.NowPlaying = true
		}
		tracks = append(tracks, curTrack)
	}

	return tracks, nil
}

func (lfm *LFM) TopAlbums() ([]Album, error) {

	type Artist struct {
		URL  string `json:"url"`
		Name string `json:"name"`
		Mbid string `json:"mbid"`
	}
	type Image struct {
		Size string `json:"size"`
		Text string `json:"#text"`
	}
	type AlbumsAttr struct {
		Rank string `json:"rank"`
	}
	type AlbumAttr struct {
		User       string `json:"user"`
		TotalPages string `json:"totalPages"`
		Page       string `json:"page"`
		PerPage    string `json:"perPage"`
		Total      string `json:"total"`
	}
	type LfmAlbum struct {
		Artist    Artist    `json:"artist"`
		Image     []Image   `json:"image"`
		Mbid      string    `json:"mbid"`
		URL       string    `json:"url"`
		Playcount string    `json:"playcount"`
		Attr      AlbumAttr `json:"@attr"`
		Name      string    `json:"name"`
	}
	type Topalbums struct {
		Albums []LfmAlbum `json:"album"`
		Attr   AlbumsAttr `json:"@attr"`
	}
	type Response struct {
		Topalbums Topalbums `json:"topalbums"`
	}

	req, _ := http.NewRequest(http.MethodGet, lfm.BaseURL, nil)
	q := req.URL.Query()
	q.Set("method", "user.gettopalbums")
	q.Set("api_key", lfm.APIKey)
	q.Set("user", lfm.User)
	q.Set("format", "json")
	q.Set("limit", "12")
	q.Set("period", "7day")
	req.URL.RawQuery = q.Encode()

	resp, err := lfm.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := Response{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	albums := []Album{}

	for _, a := range data.Topalbums.Albums {

		playCount, err := strconv.Atoi(a.Playcount)
		if err != nil {
			fmt.Println(err)
			playCount = 0
		}

		albums = append(albums, Album{
			CoverImageURL: a.Image[3].Text,
			Artist:        a.Artist.Name,
			Title:         a.Name,
			PlayCount:     playCount,
		})
	}

	return albums, nil
}
