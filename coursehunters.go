package main

import (
    "fmt"
    "golang.org/x/net/html"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

var courses map[string]bool

func init() {
    courses = make(map[string]bool)
}

func GetVideos(courseUrl string) {
    fmt.Println("Getting videos info")
    // todo validate url

    // parse from html
    resp, err := http.Get(courseUrl)
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    htmlTokens := html.NewTokenizer(resp.Body)
loop:
    for {
        tt := htmlTokens.Next()
        switch tt {
        case html.ErrorToken:
            break loop
        case html.SelfClosingTagToken:
            t := htmlTokens.Token()
            if t.Data == "link" {
                for _, attr := range t.Attr {
                    if attr.Key != "href" {
                        continue
                    }
                    // validate mp4 path
                    if !strings.HasSuffix(attr.Val, ".mp4") {
                        continue loop
                    }
                    found := courses[attr.Val]
                    if !found {
                        courses[attr.Val] = true
                    }
                }
            }
        }
    }

    //go saveToFile()
}

/*
func saveToFile() {
    json, err := json2.Marshal(courses)
    if err != nil {
        log.Fatal(err)
    }

    errW := ioutil.WriteFile("./temp/courses.json", json, 666)
    if errW != nil {
        log.Fatal(errW)
    }
}
*/
func GetNextVideo() string {
    for url, notDone := range courses {
        if !notDone {
            continue
        }
        return url
    }

    return ""
}

func emptyTempFolder() {
    files, err := filepath.Glob("./temp/*.mp4")
    if err != nil {
        log.Fatal(err)
    }
    for _, f := range files {
        if err := os.Remove(f); err != nil {
            log.Fatal(err)
        }
    }
}
