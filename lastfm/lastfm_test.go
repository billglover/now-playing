package lastfm

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestTrackInfo(t *testing.T) {

}

func TestUserInfo(t *testing.T) {

}

func TestRecentTracks(t *testing.T) {
	f, err := os.Open("test/recentTracks.json")
	if err != nil {
		t.Error("unable to open test data:", err)
	}
	resp, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error("unable to open test data:", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			t.Errorf("queried unexpected URL, got: %s", req.URL.String())
		}

		q := req.URL.Query()

		if q.Get("method") != "user.getrecenttracks" {
			t.Errorf("didn't set the method correctly, got: %s", q.Get("method"))
		}

		if q.Get("format") != "json" {
			t.Errorf("didn't set the format correctly, got: %s", q.Get("format"))
		}

		l, err := strconv.Atoi(q.Get("limit"))
		if err != nil || l != 2 {
			t.Errorf("didn't set the limit correctly, got: %s", q.Get("limit"))
		}

		if q.Get("user") != "dummy_user" {
			t.Errorf("didn't set the user correctly, got: %s", q.Get("user"))
		}

		if q.Get("api_key") != "dummy_api_key" {
			t.Errorf("didn't set API key correctly, got: %s", q.Get("api_key"))
		}

		rw.Write(resp)
	}))
	defer server.Close()

	lfm := LFM{
		Client:  server.Client(),
		APIKey:  "dummy_api_key",
		User:    "dummy_user",
		BaseURL: server.URL,
	}
	tracks, err := lfm.RecentTracks()
	if err != nil {
		t.Error("unable to parse tracks:", err)
	}

	// We expect to get 3 items back (2 + 1 now playing)
	if len(tracks) != 3 {
		t.Errorf("got unexpected number of tracks: %d", len(tracks))
	}

	nowPlaying := 0
	for _, t := range tracks {
		if t.NowPlaying == true {
			nowPlaying++
		}
	}
	if nowPlaying != 1 {
		t.Errorf("unexpected number of NowPlaying tracks, got: %d", nowPlaying)
	}
}

func TestTopAlbums(t *testing.T) {
	f, err := os.Open("test/topAlbums.json")
	if err != nil {
		t.Error("unable to open test data:", err)
	}
	resp, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error("unable to open test data:", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			t.Errorf("queried unexpected URL, got: %s", req.URL.String())
		}

		q := req.URL.Query()

		if q.Get("method") != "user.gettopalbums" {
			t.Errorf("didn't set the method correctly, got: %s", q.Get("method"))
		}

		if q.Get("format") != "json" {
			t.Errorf("didn't set the format correctly, got: %s", q.Get("format"))
		}

		l, err := strconv.Atoi(q.Get("limit"))
		if err != nil || l != 12 {
			t.Errorf("didn't set the limit correctly, got: %s", q.Get("limit"))
		}

		if q.Get("period") != "7d" {
			t.Errorf("didn't set the period correctly, got: %s", q.Get("period"))
		}

		if q.Get("user") != "dummy_user" {
			t.Errorf("didn't set the user correctly, got: %s", q.Get("user"))
		}

		if q.Get("api_key") != "dummy_api_key" {
			t.Errorf("didn't set API key correctly, got: %s", q.Get("api_key"))
		}

		rw.Write(resp)
	}))
	defer server.Close()

	lfm := LFM{
		Client:  server.Client(),
		APIKey:  "dummy_api_key",
		User:    "dummy_user",
		BaseURL: server.URL,
	}
	albums, err := lfm.TopAlbums()
	if err != nil {
		t.Error("unable to parse albums:", err)
	}

	// We expect to get 12 items back
	if len(albums) != 12 {
		t.Errorf("got unexpected number of albums: %d", len(albums))
	}
}
