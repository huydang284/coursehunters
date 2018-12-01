package main

/**
 * @website http://albulescu.ro
 * @author Cosmin Albulescu <cosmin@albulescu.ro>
 */

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path"
    "strconv"
    "time"
)

func PrintDownloadPercent(done chan int64, path string, total int64) {
    var stop = false

    for {
        select {
        case <-done:
            stop = true
        default:

            file, err := os.Open(path)
            if err != nil {
                log.Fatal(err)
            }

            fi, err := file.Stat()
            if err != nil {
                log.Fatal(err)
            }

            size := fi.Size()

            if size == 0 {
                size = 1
            }

            var percent = float64(size) / float64(total) * 100

            fmt.Printf("%.0f", percent)
            fmt.Println("%")
        }

        if stop {
            break
        }

        time.Sleep(time.Second)
    }
}

func DownloadFile(url string, dest string) {
    file := path.Base(url)

    log.Printf("Downloading file %s from %s\n", file, url)

    var localPath bytes.Buffer
    localPath.WriteString(dest)
    localPath.WriteString("/")
    localPath.WriteString(file)

    start := time.Now()

    out, err := os.Create(localPath.String())

    if err != nil {
        fmt.Println(localPath.String())
        panic(err)
    }

    defer out.Close()

    headResp, err := http.Head(url)

    if err != nil {
        panic(err)
    }

    defer headResp.Body.Close()

    size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))

    if err != nil {
        panic(err)
    }

    done := make(chan int64)

    go PrintDownloadPercent(done, localPath.String(), int64(size))

    resp, err := http.Get(url)

    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()

    n, err := io.Copy(out, resp.Body)

    if err != nil {
        panic(err)
    }

    done <- n

    elapsed := time.Since(start)
    log.Printf("Download completed in %s", elapsed)
}