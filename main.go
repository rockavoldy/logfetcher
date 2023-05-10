package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"logfetcher/helper"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

type App struct {
	Drive *drive.Service
	ctx   context.Context
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	folderID := os.Getenv("LOGFETCHER_FOLDER_ID")
	if folderID == "" {
		log.Fatalln("LOGFETCHER_FOLDER_ID is not set")
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	app := App{
		Drive: srv,
		ctx:   ctx,
	}
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	s := gocron.NewScheduler(time.UTC)
	// run it at 03:00 UTC+7 everyday, as the log is uploded at 01:30 UTC+7
	s.Every(1).Day().At("20:00:00").Do(func() {
		app.listDir(folderID)
		log.Println("Fetch new log")
	})

	s.StartImmediately().StartBlocking()
}

func (app App) listDir(FolderID string) error {
	r, err := app.Drive.Files.List().
		Corpora("user").
		IncludeItemsFromAllDrives(true).
		OrderBy("createdTime desc").
		Q(fmt.Sprintf("'%s' in parents", FolderID)).
		PageSize(1).
		SupportsAllDrives(true).
		Fields("nextPageToken, files(id, name)").Do()

	if err != nil {
		log.Fatalln(helper.HandleErr(helper.ErrNotFound, err.Error()))
	}
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			app.downloadDir(i.Id)
		}
	}

	return nil
}

func (app App) downloadDir(FolderID string) error {
	folder, err := app.Drive.Files.Get(FolderID).Do()
	if err != nil {
		return helper.HandleErr(helper.ErrNotFound, err.Error())
	}

	folderName := folder.Name

	err = os.Mkdir(folderName, 0755)
	if err != nil {
		if err == os.ErrNotExist {
			os.Mkdir(folderName, 0755)
		} else {
			log.Println("Folder already exist")
		}
	}

	for {
		res, err := app.Drive.Files.List().
			Corpora("user").
			IncludeItemsFromAllDrives(true).
			OrderBy("createdTime desc").
			Q(fmt.Sprintf("'%s' in parents", FolderID)).
			PageSize(1000).
			SupportsAllDrives(true).
			Do()

		if err != nil {
			return helper.HandleErr(helper.ErrNotFound, err.Error())
		}

		items := res.Files
		for _, item := range items {
			log.Printf("%s: %s", item.Id, item.Name)
		}
		for _, item := range items {
			if item.MimeType == "application/vnd.google-apps.folder" {
				continue
			}
			app.downloadFile(item.Id, item.Name, folderName)
			log.Println(item.Name)
		}

		break
	}
	return nil
}

func (app App) downloadFile(fileID, fileName, folderName string) error {
	res, err := app.Drive.Files.Get(fileID).Download()
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	outFile, err := os.Create(fmt.Sprintf("%s/%s", folderName, fileName))
	if err != nil {
		log.Println(err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		log.Println(err)
	}

	return nil
}
