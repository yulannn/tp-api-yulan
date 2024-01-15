package main

import (
	"html/template"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)


type Album struct {
	Name        string    `json:"name"`
	Images      []Image   `json:"images"`
	ReleaseDate string    `json:"release_date"`
	Total_tracks int 	  `json:"total_tracks"`
}

type ResponseAlbums struct {
	Items []Album `json:"items"`
}


type Track struct {
	Name         string     `json:"name"`
	Album        Album      `json:"album"`
	Artists      []Artist   `json:"artists"`
	ReleaseDate  string     `json:"release_date"`
	ExternalURLs ExternalURL `json:"external_urls"`
}


type Image struct {
	URL string `json:"url"`
}


type Artist struct {
	Name string `json:"name"`
}


type ExternalURL struct {
	Spotify string `json:"spotify"`
}



func appelAPI(token string, endpoint string) ([]byte, error) {
	urlAPI := "https://api.spotify.com/v1" + endpoint

	req, err := http.NewRequest("GET", urlAPI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getAlbums(w http.ResponseWriter, r *http.Request) ResponseAlbums{
	token := "BQAbDrgEFPF7ZEmJtwkMOkn8Y5Jo-JHssOrzQFJCskLPT71-LdYNCKjY6jTdzn1kt4ktS2L1RZGhfz6l2pCkA7XBkyDQ-maazxJ2RBwMVXiHwyrYBSw"

	albumEndpoint := "/artists/3IW7ScrzXmPvZhB27hmfgy/albums"
	albumBody, _ := appelAPI(token, albumEndpoint)

	var response ResponseAlbums
	json.Unmarshal(albumBody, &response)

	return response
}

func getTrack(w http.ResponseWriter, r *http.Request) Track {
	token := "BQAbDrgEFPF7ZEmJtwkMOkn8Y5Jo-JHssOrzQFJCskLPT71-LdYNCKjY6jTdzn1kt4ktS2L1RZGhfz6l2pCkA7XBkyDQ-maazxJ2RBwMVXiHwyrYBSw"

	trackEndpoint := "/tracks/0EzNyXyU7gHzj2TN8qYThj"
	trackBody, _ := appelAPI(token, trackEndpoint)

	var trackInfo Track
	json.Unmarshal(trackBody, &trackInfo)

	return trackInfo
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseGlob("./templates/*.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        fmt.Println("Error parsing templates:", err)
        return
    }
	err = temp.ExecuteTemplate(w, "index", nil)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        fmt.Println("Error executing template:", err)
        return
    }
}

func main() {
	temp, _ := template.ParseGlob("templates/*.html")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", homeHandler)

	http.HandleFunc("/album/jul", func(w http.ResponseWriter, r *http.Request){
		albums := getAlbums(w, r)
		temp.ExecuteTemplate(w, "album", albums)
	})

	http.HandleFunc("/track/sdm", func(w http.ResponseWriter, r *http.Request){
		track := getTrack(w, r)
		fmt.Println(track)
		temp.ExecuteTemplate(w, "track", track)
	})

	http.ListenAndServe(":8080", nil)
}
