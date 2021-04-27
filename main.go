package main

import (
	"flaky-api/house"
	"log"
	"sync"
)

const (
	TOTAL_PAGES = 10
)

func downloadHouses(page int, wg *sync.WaitGroup) {
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
	for i := 1; i <= TOTAL_PAGES; i++ {
		wg.Add(1)
		go downloadHouses(i, &wg)
	}
	wg.Wait()
}
