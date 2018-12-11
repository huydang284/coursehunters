package main

import (
    "github.com/PuerkitoBio/goquery"
    "log"
    "net/http"
)

var title string

type Lesson struct {
    title      string
    url        string
}

func getVideos(courseUrl string) (courses []Lesson) {
    log.Println("Getting lessons url...")
    // todo validate url

    // parse from html
    resp, err := http.Get(courseUrl)
    if err != nil {
        check(err)
    }

    html, err := goquery.NewDocumentFromReader(resp.Body)
    defer resp.Body.Close()
    check(err)

    title = html.Find("article h2").Text()
    html.Find("li.lessons-list__li").Each(func(i int, selection *goquery.Selection) {
        lessonTitle := selection.Find("[itemprop=\"name\"]").Text()
        url, _ := selection.Find("[itemprop=\"url\"]").Attr("href")

        courses = append(courses, Lesson{
            lessonTitle,
            url,
        })
    })
    return
}
