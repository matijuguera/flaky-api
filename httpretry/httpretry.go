package httpretry

import (
	"errors"
	apierror "flaky-api/apierror"
	"fmt"
	"net/http"
	"time"
)

func Get(URL string, retryDuration time.Duration, retries int) (*http.Response, error) {
	errs := errors.New("max amount of retries reached, no response was found")

	for i := 0; i <= retries; i++ {

		resp, err := http.Get(URL)
		if err != nil {
			errs = fmt.Errorf("%w \n -%v", errs, apierror.NewAPIError(resp.StatusCode, "http retry get error", URL, resp.Status))
			continue
		}

		if resp != nil && resp.StatusCode == 200 {
			return resp, nil
		} else {
			errs = fmt.Errorf("%w \n -%v", errs, apierror.NewAPIError(resp.StatusCode, "http retry status error", URL, resp.Status))
		}

		//avoids sleep at the last iteration
		if i != retries {
			time.Sleep(retryDuration)
		}
	}

	return nil, errs
}
