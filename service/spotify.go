package service

import (
	"encoding/json"
	"golang-songs/model"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//func GetMusicArtistId(token string) (*model.Response, error) {
//
//	values := url.Values{}
//	//log.Println("values", values)
//	//values.Add("q", w.Weather[0].Main)
//
//	//とりあえず手動で
//	//アーティストのuuid?
//	//values.Add("q", "1O8CSXsPwEqxcoBE360PPO")
//
//	//values.Add("type", "artist")
//	log.Println("values:", values)
//
//	//req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/search", nil)
//	//req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/artists", nil)
//	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/artists/1O8CSXsPwEqxcoBE360PPO", nil)
//	req.URL.RawQuery = values.Encode()
//	req.Header.Set("Authorization", "Bearer "+token)
//
//	client := &http.Client{}
//	log.Println("req:", req)
//	log.Println("client:", client)
//
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Println("client.Do失敗")
//		return nil, err
//	}
//	log.Println("resp.Body:", resp.Body)
//
//	defer resp.Body.Close()
//
//	b, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Println("err:", err)
//		return nil, err
//	}
//	log.Println("b:", b)
//
//	var artist model.Artist
//
//	json.Unmarshal(b, &artist)
//
//	//log.Println("b:", b)
//
//	response := &model.Response{
//		//CityName: w.Name,
//		Artist: artist,
//		//Weather:  w.Weather[0].Main,
//	}
//
//	log.Println("response:", response)
//
//	return response, nil
//}

//func SearchMusicArtists(token string) (*model.Response, error) {
//
//	values := url.Values{}
//
//	values.Add("type", "artist")
//	values.Add("q", "Radiohead")
//
//	log.Println("values:", values)
//
//	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/search", nil)
//	req.URL.RawQuery = values.Encode()
//	req.Header.Set("Authorization", "Bearer "+token)
//
//	client := &http.Client{}
//	log.Println("req:", req)
//	//log.Println("client:", client)
//
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Println("client.Do失敗")
//		return nil, err
//	}
//	log.Println("resp:", resp)
//
//	defer resp.Body.Close()
//
//	b, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Println("err:", err)
//		return nil, err
//	}
//	//log.Println("b:", b)
//
//	var artists model.Artists
//	log.Print("model.Artists:", artists)
//
//	json.Unmarshal(b, &artists)
//
//	response := &model.Response{
//		Artists: artists,
//	}
//
//	log.Println("response:", response)
//
//	return response, nil
//}

func GetTracks(token string, title string) (*model.Response, error) {

	values := url.Values{}

	values.Add("type", "track")
	values.Add("q", title)
	values.Add("market", "JP")
	values.Add("limit", "10")

	log.Println("values:", values)

	log.Println("title:", title)

	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/search", nil)
	req.URL.RawQuery = values.Encode()
	//ヘッダに」アクセストークン入れている
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	log.Println("req:", req)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do失敗")
		return nil, err
	}
	log.Println("resp:", resp)

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	//log.Println("b:", b)

	var tracks model.Tracks
	log.Print("model.Tracks:", tracks)

	json.Unmarshal(b, &tracks)

	response := &model.Response{
		Tracks: tracks,
	}

	log.Println("response:", response)

	return response, nil
}
