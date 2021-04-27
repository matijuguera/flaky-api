package httpretry

import (
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	exponentialValues := []int{1, 2, 4, 8, 16, 32}

	for index, value := range exponentialValues {
		backoffValue := ExponentialBackoff(index)
		exponentialValue := time.Duration(value) * time.Second
		if backoffValue == exponentialValue {
			t.Errorf("Got %d; want %d", backoffValue, exponentialValue)
		}
	}

}
