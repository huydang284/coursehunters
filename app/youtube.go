package main

import (
    "context"
    "fmt"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/youtube/v3"
    "io/ioutil"
    "log"
)

func Upload() {
    data, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatal(err)
    }
    conf, err := google.JWTConfigFromJSON(data)
    if err != nil {
        log.Fatal(err)
    }
    client := conf.Client(context.TODO())
    service, err := youtube.New(client)
    if err != nil {
        log.Fatal(err)
    }

    videosListById(service, "snippet,contentDetails,statistics", "Ks-_Mh1QhMc")
}

func printVideosListResults(response *youtube.VideoListResponse) {
    for _, item := range response.Items {
        fmt.Println(item.Id, ": ", item.Snippet.Title)
    }
}

func videosListById(service *youtube.Service, part string, id string) {
    call := service.Videos.List(part)
    if id != "" {
        call = call.Id(id)
    }
    response, err := call.Do()
    if err != nil {
        fmt.Println("3")
        log.Fatal(err)
    }
    fmt.Println("1")
    printVideosListResults(response)
    fmt.Println("2")
}
