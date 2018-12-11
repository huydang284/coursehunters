package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/user"
    "path/filepath"
    "strings"
)

func main() {
    var dir string
    var isUpload bool
    var playlist string
    var help bool
    flag.StringVar(&dir, "dir", "./", "Execution directory")
    flag.BoolVar(&isUpload, "upload", false, "Upload directory to Youtube")
    flag.StringVar(&playlist, "playlist", "Playlist", "Youtube playlist")
    flag.BoolVar(&help, "help", false, "Command help")
    flag.Parse()

    if isUpload {
        uploadDirectoryToYoutube(dir, playlist)
        return
    }

    if help {
        printHelp()
        return
    }

    defaultProcess(playlist)
}

func defaultProcess(playlist string) {
    fmt.Print("Please enter coursehunters course url: ")
    var courseUrl string
    var downloaded []string
    if _, err := fmt.Scan(&courseUrl); err != nil {
        log.Fatalln("Cannot get course url")
    }

    service := initService()
    lessons := getVideos(courseUrl)

    if playlist == "Playlist" {
        playlist = title
    }
    playlistId := getPlaylist(service, playlist)
    log.Println("Creating temp directory")
    tempDir := "./temp"
    err := os.Mkdir(tempDir, 0777)
    check(err)

    jsonFile := tempDir + "/downloaded.json"
    data, err := ioutil.ReadFile(jsonFile)

    if err == nil {
        err = json.Unmarshal(data, &downloaded)
        check(err)
    }

    for _, lesson := range lessons {
        if inSlice(lesson.url, downloaded) {
            continue
        }
        // download video from coursehunters
        downloadedFile := downloadFile(lesson.url, tempDir)
        // upload video to youtube
        videoId := uploadVideo(service, tempDir+"/"+downloadedFile, lesson.title)
        // add to playlist
        addVideoToPlaylist(service, playlistId, videoId)
        // append downloaded
        downloaded = append(downloaded, lesson.url)
        downloadedJson, err := json.Marshal(downloaded)
        check(err)
        err = ioutil.WriteFile(jsonFile, downloadedJson, 0777)
    }

    check(os.Remove(tempDir))
    log.Println("Done")
}

func uploadDirectoryToYoutube(directory string, playlistName string) {
    var err error
    var uploaded []string
    usr, _ := user.Current()
    homeDir := usr.HomeDir

    if strings.HasPrefix(directory, "~/") {
        directory = filepath.Join(homeDir, directory[2:])
    }

    if _, err := os.Stat(directory); os.IsNotExist(err) {
        log.Fatalf("Directory %s is not exists", directory)
    }

    service := initService()
    // create playlist
    log.Printf("Get playlist \"%s\"\n", playlistName)
    playlistId := getPlaylist(service, playlistName)

    // log file
    jsonFile := strings.TrimRight(directory, "/") + "/uploaded.json"
    data, err := ioutil.ReadFile(jsonFile)

    if err == nil {
        err = json.Unmarshal(data, &uploaded)
        check(err)
    }

    log.Println("Uploading videos")
    err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        log.Println("Scanned file: " + path)
        if filepath.Ext(path) != ".mp4" {
            log.Println("> Invalid file type")
            return nil
        }
        if inSlice(path, uploaded) {
            log.Println("> This file was uploaded")
            return nil
        }
        // upload videos to Youtube
        videoTitle := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
        videoId := uploadVideo(service, path, videoTitle)
        // add video to playlist
        addVideoToPlaylist(service, playlistId, videoId)
        // save uploaded files for resuming
        uploaded = append(uploaded, path)
        uploadedJson, err := json.Marshal(uploaded)
        check(err)
        err = ioutil.WriteFile(jsonFile, uploadedJson, 0777)
        check(err)
        return nil
    })

    check(err)
    // delete uploaded files
    check(os.Remove(jsonFile))
    log.Println("Done")
}

func printHelp() {
    fmt.Println("Coursehunters tool by Huy Dang\n\n" +
        "Usage:\n" +
        "> Download videos from coursehunters then uploading to Youtube:\n" +
        "\t ./coursehunters [-playlist=\"Your playlist name\"]\n\n" +
        "> Upload videos from directory:\n" +
        "\t ./coursehunters -upload -dir=~/sample [-playlist=\"Your playlist name\"]")
}
