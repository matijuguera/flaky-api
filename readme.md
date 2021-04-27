# Flaky API: Dealing with errors when fetching data

This is a take-home interview for HomeVision that focuses primarily on writing clean code that accomplishes a very practical task.

We have a simple pagianted API hosted at [http://app-homevision-staging.herokuapp.com/api_project/houses?page=1](http://app-homevision-staging.herokuapp.com/api_project/houses?page=1) that returns a list of houses along with some metadata. The task will be to write a script that accomplishes the following tasks:

1. Requests the first 10 pages of results from the API
2. Parses the JSON returned by the API
3. Downloads the photo for each house and saves it in a file with the format: `id-[NNN]-[address].[ext]`

There are a few gotchas to watch out for:

1. This is a _flaky_ API! That means that it will likely fail with a non-200 response code. Your code _must_ handle these errors correctly so that all photos are downloaded
2. Downloading photos is slow so please think a bit about how you would optimize your downloads, making use of concurrency

# `go run <project folder>`

Use `go run <project folder>` to run the main.go
