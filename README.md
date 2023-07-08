# Log fetcher

Fetch log from google drive and save it to local

## Usage
1. Download credential json from oAuth2 and save it to `credentials.json`
2. Run the program
    ```bash
    env LOGFETCHER_FOLDER_ID=<gdrive_folder_id> ./logfetcher
    ```
3. The program will ask you to open a link in browser, open it and login with your google account
4. After finished, you will be redirect to localhost blablabla, check the url and copy the value from parameter `code` from url, then paste it to the same terminal with the instruction
5. The program will create token.json if it's correct, and the log will be fetched at the time it set (by default, every 03:00 UTC+7)

## Build

```bash
# install dependencies first
go mod tidy

# for Apple Silicon
env GOOS=darwin GOARCH=arm64 go build -o logfetcher-darwin-arm64
# for Apple intel
env GOOS=darwin GOARCH=amd64 go build -o logfetcher-darmin-amd64
# for Linux with intel/amd processor (x86_64)
env GOOS=linux GOARCH=amd64 go build -o logfetcher-linux-amd64
# for Linux with arm64 processor (raspi 3, raspi 4, rockpiS, etc.)
env GOOS=linux GOARCH=arm64 go build -o logfetcher-linux-arm64
# for Linux with arm processor (raspi 1, raspi 2, etc.)
env GOOS=linux GOARCH=arm go build -o logfetcher-linux-arm
```

> Something to note, as for now, it takes the filename as a date to check if it's older enough since it's the case at the time. You can check how it's work in the function `func (app App) DeleteOldLog(interval IntervalDays)`.
