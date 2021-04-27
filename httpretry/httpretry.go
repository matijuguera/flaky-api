package httpretry

import (
	"errors"
	apierror "flaky-api/apierror"
	"fmt"
	"net/http"
	"time"
)

type BackoffStrategy func(retry int) time.Duration

func ExponentialBackoff(i int) time.Duration {
	return time.Duration(1<<uint(i)) * time.Second
}

func LinearBackoff(i int) time.Duration {
	return time.Duration(i) * time.Second
}

func DefaultBackoff(_ int) time.Duration {
	return 1 * time.Second
}

func Get(URL string, backoffStrategy BackoffStrategy, retries int) (*http.Response, error) {
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
			time.Sleep(backoffStrategy(i))
		}
	}

	return nil, errs
}
