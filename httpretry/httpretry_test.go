package httpretry

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ClientMock struct {
	http.Client
	response *http.Response
	err      error
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return c.response, c.err
}

func TestExponentialBackoff(t *testing.T) {
	exponentialValues := []int{1, 2, 4, 8, 16, 32}

	for index, value := range exponentialValues {
		backoffValue := ExponentialBackoff(index)
		exponentialValue := time.Duration(value) * time.Second
		assert.Equal(t, exponentialValue, backoffValue, fmt.Sprintf("Got %d; want %d", backoffValue, exponentialValue))
	}

}

func TestLinearBackoff(t *testing.T) {
	linearValues := []int{0, 1, 2, 3, 4}

	for index, value := range linearValues {
		backoffValue := LinearBackoff(index)
		linearValue := time.Duration(value) * time.Second
		assert.Equal(t, linearValue, backoffValue, fmt.Sprintf("Got %d; want %d", backoffValue, linearValue))
	}
}

func TestDefaultBackoff(t *testing.T) {
	backoffValue := DefaultBackoff(0)
	defaultValue := 1 * time.Second
	assert.Equal(t, defaultValue, backoffValue, fmt.Sprintf("Got %d; want %d", backoffValue, defaultValue))
}

func TestSuccessResponse(t *testing.T) {
	//build http mock
	newClientMock := &ClientMock{}
	json := `{"id":0}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	newClientMock.response = &http.Response{
		StatusCode: 200,
		Body:       r,
	}

	//do fake request
	httpRetryClient := New(newClientMock)
	res, _ := httpRetryClient.Get("mock")
	assert.Equal(t, 200, res.StatusCode, fmt.Sprintf("got %d; want 200", res.StatusCode))
}

func TestWrongResponse(t *testing.T) {
	//build http mock
	newClientMock := &ClientMock{}
	newClientMock.err = errors.New("Unexpected error!")

	//do fake request
	httpRetryClient := New(newClientMock)
	_, err := httpRetryClient.Get("mock")

	assert.Error(t, err, "expects an error")
}

func TestMaxRetriesMustBeGreaterThanZero(t *testing.T) {
	//build http retry client
	newClientMock := &ClientMock{}
	httpRetryClient := New(newClientMock)
	httpRetryClient.MaxRetries = 0

	_, err := httpRetryClient.Get("mockurl")

	assert.Equal(t, ErrMaxRetriesGreaterThanZero, err, fmt.Sprintf("got %v; want ErrMaxRetriesGreaterThanZero", err.Error()))
}
