package httpRetry

import (
	"errors"
	apierror "flaky-api/apierror"
	"log"
	"net/http"
	"time"
)

func Get(URL string, retryDuration time.Duration, retries int) (*http.Response, error) {
	for i := 0; i <= retries; i++ {

		resp, err := http.Get(URL)
		if err != nil {
			log.Println(apierror.NewAPIError(resp.StatusCode, "http retry get error", URL, resp.Status).Error())
			continue
		}

		if resp != nil && resp.StatusCode == 200 {
			return resp, nil
		} else {
			log.Println(apierror.NewAPIError(resp.StatusCode, "http retry status error", URL, resp.Status).Error())
		}

		//avoids sleep at the last iteration
		if i != retries {
			time.Sleep(retryDuration)
		}
	}

	return nil, errors.New("max amount of retries reached, no response was found")
}
