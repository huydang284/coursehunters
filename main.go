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

    config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
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
    playlistId := createPlaylist(service, title)

    for {
        lessonIndex := getNextLessonIndex()

        if lessonIndex == -1 {
            break
        }
        lesson := courses[lessonIndex]
        // download video from coursehunters
        downloadedFile := DownloadFile(lesson.url, "./temp/")
        // upload video to youtube
        videoId := uploadVideo(service, "./temp/"+downloadedFile, lesson.title)
        // add to playlist
        addVideoToPlaylist(service, playlistId, videoId)
        courses[lessonIndex].downloaded = true
    }
    // empty temp folders
    emptyTempFolder()
}
