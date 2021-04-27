package house

import (
	"encoding/json"
	"flaky-api/downloader"
	"flaky-api/httpretry"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

type House struct {
	Id        int    `json:"id"`
	Address   string `json:"address"`
	Homeowner string `json:"homeowner"`
	Price     int    `json:"price"`
	PhotoURL  string `json:"photoURL"`
}

type HousesResponse struct {
	Houses []House `json:"houses"`
	Ok     bool    `json:"ok"`
}

const (
	HOMEVISION_ENDPOINT    = "http://app-homevision-staging.herokuapp.com/api_project"
	PHOTOS_REPOSITORY_PATH = "photos-repository/"
	DEFAULT_RETRY_DURATION = 2 * time.Second
	DEFAULT_MAX_RETRIES    = 10
)

func Get(page int) ([]House, error) {
	//do request
	res, err := httpretry.Get(fmt.Sprintf("%s/houses?page=%d", HOMEVISION_ENDPOINT, page), DEFAULT_RETRY_DURATION, DEFAULT_MAX_RETRIES)
	if err != nil {
		return []House{}, err
	}
	defer res.Body.Close()

	//to json
	housesResponse := HousesResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []House{}, err
	}
	if err := json.Unmarshal(body, &housesResponse); err != nil {
		return []House{}, err
	}

	return housesResponse.Houses, nil
}

func Download(h House, wgFile *sync.WaitGroup) {
	defer wgFile.Done()
	if err := downloader.DownloadFile(h.PhotoURL, PHOTOS_REPOSITORY_PATH, fmt.Sprintf("id-%d-%s.jpg", h.Id, h.Address)); err != nil {
		log.Printf("error downloading house: %v", err)
	}
}
