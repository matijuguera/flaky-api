package house

import (
	"encoding/json"
	"flaky-api/downloader"
	"flaky-api/httpretry"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	homevisionEndpoint   = "http://app-homevision-staging.herokuapp.com/api_project/"
	PhotosRepositoryPath = "photos-repository/"
)

var HttpRetryClient = newHttpRetryClient()

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

func newHttpRetryClient() *httpretry.Client {
	client := &http.Client{}
	client.Timeout = 30 * time.Second

	httpRetryClient := httpretry.New(client)
	httpRetryClient.Backoff = httpretry.LinearBackoff
	httpRetryClient.MaxRetries = 10

	return httpRetryClient
}

func Get(page int) ([]House, error) {
	//do request
	res, err := HttpRetryClient.Get(fmt.Sprintf("%shouses?page=%d", homevisionEndpoint, page))
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

func ConcurrentDownload(h House, wgFile *sync.WaitGroup) {
	defer wgFile.Done()
	if err := downloader.DownloadFile(h.PhotoURL, PhotosRepositoryPath, h.GetFilename()); err != nil {
		log.Printf("error downloading house [%d]: %v", h.Id, err)
	}
}

func (h *House) GetFilename() string {
	return fmt.Sprintf("id-%d-%s.jpg", h.Id, h.Address)
}
