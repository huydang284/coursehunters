package main

import (
    "golang.org/x/net/html"
    "net/http"
    "strings"
)

func GetVideos(courseUrl string) map[string]bool {
    // todo validate url

    courses := make(map[string]bool)
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

    return courses
}

func DownloadVideos(courseUrl string) {
    courses := GetVideos(courseUrl)
    for url := range courses {
        DownloadFile(url, "./../temp/")
        break
    }
}
