package downloader

import (
	"flaky-api/apierror"
	"io"
	"net/http"
	"os"
)

func DownloadFile(URL, filePath string, fileName string) error {
	//Get the response bytes from the url
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return apierror.NewAPIError(resp.StatusCode, "downloading file error", URL, resp.Status)
	}

	//Create a empty file
	file, err := os.Create(filePath + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}
