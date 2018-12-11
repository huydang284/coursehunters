package main

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

            fmt.Printf("\r%s", generateProgressBar(percent))
        }

        if stop {
            break
        }

        time.Sleep(time.Second)
    }
}

func generateProgressBar(percent float64) (progress string) {
    done := int(percent / 2)

    progress = "["
    for i := 0; i < done-1; i++ {
        progress += "="
    }
    progress += ">"
    for i := 0; i < 50-done-1; i++ {
        progress += " "
    }
    progress += "] " + fmt.Sprintf("%.0f%%", percent)

    return progress
}

func downloadFile(url string, dest string) string {
    file := path.Base(url)

    log.Printf("Downloading file %s from %s\n", file, url)

    var localPath bytes.Buffer
    localPath.WriteString(dest)
    localPath.WriteString("/")
    localPath.WriteString(file)

    start := time.Now()

    out, err := os.Create(localPath.String())
    check(err)
    defer out.Close()

    headResp, err := http.Head(url)
    check(err)

    defer headResp.Body.Close()

    size, err := strconv.Atoi(headResp.Header.Get("Content-Length"))
    check(err)

    done := make(chan int64)

    go PrintDownloadPercent(done, localPath.String(), int64(size))

    resp, err := http.Get(url)
    check(err)
    defer resp.Body.Close()

    n, err := io.Copy(out, resp.Body)
    check(err)

    done <- n

    elapsed := time.Since(start)
    log.Println(fmt.Sprintf("Download completed in %s", elapsed))
    log.Println("Uploading to Youtube...")

    return file
}
