# Flaky API: Dealing with errors when fetching data

This is a take-home interview for HomeVision that focuses primarily on writing clean code that accomplishes a very practical task.

We have a simple pagianted API hosted at [http://app-homevision-staging.herokuapp.com/api_project/houses?page=1](http://app-homevision-staging.herokuapp.com/api_project/houses?page=1) that returns a list of houses along with some metadata. The task will be to write a script that accomplishes the following tasks:

1. Requests the first 10 pages of results from the API
2. Parses the JSON returned by the API
3. Downloads the photo for each house and saves it in a file with the format: `id-[NNN]-[address].[ext]`

There are a few gotchas to watch out for:

1. This is a _flaky_ API! That means that it will likely fail with a non-200 response code. Your code _must_ handle these errors correctly so that all photos are downloaded
2. Downloading photos is slow so please think a bit about how you would optimize your downloads, making use of concurrency

# Code organization

I use the code organization given in the GolangBootcamp book http://www.golangbootcamp.com/book/methods#cid33

# Http Retry Package

I created my own custom http-retry package, there you will find different Backoff Strategy

- ExponentialBackoff returns ever increasing backoffs by a power of 2

- LinearBackoff returns increasing durations, each a second longer than the last

- DefaultBackoff always returns 1 second

Example of use:

```
package house

func newHttpClient() *httpretry.Client {
	client := &http.Client{}
	httpRetryClient := httpretry.New(client)
	httpRetryClient.Backoff = httpretry.LinearBackoff <------

	return httpRetryClient
}
```

# Http Retry Package - Test - Testify

I tested the http retry package with the library https://github.com/stretchr/testify

To run the test `go test ./...`

# Run the project

Use `go run <project folder>` to run the main.go

# Photos

After running the script, you will see all the photos inside the folder "/photos-repository"
