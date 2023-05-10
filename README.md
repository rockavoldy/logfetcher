# Log fetcher

Fetch log from google drive and save it to local

## Usage
1. Download credential json from oAuth2 and save it to `credentials.json`
2. Run the program
    ```bash
    env LOGFETCHER_FOLDER_ID=folder_id_that_hold_all_log ./logfetcher
    ```
3. The program will ask you to open a link in browser, open it and login with your google account
4. After finished, you will be redirected to localhost blablabla, check the browser and copy the value from parameter `code` from url, then paste it to the same terminal with the link instruction
5. The program will create token.json if it's good, and the log will be fetched at the time it set (by default, every 03:00 UTC+7)

## Build

```bash
# install dependencies first
go mod tidy

# for Apple Silicon
env GOOS=darwin GOARCH=arm64 go build -o logfetcher
# for Apple intel
env GOOS=darwin GOARCH=amd64 go build -o logfetcher
# for Linux with intel/amd processor (x86_64)
env GOOS=linux GOARCH=amd64 go build -o logfetcher
# for Linux with arm64 processor (raspi 3, raspi 4, rockpiS, etc.)
env GOOS=linux GOARCH=arm64 go build -o logfetcher
# for Linux with arm processor (raspi 1, raspi 2, etc.)
env GOOS=linux GOARCH=arm go build -o logfetcher
```
