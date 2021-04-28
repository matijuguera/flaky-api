package main

import (
	"flaky-api/house"
	"fmt"
	"log"
	"sync"
)

const (
	totalPages = 10
)

func concurrentDownloadHouses(page int, wg *sync.WaitGroup) {
	defer wg.Done()

	houses, err := house.Get(page)
	if err != nil {
		log.Printf("error getting houses page %d: %s \n", page, err)
		return
	}

	var wgFile sync.WaitGroup
	for _, h := range houses {
		wgFile.Add(1)
		go house.ConcurrentDownload(h, &wgFile)
	}
	wgFile.Wait()
}

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= totalPages; i++ {
		wg.Add(1)
		go concurrentDownloadHouses(i, &wg)
	}
	wg.Wait()
	fmt.Printf("Look inside the dir %s\n", house.PhotosRepositoryPath)
}
