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

func getAlbums(w http.ResponseWriter, r *http.Request) {
	token := "BQAP_26A1KQI1pYQE6PfZ4C7K_HS2mw26vJac8KR6yljSta2-IANOvRva5a388Ijysm2ceHXj841PTDm1dMrUP5rfRp4krGSuYbHbwwNydWf-xeRHT4"

	albumEndpoint := "/artists/3IW7ScrzXmPvZhB27hmfgy/albums"
	albumBody, err := appelAPI(token, albumEndpoint)
	if err != nil {
		fmt.Println("Erreur lors de la requête des albums:", err)
		http.Error(w, "Erreur lors de la requête des albums", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Response JSON: %s\n", albumBody)

	var response ResponseAlbums
	err = json.Unmarshal(albumBody, &response)
	if err != nil {
		fmt.Println("Erreur lors de la désérialisation JSON des albums:", err)
		http.Error(w, "Erreur lors de la désérialisation JSON des albums", http.StatusInternalServerError)
		return
	}

	for _, album := range response.Items {
		fmt.Fprintf(w, "Nom de l'album: %s\n", album.Name)
		fmt.Fprintf(w, "Image de couverture: %s\n", album.Images[0].URL)
		fmt.Fprintf(w, "Date de sortie: %s\n", album.ReleaseDate)
		fmt.Fprintf(w, "Nombre de titre: %d\n", album.Total_tracks)
		fmt.Fprintln(w, "------------------------")
	}
}

func getTrack(w http.ResponseWriter, r *http.Request) {
	token := "BQAP_26A1KQI1pYQE6PfZ4C7K_HS2mw26vJac8KR6yljSta2-IANOvRva5a388Ijysm2ceHXj841PTDm1dMrUP5rfRp4krGSuYbHbwwNydWf-xeRHT4"

	trackEndpoint := "/tracks/0EzNyXyU7gHzj2TN8qYThj"
	trackBody, err := appelAPI(token, trackEndpoint)
	if err != nil {
		fmt.Println("Erreur lors de la requête de la piste:", err)
		http.Error(w, "Erreur lors de la requête de la piste", http.StatusInternalServerError)
		return
	}

	var trackInfo Track
	err = json.Unmarshal(trackBody, &trackInfo)
	if err != nil {
		fmt.Println("Erreur lors de la désérialisation JSON de la piste:", err)
		http.Error(w, "Erreur lors de la désérialisation JSON de la piste", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Nom de la piste: %s\n", trackInfo.Name)
	fmt.Fprintf(w, "Nom de l'album: %s\n", trackInfo.Album.Name)
	fmt.Fprintf(w, "Artiste: %s\n", trackInfo.Artists[0].Name)
	fmt.Fprintf(w, "Date de sortie: %s\n", trackInfo.ReleaseDate)
	fmt.Fprintf(w, "URL Spotify: %s\n", trackInfo.ExternalURLs.Spotify)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Erreur lors de l'analyse du modèle HTML", http.StatusInternalServerError)
		return
	}

	

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Erreur lors du rendu du modèle HTML", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/album/jul", getAlbums)
	http.HandleFunc("/track/sdm", getTrack)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":8080", nil)
}
