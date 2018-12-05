package main

import (
    "fmt"
    "github.com/PuerkitoBio/goquery"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

var courses []Lesson
var title string

type Lesson struct {
    title      string
    url        string
    downloaded bool
}

func GetVideos(courseUrl string) {
    fmt.Println("Getting videos info...")
    // todo validate url

    // parse from html
    resp, err := http.Get(courseUrl)
    if err != nil {
        panic(err)
    }

    html, err := goquery.NewDocumentFromReader(resp.Body)
    defer resp.Body.Close()

    if err != nil {
        panic(err)
    }

    //defer resp.Body.Close()

    title = html.Find("article h2").Text()
    html.Find("li.lessons-list__li").Each(func(i int, selection *goquery.Selection) {
        lessonTitle := selection.Find("[itemprop=\"name\"]").Text()
        url, _ := selection.Find("[itemprop=\"url\"]").Attr("href")

        courses = append(courses, Lesson{
            lessonTitle,
            url,
            false,
        })
    })
}

func getNextLessonIndex() int {
    for index := range courses {
        if courses[index].downloaded {
            continue
        }
        return index
    }
    return -1
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
