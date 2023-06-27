# How to run test

```bash
go test ./... -v
```

# Example Request

```json
{
    "winningNumbers": "3 9 28 32 37 46",
    "additionalNumber": "7",
    "bets": [
       "8 14 19 22 26 31",
       "4 12 19 32 35 40",
       "6 7 23 28 33 46"
    ]
}
```
# How to build for deployment

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
```
