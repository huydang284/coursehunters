package main

import (
    "fmt"
    "golang.org/x/net/context"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/youtube/v3"
    "io/ioutil"
    "log"
)

func main() {
    ctx := context.Background()

    b, err := ioutil.ReadFile("client_secret.json")
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

    config, err := google.ConfigFromJSON(b, youtube.YoutubeUploadScope)
    if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
    }

    client := getClient(ctx, config)
    service, err := youtube.New(client)

    handleError(err, "Error creating YouTube client")

    fmt.Print("Please enter coursehunters course url: ")
    var courseUrl string
    if _, err := fmt.Scan(&courseUrl); err != nil {
        log.Fatal("Cannot get course url")
    }

    GetVideos(courseUrl)
    for {
        url := GetNextVideo()
        if url == "" {
            break
        }
        downloadedFile := DownloadFile(url, "./temp/")
        uploadVideo(service, "./temp/"+downloadedFile)
        courses[url] = false // false as DONE
    }
    // empty temp folders
    emptyTempFolder()
}
