# Random tasks from interviews

## `cmd/requests`

### Goal
Need to create a CLI utility, which will run parallel/concurrent HTTP GET requests, slowly increasing amount of them
and printing how many of them are still running on the moment of first Timeout.

### Usage
```
-m int
    int value by witch we will multiply amount of requests made on previous iteration (default 2)
-p int
    how often in milliseconds we need to increase amount of HTTP requests (default 1000)
-t int
    timeout in milliseconds from net.Dialer till end of response from remote end (default 100)
-url string
    url to which GET requests will be sent (default "https://ya.ru")
```

Examples: 
- `go run cmd/requests/main.go` - with default parameters
- `go run cmd/requests/main.go -t 600 -m 2 -p 500 -url https://ya.ru`

## `cmd/policedpts`

### Goal
During 15 min programming part of the interview we need to create a CLI utility, using these docs https://data.police.uk/docs/,
which will output each police force with its phone number in csv format, if it has a facebook account. Fields that
we need in output: Force name, phone number, Facebook URL.

### Usage

Examples:
- `go run cmd/policedpts/main.go`
